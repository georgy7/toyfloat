// Package toyfloat provides tiny (3 to 16 bits)
// floating-point number formats for serialization.
package toyfloat

import (
	"errors"
	"math"
)

// Type is a reusable immutable set of encoder settings.
type Type struct {
	mSize               uint8
	minus, mMask, xMask uint16
	minValue, maxValue  float64
	esFactor, dsFactor  float64
	xBoundary           float64
	scale               []float64
	bitmask             uint16
}

// NewTypeX2 makes a type with 2-bit exponent with default settings.
func NewTypeX2(length int, signed bool) (Type, error) {
	return NewType(uint8(length), 3, 2, -3, signed)
}

// NewTypeX3 makes a type with 3-bit exponent with default settings.
func NewTypeX3(length int, signed bool) (Type, error) {
	return NewType(uint8(length), 2, 3, -6, signed)
}

// NewTypeX4 makes a type with 4-bit exponent with default settings.
func NewTypeX4(length int, signed bool) (Type, error) {
	return NewType(uint8(length), 2, 4, -8, signed)
}

// NewType allows creating custom types.
// Use it at your own risk.
// The argument minX is the minimum power of the exponential part of a number.
// The argument xSize is the number of bits that encode the power.
// So the maximum power equals minX+(2^xSize)-1,
// and the maximum exponential part equals xBase^(minX+(2^xSize)-1).
func NewType(length, xBase, xSize uint8, minX int, signed bool) (Type, error) {
	if (xBase != 2) && (xBase != 3) {
		return Type{},
			errors.New("only base 2 and base 3 exponents are supported rn")
	}

	if minX >= 0 {
		return Type{}, errors.New("c=1/(1-xBase^minX)" +
			" where it is assumed that minX is not positive" +
			" so that с makes sense")
	}

	return newSettings(length, xBase, xSize, minX, signed)
}

// Encode converts a number to its binary representation for this type.
// You cannot compare such values directly because they are "sign–magnitude".
// Of course, they have zeros in extra most-significant bits.
func (t *Type) Encode(v float64) uint16 {
	return encode(v, t)
}

// Decode is just method Encode in reverse.
// It ignores values of extra most-significant bits.
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
// This does not work for the comparable form.
func (t *Type) Abs(x uint16) uint16 {
	return x & (^t.minus)
}

// ToComparable returns a representation close to "ones' complement",
// except for its sign bit reversed.
// Thus, all zeros mean the lowest value, and all ones mean the maximum.
// Programming languages such as C, C++, Go define unsigned integer overflow,
// which allows this form to be used for delta encoding without branching.
func (t *Type) ToComparable(tf uint16) uint16 {
	var r uint16
	if 0 == tf&t.minus {
		// It's true for both positive signed and unsigned numbers.
		r = t.minus | tf
	} else {
		// Negative, including -0.
		r = ^tf
	}
	return r & t.bitmask
}

// FromComparable is ToComparable in reverse.
// Note, that it does not reset extra bits (for performance reasons).
func (t *Type) FromComparable(c uint16) uint16 {
	// Sign bit are inverted here, so it is
	// not equal to its bitmask for a negative number.
	// Also, variable "minus" equals zero for unsigned values,
	// "0 != 0" is always false.
	if t.minus != c&t.minus {
		// Negative, including -0.
		return ^c
	}
	return (^t.minus) & c
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

func newSettings(length, xBase, xSize uint8, minX int, signed bool) (Type, error) {
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

	if (xSize >= 16) || (mSize >= 16) {
		return Type{}, errors.New("library Toyfloat is broken")
	}

	settings := Type{
		mSize: mSize,
		minus: uint16(0),
		mMask: (uint16(1) << mSize) - 1,
		xMask: (uint16(1) << xSize) - 1,
	}

	if signed {
		settings.minus = uint16(1) << (length - 1)
	}

	settings.bitmask =
		settings.minus | (settings.xMask << settings.mSize) | settings.mMask

	// multiplier to encode the significand
	settings.esFactor = powerOfTwo(mSize) / float64(xBase-1)
	// multiplier to decode it
	settings.dsFactor = 1.0 / settings.esFactor

	f64Base := float64(xBase)
	settings.xBoundary = makeExponentBoundary(powerOfTwo(mSize), f64Base)

	settings.scale = make([]float64, int(1)<<xSize)
	maxX := minX + len(settings.scale) - 1

	{
		denominator := f64Base
		for x := -1; x >= minX; x-- {
			settings.scale[x-minX] = 1.0 / denominator
			denominator *= f64Base
		}
		for x := 0; x <= maxX; x++ {
			settings.scale[x-minX] = math.Pow(f64Base, float64(x))
		}
	}

	mMax := powerOfTwo(mSize) - 1.0
	maxScale := get(settings.scale, settings.xMask)
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

	scale := get(s.scale, (tf>>s.mSize)&(s.xMask))

	significand := decodeSignificand(float64(tf&s.mMask), s.dsFactor)

	absValue := ((significand * scale) - a) * c

	// The problem only appeared with base three exponent.
	if ((1.0 - 1e-14) < absValue) && (absValue < (1.0 + 1e-14)) {
		absValue = 1.0
	}

	if isNegative(tf, s.minus) {
		return -absValue
	}
	return absValue
}

func encodeInnerValue(inner float64, s *Type) uint16 {
	binaryExponent, inverseScale := getBinaryExponent(inner, s)
	denominator := s.esFactor

	// math.Round(x) = math.Floor(x + 0.5), x >= 0
	const rounding = 0.499999999999

	// I need to find m from (1+(b-1)(m/2^M))(b^x), which is named "inner" here.
	//
	// "inverseScale" = 1/(b^x)
	// "denominator" = 1/((b-1)(1/(2^M))). It is reversed dsFactor.
	// It's called denominator because it's equals 2^M for base 2 exponents.
	// It's an integer power of two for both base 2 and base 3 exponents.
	// So,
	// inner = (1+(b-1)(m/2^M)) * (b^x)
	// inner = (1+(b-1)(m/2^M)) * (1/inverseScale)
	// inner * inverseScale = 1 + (b-1)(m/2^M)
	// (inner * inverseScale) - 1 = (b-1)(m/2^M)
	// (inner * inverseScale) - 1 = (b-1)(1/2^M)m
	// (inner * inverseScale) - 1 = (1/denominator)m
	// ((inner * inverseScale) - 1) * denominator = m
	// m = ((inner * inverseScale) - 1) * denominator
	// m = ((inner * inverseScale * denominator) - denominator)
	// m = inner * inverseScale * denominator - denominator
	//
	// inner >= 0 always for this method,
	// denominator is a natural number,
	// inverseScale is a positive rational number.
	// Thus, m can not be negative.
	//
	// Method getBinaryExponent ensures that m < (2^M)-0.5,
	// so, in theory, m can be safely rounded to the closest integer,
	// however since this is floating-point arithmetic,
	// I am afraid, it might somehow be >= (2^M)-0.5.
	// So I use the constant slightly less than one-half for rounding.
	significand := denominator*inner*inverseScale - denominator + rounding

	return uint16(significand) | binaryExponent
}

func decodeSignificand(m, dsFactor float64) float64 {
	return 1.0 + m*dsFactor
}

func makeExponentBoundary(twoPowerMSize, base float64) float64 {
	// This is the part (1 + (b - 1) * m/(2^M)) of the formula,
	// that should be rounded to a greater exponent.
	// Maximum integer m equals 2^M - 1.
	// So, a floating-point m before rounding must be less than 2^M - 0.5.
	mDiv2m := (twoPowerMSize - 0.5) / twoPowerMSize
	return 1 + (base-1)*mDiv2m
}

func getBinaryExponent(absValue float64, s *Type) (uint16, float64) {
	xb := s.xBoundary

	// It's biased in the sense
	// that zero means the minimum exponent of the type.
	biasedExponent := s.xMask

	// This method must find such minimum x, that m < (2^M)-0.5.
	// So, this loop finds the maximum x for the following condition:
	// absValue >= (1 + (b-1) * ((2^M)-0.5)/(2^M)) * b^(x-1)
	// Thus, m values >= (2^M)-0.5 will lead to selection of a larger x.
	// Then, the method will return the corresponding scale,
	// and this will result in a lower m (non-negative anyway).
	//
	// The case, with m >= (2^M)-0.5 for the maximum possible x,
	// is filtered in the beginning of method "encode".
	// If those checks fail to filter out values that are out of range,
	// it will lead to an integer overflow.
	for (biasedExponent > 0) && (xb*get(s.scale, biasedExponent-1) > absValue) {
		biasedExponent--
	}

	// This is an exponential part of encoded number: b^x.
	scale := get(s.scale, biasedExponent)
	// By some reason, multiplying by a non-constant inverse number
	// is faster, than division on my computer. So I return the inverse scale.
	return biasedExponent << s.mSize, 1.0 / scale
}

func encodeDelta(last, x uint16, s *Type) int {
	a := int(s.ToComparable(last))
	b := int(s.ToComparable(x))
	return b - a
}

func decodeDelta(last uint16, delta int, s *Type) uint16 {
	lastComparable := int(s.ToComparable(last))

	r := uint16(0)
	if delta > int(s.bitmask)-lastComparable {
		r = s.bitmask
	} else if delta >= -lastComparable {
		r = uint16(lastComparable + delta)
	}

	return s.FromComparable(r)
}

func isNegative(tf, minus uint16) bool {
	return 0b0 != tf&minus
}

func powerOfTwo(x uint8) float64 {
	return float64(int(1) << x)
}

func get(s []float64, i uint16) float64 {
	maxIndex := len(s) - 1
	if maxIndex < 0 {
		return 0.0
	}

	if int(i) > maxIndex {
		return s[maxIndex]
	}
	return s[i]
}
