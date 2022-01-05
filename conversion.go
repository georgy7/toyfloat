// Package toyfloat provides tiny (4 to 16 bits)
// floating-point number formats for serialization.
package toyfloat

import (
	"errors"
	"math"
)

// Type is a reusable immutable set of encoder settings.
type Type struct {
	settings settings
}

func NewTypeX2(length int, signed bool) (Type, error) {
	s, e := newSettings(length, x2(), signed)
	return Type{s}, e
}

func NewTypeX3(length int, signed bool) (Type, error) {
	s, e := newSettings(length, x3(), signed)
	return Type{s}, e
}

func NewTypeX4(length int, signed bool) (Type, error) {
	s, e := newSettings(length, x4(), signed)
	return Type{s}, e
}

func (t *Type) Encode(v float64) uint16 {
	return encode(v, &t.settings)
}

func (t *Type) Decode(x uint16) float64 {
	return decode(x, &t.settings)
}

func (t *Type) GetIntegerDelta(last uint16, x uint16) int {
	return encodeDelta(last, x, &t.settings)
}

func (t *Type) UseIntegerDelta(last uint16, delta int) uint16 {
	return decodeDelta(last, delta, &t.settings)
}

// ----------------
// Implementation:

type xConstants struct {
	xMask               uint16
	xSize               int
	minExponent         int
	maxExponent         int
	twoPowerMinExponent float64
	twoPowerMaxExponent float64
}

type settings struct {
	mSize  int
	xShift int
	minus  uint16
	mMask  uint16
	xc     xConstants
}

func x4() xConstants {
	return xConstants{
		xMask:               0b1111,
		xSize:               4,
		minExponent:         -8,
		maxExponent:         -8 + 15,
		twoPowerMinExponent: 1.0 / 256.0,
		twoPowerMaxExponent: 128.0,
	}
}

func x3() xConstants {
	return xConstants{
		xMask:               0b111,
		xSize:               3,
		minExponent:         -6,
		maxExponent:         -6 + 7,
		twoPowerMinExponent: 1.0 / 64.0,
		twoPowerMaxExponent: 2.0,
	}
}

func x2() xConstants {
	return xConstants{
		xMask:               0b11,
		xSize:               2,
		minExponent:         -3,
		maxExponent:         0,
		twoPowerMinExponent: 1.0 / 8.0,
		twoPowerMaxExponent: 1.0,
	}
}

func makeSettings(mSize int, minus, mMask uint16, xc xConstants) settings {
	return settings{mSize, mSize, minus, mMask, xc}
}

func newSettings(length int, xc xConstants, signed bool) (settings, error) {
	if length > 16 {
		return settings{}, errors.New("maximum length is 16 bits")
	}

	mSize := length - xc.xSize
	if signed {
		mSize -= 1
	}

	if mSize < 1 {
		return settings{}, errors.New("mantissa must be at least 1 bit wide")
	}

	minus := uint16(0)
	if signed {
		minus = uint16(1) << (length - 1)
	}

	mMask := makeBitMask(mSize)
	return makeSettings(mSize, minus, mMask, xc), nil
}

func isNegative(tf uint16, settings *settings) bool {
	return 0b0 != tf&(settings.minus)
}

func abs(tf uint16, settings *settings) uint16 {
	return tf & (^settings.minus)
}

func encode(value float64, settings *settings) uint16 {

	if math.IsNaN(value) {
		return 0x0
	}

	a := settings.xc.twoPowerMinExponent
	reversedB := 1.0 - a

	if value >= 0 {
		return encodeInnerValue(value*reversedB+a, settings)
	} else if 0b0 == settings.minus {
		return 0x0
	} else {
		return settings.minus | encodeInnerValue(-value*reversedB+a, settings)
	}
}

func decode(tf uint16, settings *settings) float64 {

	exponentOffset := -settings.xc.minExponent
	a := settings.xc.twoPowerMinExponent
	reversedB := 1.0 - a
	b := 1.0 / reversedB

	xShift, xMask := settings.xShift, settings.xc.xMask
	mMask, mSize := settings.mMask, settings.mSize

	x := int((tf>>xShift)&xMask) - exponentOffset

	significand := 1.0 + float64(tf&mMask)/powerOfTwo(mSize)
	characteristic := powerOfTwo(x)

	r := significand * characteristic

	return (r - a) * b * sign(tf, settings.minus)
}

func sign(x uint16, minus uint16) float64 {
	if 0b0 == x&minus {
		return 1
	} else {
		return -1
	}
}

func encodeInnerValue(inner float64, s *settings) uint16 {

	xShift, xMask := s.xShift, s.xc.xMask
	mMask, mSize := s.mMask, s.mSize

	exponentOffset := -s.xc.minExponent
	maxValue := (xMask << xShift) | mMask

	twoPowerM := powerOfTwo(mSize)
	internalMaximum := (1 + (twoPowerM-1)/twoPowerM) * s.xc.twoPowerMaxExponent
	if inner >= internalMaximum {
		return maxValue
	}

	x := getExponent(inner, s)
	binaryExponent := uint16(x+exponentOffset) << xShift

	characteristic := powerOfTwo(x)
	normalized := inner / characteristic

	mFloat := (normalized - 1.0) * twoPowerM
	mMax := twoPowerM - 1.0
	binarySignificand := uint16(math.Min(math.Round(mFloat), mMax))

	return binarySignificand | binaryExponent
}

func getExponent(innerValue float64, s *settings) int {
	modulus := math.Abs(innerValue)

	eps := 0.5 / powerOfTwo(s.mSize)
	factor := 2.0 - eps

	for exp := s.xc.maxExponent; exp > s.xc.minExponent; exp-- {
		if factor*powerOfTwo(exp-1) <= modulus {
			return exp
		}
	}

	return s.xc.minExponent
}

func makeBitMask(bits int) uint16 {
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

func min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func encodeDelta(last uint16, x uint16, settings *settings) int {
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

func decodeDelta(last uint16, delta int, s *settings) uint16 {
	if delta == 0 {
		return last
	}

	absLast := int(abs(last, s))

	xShift, xMask := s.xShift, s.xc.xMask
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
