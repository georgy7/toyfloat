// Package toyfloat provides tiny (3 to 16 bits)
// floating-point number formats for serialization.
package toyfloat

import (
	"errors"
	"math"
)

// Type is a reusable immutable set of encoder settings.
type Type struct {
	mSize         uint8
	minus         uint16
	mMask         uint16
	twoPowerMSize float64
	dsFactor      float64
	boundary      float64
	xc            xConstants
}

func NewTypeX2(length int, signed bool) (Type, error) {
	return newSettings(uint8(length), 2, -3, true, signed)
}

func NewTypeX3(length int, signed bool) (Type, error) {
	return newSettings(uint8(length), 3, -6, false, signed)
}

func NewTypeX4(length int, signed bool) (Type, error) {
	return newSettings(uint8(length), 4, -8, false, signed)
}

func (t *Type) Encode(v float64) uint16 {
	return encode(v, t)
}

func (t *Type) Decode(x uint16) float64 {
	return decode(x, t)
}

func (t *Type) GetIntegerDelta(last uint16, x uint16) int {
	return encodeDelta(last, x, t.minus)
}

func (t *Type) UseIntegerDelta(last uint16, delta int) uint16 {
	return decodeDelta(last, delta, t)
}

// ----------------
// Implementation:

const maxPossibleScaleIndex = 0xF
const scaleArraySize = maxPossibleScaleIndex + 1

type xConstants struct {
	xSize uint8
	xMask uint16
	base3 bool
	scale [scaleArraySize]float64
}

func newSettings(length, xSize uint8, minX int, b3, signed bool) (Type, error) {
	if length > 16 {
		return Type{}, errors.New("maximum length is 16 bits")
	}

	signSize := uint8(0)
	if signed {
		signSize = 1
	}

	if length <= xSize+signSize {
		return Type{}, errors.New("mantissa must be at least 1 bit wide")
	}

	mSize := length - (xSize + signSize)

	minus := uint16(0)
	if signed {
		minus = uint16(1) << (length - 1)
	}

	if (xSize >= 16) || ((int(1) << xSize) > scaleArraySize) {
		return Type{}, errors.New("such big exponents are not supported")
	}

	maxX := minX + (int(1) << xSize) - 1

	var scale [scaleArraySize]float64
	for x := minX; x <= maxX; x++ {
		base := 2.0
		if b3 {
			base = 3.0
		}
		scale[x-minX] = math.Pow(base, float64(x))
	}

	return Type{
		mSize:         mSize,
		minus:         minus,
		mMask:         getOneBits(mSize),
		twoPowerMSize: powerOfTwo(mSize),
		dsFactor:      makeDecodeSignificandFactor(mSize, b3),
		boundary:      makeBoundaryBetweenExponents(mSize, b3),
		xc: xConstants{
			xSize: xSize,
			xMask: getOneBits(xSize),
			base3: b3,
			scale: scale,
		}}, nil
}

func getOneBits(bits uint8) uint16 {
	if bits >= 16 {
		return ^uint16(0)
	}
	return (uint16(1) << bits) - 1
}

func encode(value float64, settings *Type) uint16 {
	if math.IsNaN(value) {
		return 0x0
	}

	negativeValue := value < 0

	if negativeValue && (0b0 == settings.minus) {
		return 0x0
	}

	a := settings.xc.scale[0]
	vReversedC := value * (1.0 - a)

	if negativeValue {
		return settings.minus | encodeInnerValue(a-vReversedC, settings)
	}
	return encodeInnerValue(a+vReversedC, settings)
}

func decode(tf uint16, s *Type) float64 {
	a := s.xc.scale[0]
	c := 1.0 / (1.0 - a)

	scale := s.xc.scale[(tf>>s.mSize)&s.xc.xMask]

	significand := decodeSignificand(float64(tf&s.mMask), s.dsFactor)

	r := significand * scale
	r = (r - a) * c

	// The problem only appeared with base three exponent.
	if math.Abs(r-1.0) < 1e-14 {
		r = 1.0
	}

	if isNegative(tf, s.minus) {
		return -r
	}
	return r
}

func encodeInnerValue(inner float64, s *Type) uint16 {
	mMax := s.twoPowerMSize - 1.0

	maxScaleIndex := getMaxScaleIndex(s)
	internalMaximum := decodeSignificand(mMax, s.dsFactor)
	internalMaximum *= s.xc.scale[maxScaleIndex]

	if inner >= internalMaximum {
		return (s.xc.xMask << s.mSize) | s.mMask
	}

	binaryExponent, inverseScale := getBinaryExponent(inner, maxScaleIndex, s)

	significand := inner * inverseScale
	mFloat := encodeSignificand(significand, s)

	// math.Round(x) = math.Floor(x + 0.5), x >= 0
	binarySignificand := uint16(mFloat + 0.5)

	// If method getBinaryExponent works as intended,
	// this is always false.
	if binarySignificand > s.mMask {
		binarySignificand = s.mMask
	}

	return binarySignificand | binaryExponent
}

func encodeSignificand(significand float64, s *Type) float64 {
	mFloat := significand - 1.0
	if s.xc.base3 {
		return mFloat * powerOfTwo(s.mSize-1)
	}
	return mFloat * s.twoPowerMSize
}

func decodeSignificand(m, dsFactor float64) float64 {
	return 1.0 + m*dsFactor
}

func makeDecodeSignificandFactor(mSize uint8, base3 bool) float64 {
	if base3 {
		return 1.0 / powerOfTwo(mSize-1)
	}
	return 1.0 / powerOfTwo(mSize)
}

func makeBoundaryBetweenExponents(mSize uint8, base3 bool) float64 {
	base := 2.0
	if base3 {
		base = 3.0
	}

	mDiv2m := (powerOfTwo(mSize) - 0.5) / powerOfTwo(mSize)

	// This is the part (1 + (b - 1) * m/(2^M)) of the formula,
	// that should be rounded to a greater exponent.
	return 1 + (base-1)*mDiv2m
}

func getBinaryExponent(inner float64, maxScaleIndex uint8, s *Type) (uint16, float64) {
	modulus := math.Abs(inner)
	factor := s.boundary

	start := uint16(maxScaleIndex)
	scale := s.xc.scale[start]

	var result uint16
	for result = start; result > 0; result-- {
		smallerScale := s.xc.scale[result-1]
		if factor*smallerScale <= modulus {
			break
		}
		scale = smallerScale
	}

	return result << s.mSize, 1.0 / scale
}

func getMaxScaleIndex(s *Type) uint8 {
	r := (uint8(1) << s.xc.xSize) - 1
	return r & maxPossibleScaleIndex
}

func encodeDelta(last, x, minus uint16) int {
	lastIsNegative := isNegative(last, minus)
	xIsNegative := isNegative(x, minus)
	sameSign := lastIsNegative == xIsNegative

	absLast := int(abs(last, minus))
	absX := int(abs(x, minus))

	if sameSign {
		diff := absX - absLast
		if lastIsNegative {
			return -diff
		}
		return diff
	}

	// the additional 1 is minus zero
	sum := absX + 1 + absLast
	if lastIsNegative {
		return sum
	}
	return -sum
}

func decodeDelta(last uint16, delta int, s *Type) uint16 {
	if delta == 0 {
		return last
	}

	lastIsNegative := isNegative(last, s.minus)

	diffOrNegativeSum := delta
	if lastIsNegative {
		diffOrNegativeSum = -diffOrNegativeSum
	}

	// diff[OrNegativeSum] + absLast = absX (sameSign is true)
	// [diffOr]NegativeSum + absLast = -absX - 1

	r := diffOrNegativeSum + int(abs(last, s.minus))
	sameSign := r >= 0

	maxValue := (s.xc.xMask << s.mSize) | s.mMask

	if sameSign {
		rLimited := min16u(uint16(r), maxValue)
		if lastIsNegative {
			return s.minus | rLimited
		}
		return rLimited
	}

	if !lastIsNegative && (s.minus == 0b0) {
		return 0b0
	}

	rLimited := min16u(uint16(-r-1), maxValue)

	if lastIsNegative {
		return rLimited
	}

	return s.minus | rLimited
}

func isNegative(tf, minus uint16) bool {
	return 0b0 != tf&minus
}

func abs(tf, minus uint16) uint16 {
	return tf & (^minus)
}

func powerOfTwo(x uint8) float64 {
	return float64(int(1) << x)
}

func min16u(a, b uint16) uint16 {
	if a < b {
		return a
	}
	return b
}
