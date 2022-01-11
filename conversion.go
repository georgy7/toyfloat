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
	return encodeDelta(last, x, t)
}

func (t *Type) UseIntegerDelta(last uint16, delta int) uint16 {
	return decodeDelta(last, delta, t)
}

// Abs returns encoded absolute value of encoded argument.
func (t *Type) Abs(x uint16) uint16 {
	return x & (^t.minus)
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

	if b3 {
		denominator := 3.0
		for x := -1; x >= minX; x-- {
			settings.scale[x-minX] = 1.0 / denominator
			denominator *= 3.0
		}
		for x := 0; x <= maxX; x++ {
			settings.scale[x-minX] = math.Pow(base, float64(x))
		}
	} else {
		for x := minX; x < 0; x++ {
			settings.scale[x-minX] = 1.0 / powerOfTwo(uint8(-x))
		}
		for x := 0; x <= maxX; x++ {
			settings.scale[x-minX] = powerOfTwo(uint8(x))
		}
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
	if ((1.0 - 1e-14) < r) && (r < (1.0 + 1e-14)) {
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

func encodeDelta(last, x uint16, s *Type) int {
	filter := s.minus | (s.xMask << s.mSize) | s.mMask
	a := int(toSimple(last, s.minus) & filter)
	b := int(toSimple(x, s.minus) & filter)
	return b - a
}

func decodeDelta(last uint16, delta int, s *Type) uint16 {
	filter := s.minus | (s.xMask << s.mSize) | s.mMask
	filteredSimpleLast := int(toSimple(last, s.minus) & filter)

	r := uint16(0)
	if delta > int(filter)-filteredSimpleLast {
		r = filter
	} else if delta >= -filteredSimpleLast {
		r = uint16(filteredSimpleLast + delta)
	}

	return fromSimple(r, s.minus)
}

func toSimple(tf, minus uint16) uint16 {
	if 0 == tf&minus {
		return minus | tf
	}
	return ^tf
}

func fromSimple(simple, minus uint16) uint16 {
	if minus != simple&minus {
		return ^simple
	}
	return (^minus) & simple
}

func isNegative(tf, minus uint16) bool {
	return 0b0 != tf&minus
}

func powerOfTwo(x uint8) float64 {
	return float64(int(1) << x)
}
