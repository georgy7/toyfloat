// Package toyfloat provides tiny (3 to 16 bits)
// floating-point number formats for serialization.
package toyfloat

import (
	"errors"
	"math"
)

// Type is a reusable immutable set of encoder settings.
type Type struct {
	mSize     int
	minus     uint16
	mMask     uint16
	twoPowerM float64
	xc        xConstants
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
	return encodeDelta(last, x, t)
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
		mSize -= 1
	}

	if mSize < 1 {
		return Type{}, errors.New("mantissa must be at least 1 bit wide")
	}

	minus := uint16(0)
	if signed {
		minus = uint16(1) << (length - 1)
	}

	mMask := getOnes(mSize)
	maxX := minX + (int(1) << xSize) - 1

	if xSize > 4 {
		return Type{}, errors.New("such big exponents are not supported")
	}

	var scales [16]float64
	for x := minX; x <= maxX; x++ {
		scales[x-minX] = calculateScale(b3, x)
	}

	return Type{mSize, minus, mMask,
		powerOfTwo(mSize),
		xConstants{
			xMask:       getOnes(xSize),
			xSize:       xSize,
			minExponent: minX,
			maxExponent: maxX,
			base3:       b3,
			scales:      scales,
		}}, nil
}

func isNegative(tf uint16, settings *Type) bool {
	return 0b0 != tf&(settings.minus)
}

func abs(tf uint16, settings *Type) uint16 {
	return tf & (^settings.minus)
}

func encode(value float64, settings *Type) uint16 {
	if math.IsNaN(value) {
		return 0x0
	}

	nv := value < 0

	if nv && (0b0 == settings.minus) {
		return 0x0
	}

	a := settings.xc.scales[0]
	vReversedC := value * (1.0 - a)

	if nv {
		return settings.minus | encodeInnerValue(a-vReversedC, settings)
	} else {
		return encodeInnerValue(a+vReversedC, settings)
	}
}

func decode(tf uint16, settings *Type) float64 {
	a := settings.xc.scales[0]
	c := 1.0 / (1.0 - a)

	xShift, xMask := settings.mSize, settings.xc.xMask
	mMask := settings.mMask

	x := int((tf>>xShift)&xMask) + settings.xc.minExponent

	significand := decodeSignificand(
		float64(tf&mMask), settings.twoPowerM, settings.xc.base3)

	r := significand * getScale(x, settings)
	r = (r - a) * c

	// The problem only appeared with base three exponent.
	if math.Abs(r-1.0) < 1e-14 {
		r = 1.0
	}

	if 0b0 != tf&settings.minus {
		return -r
	}

	return r
}

func encodeInnerValue(inner float64, s *Type) uint16 {

	xShift, xMask := s.mSize, s.xc.xMask
	xBias := s.xc.minExponent

	mMax := s.twoPowerM - 1.0

	internalMaximum := decodeSignificand(mMax, s.twoPowerM, s.xc.base3)
	internalMaximum *= getScale(s.xc.maxExponent, s)

	if inner >= internalMaximum {
		return (xMask << xShift) | s.mMask
	}

	x, scale := getExponent(inner, s)
	binaryExponent := uint16(x-xBias) << xShift

	normalized := inner / scale

	mFloat := encodeSignificand(normalized, s)

	binarySignificand := uint16(math.Min(math.Round(mFloat), mMax))

	return binarySignificand | binaryExponent
}

func encodeSignificand(normalized float64, s *Type) float64 {
	mFloat := normalized - 1.0
	if s.xc.base3 {
		return mFloat * powerOfTwo(s.mSize-1)
	} else {
		return mFloat * s.twoPowerM
	}
}

func decodeSignificand(m, twoPowerM float64, base3 bool) float64 {
	r := m
	if base3 {
		r /= twoPowerM * 0.5
	} else {
		r /= twoPowerM
	}
	return 1.0 + r
}

func getExponent(innerValue float64, s *Type) (int, float64) {
	modulus := math.Abs(innerValue)

	eps := 0.5 / s.twoPowerM
	var factor float64

	if s.xc.base3 {
		factor = 3.0 - (2.0 * eps)
	} else {
		factor = 2.0 - eps
	}

	scale := getScale(s.xc.maxExponent, s)

	for exp := s.xc.maxExponent; exp > s.xc.minExponent; exp-- {
		smallerScale := getScale(exp-1, s)
		if factor*smallerScale <= modulus {
			return exp, scale
		} else {
			scale = smallerScale
		}
	}

	return s.xc.minExponent, scale
}

func getOnes(bits int) uint16 {
	count := bits
	if count > 16 {
		count = 16
	}

	const bit uint16 = 1

	result := uint16(0)
	for i := 0; i < count; i++ {
		result |= bit << i
	}

	return result
}

func powerOfTwo(x int) float64 {
	if x < 0 {
		return 1.0 / float64(int(1)<<-x)
	} else {
		return float64(int(1) << x)
	}
}

func powerOfThree(x int) float64 {
	root := x < 0

	var absX int
	if root {
		absX = -x
	} else {
		absX = x
	}

	rInt := 1
	for i := 0; i < absX; i++ {
		rInt *= 3
	}

	if root {
		return 1.0 / float64(rInt)
	} else {
		return float64(rInt)
	}
}

func calculateScale(base3 bool, x int) float64 {
	if base3 {
		return powerOfThree(x)
	} else {
		return powerOfTwo(x)
	}
}

func getScale(x int, s *Type) float64 {
	index := x - s.xc.minExponent
	if index <= 0 {
		return s.xc.scales[0]
	} else {
		return s.xc.scales[index%len(s.xc.scales)]
	}
}

func min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func encodeDelta(last uint16, x uint16, settings *Type) int {
	lastIsNegative := isNegative(last, settings)
	xIsNegative := isNegative(x, settings)
	sameSign := lastIsNegative == xIsNegative

	absLast := int(abs(last, settings))
	absX := int(abs(x, settings))

	if sameSign {
		diff := absX - absLast
		if xIsNegative {
			return -diff
		} else {
			return diff
		}
	} else {
		// the additional 1 is minus zero
		sum := absX + 1 + absLast
		if lastIsNegative {
			return sum
		} else {
			return -sum
		}
	}
}

func decodeDelta(last uint16, delta int, s *Type) uint16 {
	if delta == 0 {
		return last
	}

	absLast := int(abs(last, s))

	xShift, xMask := s.mSize, s.xc.xMask
	mMask := s.mMask
	maxValue := (xMask << xShift) | mMask

	if isNegative(last, s) {
		absX := min(absLast-delta, int(maxValue))
		if absX >= 0 {
			return s.minus | uint16(absX)
		} else {
			return uint16(min(-(absX + 1), int(maxValue)))
		}
	} else {
		absX := min(absLast+delta, int(maxValue))
		if absX >= 0 {
			return uint16(absX)
		} else if s.minus == 0b0 {
			return 0b0
		} else {
			return s.minus | uint16(min(-(absX+1), int(maxValue)))
		}
	}
}
