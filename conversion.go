// Package toyfloat provides tiny (3 to 16 bits)
// floating-point number formats for serialization.
package toyfloat

import (
	"errors"
	"math"
)

const maxPossibleScaleIndex = 0xF
const scaleArraySize = maxPossibleScaleIndex + 1

// Type is a reusable immutable set of encoder settings.
type Type struct {
	mSize               uint8
	minus, mMask        uint16
	encodingDenominator float64
	minValue, maxValue  float64
	dsFactor, boundary  float64
	xMask               uint16
	base3               bool
	scale               [scaleArraySize]float64
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

// Abs returns encoded absolute value of encoded argument.
func (t *Type) Abs(x uint16) uint16 {
	return abs(x, t.minus)
}

// MinValue returns zero for unsigned types and negative
// with maximum absolute value for signed types.
func (t *Type) MinValue() float64 {
	return t.minValue
}

// MaxValue returns maximum value of the type.
func (t *Type) MaxValue() float64 {
	return t.maxValue
}

// ----------------
// Implementation:

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

	if (xSize >= 16) || ((int(1) << xSize) > scaleArraySize) || (mSize >= 16) {
		return Type{}, errors.New("library Toyfloat is broken")
	}

	settings := Type{
		mSize:    mSize,
		minus:    uint16(0),
		mMask:    (uint16(1) << mSize) - 1,
		dsFactor: makeDecodeSignificandFactor(mSize, b3),
		xMask:    (uint16(1) << xSize) - 1,
		base3:    b3,
	}

	if signed {
		settings.minus = uint16(1) << (length - 1)
	}

	maxX := minX + (int(1) << xSize) - 1

	base := 2.0
	settings.encodingDenominator = powerOfTwo(mSize)
	if b3 {
		base = 3.0
		settings.encodingDenominator = powerOfTwo(mSize - 1)
	}

	settings.boundary = makeExponentBoundary(powerOfTwo(mSize), base)

	for x := minX; x <= maxX; x++ {
		settings.scale[x-minX] = math.Pow(base, float64(x))
	}

	mMax := powerOfTwo(mSize) - 1.0
	maxScale := settings.scale[settings.xMask&maxPossibleScaleIndex]
	internalMaximum := decodeSignificand(mMax, settings.dsFactor) * maxScale

	a := settings.scale[0]
	c := 1.0 / (1.0 - a)

	settings.maxValue = (internalMaximum - a) * c
	settings.minValue = 0.0
	if signed {
		settings.minValue = -settings.maxValue
	}

	return settings, nil
}

func encode(value float64, settings *Type) uint16 {
	if math.IsNaN(value) {
		return 0x0
	} else if value > settings.maxValue {
		return (settings.xMask << settings.mSize) | settings.mMask
	} else if value < 0 {
		if 0b0 == settings.minus {
			return 0x0
		} else if value < settings.minValue {
			absValue := (settings.xMask << settings.mSize) | settings.mMask
			return settings.minus | absValue
		}
	}

	a := settings.scale[0]
	vReversedC := value * (1.0 - a)

	if value < 0 {
		return settings.minus | encodeInnerValue(a-vReversedC, settings)
	}
	return encodeInnerValue(a+vReversedC, settings)
}

func decode(tf uint16, s *Type) float64 {
	a := s.scale[0]
	c := 1.0 / (1.0 - a)

	scale := s.scale[(tf>>s.mSize)&(s.xMask&maxPossibleScaleIndex)]

	significand := decodeSignificand(float64(tf&s.mMask), s.dsFactor)

	r := ((significand * scale) - a) * c

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
	binaryExponent, inverseScale := getBinaryExponent(inner, s)
	denominator := s.encodingDenominator

	// math.Round(x) = math.Floor(x + 0.5), x >= 0
	const rounding = 0.499999999999

	significand := denominator*inner*inverseScale - denominator + rounding

	return uint16(significand) | binaryExponent
}

func decodeSignificand(m, dsFactor float64) float64 {
	return 1.0 + m*dsFactor
}

func makeDecodeSignificandFactor(mSize uint8, base3 bool) float64 {
	if base3 {
		// (b-1) * 1/(2^M)
		// (3-1) * 1/(2^M)
		// (2^1) * 1/(2^M)
		// (2^1) / (2^M)
		// (2^1) * (2^-M)
		// 2^(1-M)
		// 1 / 2^(-(1-M))
		// 1 / 2^(M-1))
		return 1.0 / powerOfTwo(mSize-1)
	}
	return 1.0 / powerOfTwo(mSize)
}

func makeExponentBoundary(twoPowerMSize, base float64) float64 {
	// This is the part (1 + (b - 1) * m/(2^M)) of the formula,
	// that should be rounded to a greater exponent.
	mDiv2m := (twoPowerMSize - 0.5) / twoPowerMSize
	return 1 + (base-1)*mDiv2m
}

func getBinaryExponent(inner float64, s *Type) (uint16, float64) {
	absValue := math.Abs(inner)
	factor := s.boundary

	result := s.xMask & maxPossibleScaleIndex

	for (result > 0) && (factor*s.scale[result-1] > absValue) {
		result--
	}

	scale := s.scale[result]
	return result << s.mSize, 1.0 / scale
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

	maxValue := (s.xMask << s.mSize) | s.mMask

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
