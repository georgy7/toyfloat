// Package toyfloat provides tiny (3 to 16 bits)
// floating-point number formats for serialization.
package toyfloat

import (
	"errors"
	"math"
)

// Type is a reusable immutable set of encoder settings.
type Type struct {
	mSize             int
	minus             uint16
	mMask             uint16
	twoPowerM         float64
	dsFactor          float64
	getExponentFactor float64
	xc                xConstants
}

func NewTypeX2(length int, signed bool) (Type, error) {
	return newSettings(length, 2, -3, true, signed)
}

func NewTypeX3(length int, signed bool) (Type, error) {
	return newSettings(length, 3, -6, false, signed)
}

func NewTypeX4(length int, signed bool) (Type, error) {
	return newSettings(length, 4, -8, false, signed)
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

type xConstants struct {
	xMask       uint16
	xSize       int
	minExponent int
	maxExponent int
	base3       bool
	scales      [16]float64
}

func newSettings(length int, xSize, minX int, b3, signed bool) (Type, error) {
	if length > 16 {
		return Type{}, errors.New("maximum length is 16 bits")
	}

	mSize := length - xSize
	if signed {
		mSize--
	}

	if mSize < 1 {
		return Type{}, errors.New("mantissa must be at least 1 bit wide")
	}

	minus := uint16(0)
	if signed {
		minus = uint16(1) << (length - 1)
	}

	maxX := minX + (int(1) << xSize) - 1

	if xSize > 4 {
		// Because of scales array size (2^4 = 16).
		return Type{}, errors.New("such big exponents are not supported")
	}

	var scales [16]float64
	for x := minX; x <= maxX; x++ {
		base := 2.0
		if b3 {
			base = 3.0
		}
		scales[x-minX] = math.Pow(base, float64(x))
	}

	return Type{
		mSize:             mSize,
		minus:             minus,
		mMask:             getOneBits(mSize),
		twoPowerM:         powerOfTwo(mSize),
		dsFactor:          makeDecodeSignificandFactor(mSize, b3),
		getExponentFactor: makeGetExponentFactor(mSize, b3),
		xc: xConstants{
			xMask:       getOneBits(xSize),
			xSize:       xSize,
			minExponent: minX,
			maxExponent: maxX,
			base3:       b3,
			scales:      scales,
		}}, nil
}

func getOneBits(bits int) uint16 {
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

	a := settings.xc.scales[0]
	vReversedC := value * (1.0 - a)

	if negativeValue {
		return settings.minus | encodeInnerValue(a-vReversedC, settings)
	}
	return encodeInnerValue(a+vReversedC, settings)
}

func decode(tf uint16, s *Type) float64 {
	a := s.xc.scales[0]
	c := 1.0 / (1.0 - a)

	scale := s.xc.scales[int((tf>>s.mSize)&s.xc.xMask)]

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
	mMax := s.twoPowerM - 1.0

	internalMaximum := decodeSignificand(mMax, s.dsFactor)
	internalMaximum *= s.xc.scales[getMaxScaleIndex(s)]

	if inner >= internalMaximum {
		return (s.xc.xMask << s.mSize) | s.mMask
	}

	x, scale := getExponent(inner, s)

	xBias := s.xc.minExponent
	binaryExponent := uint16(x-xBias) << s.mSize

	significand := inner / scale
	mFloat := encodeSignificand(significand, s)
	binarySignificand := uint16(math.Round(mFloat))

	intMMax := (uint16(1) << s.mSize) - 1
	if binarySignificand > intMMax {
		binarySignificand = intMMax
	}

	return binarySignificand | binaryExponent
}

func encodeSignificand(significand float64, s *Type) float64 {
	mFloat := significand - 1.0
	if s.xc.base3 {
		return mFloat * powerOfTwo(s.mSize-1)
	}
	return mFloat * s.twoPowerM
}

func decodeSignificand(m, dsFactor float64) float64 {
	return 1.0 + m*dsFactor
}

func makeDecodeSignificandFactor(mSize int, base3 bool) float64 {
	if base3 {
		return 1.0 / powerOfTwo(mSize-1)
	}
	return 1.0 / powerOfTwo(mSize)
}

func makeGetExponentFactor(mSize int, base3 bool) float64 {
	eps := 0.5 / powerOfTwo(mSize)

	base := 2.0
	if base3 {
		base = 3.0
	}

	return base - ((base - 1) * eps)
}

func getExponent(innerValue float64, s *Type) (int, float64) {
	modulus := math.Abs(innerValue)
	factor := s.getExponentFactor

	start := getMaxScaleIndex(s)
	scale := s.xc.scales[start]

	for i := start; i >= 1; i-- {
		smallerScale := s.xc.scales[i-1]
		if factor*smallerScale <= modulus {
			return s.xc.minExponent + i, scale
		}

		scale = smallerScale
	}

	return s.xc.minExponent, scale
}

func getMaxScaleIndex(s *Type) int {
	r := s.xc.maxExponent - s.xc.minExponent
	if r >= len(s.xc.scales) {
		return len(s.xc.scales) - 1
	}
	return r
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

	// diff[OrNegativeSum] + absLast = absX	(sameSign)
	// [diffOr]NegativeSum + absLast = -absX - 1

	r := diffOrNegativeSum + int(abs(last, s.minus))
	sameSign := r >= 0

	maxValue := int((s.xc.xMask << s.mSize) | s.mMask)

	if sameSign {
		rLimited := uint16(min(r, maxValue))
		if lastIsNegative {
			return s.minus | rLimited
		}
		return rLimited
	}

	if !lastIsNegative && (s.minus == 0b0) {
		return 0b0
	}

	rLimited := uint16(min(-r-1, maxValue))

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

func powerOfTwo(x int) float64 {
	if x < 0 {
		return 1.0 / float64(int(1)<<-x)
	}
	return float64(int(1) << x)
}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}
