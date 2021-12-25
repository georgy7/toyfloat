package impl

import "math"

type XConstants struct {
	xMask               uint16
	minExponent         int
	maxExponent         int
	twoPowerMinExponent float64
	twoPowerMaxExponent float64
}

type Settings struct {
	mSize  int
	xShift int
	minus  uint16
	mMask  uint16
	xc     XConstants
}

func X4() XConstants {
	return XConstants{
		xMask:               0b1111,
		minExponent:         -8,
		maxExponent:         -8 + 15,
		twoPowerMinExponent: 1.0 / 256.0,
		twoPowerMaxExponent: 128.0,
	}
}

func X3() XConstants {
	return XConstants{
		xMask:               0b111,
		minExponent:         -6,
		maxExponent:         -6 + 7,
		twoPowerMinExponent: 1.0 / 64.0,
		twoPowerMaxExponent: 2.0,
	}
}

func MakeSettings(mSize, xShift int, minus, mMask uint16, xc XConstants) Settings {
	return Settings{mSize, xShift, minus, mMask, xc}
}

func Encode(value float64, settings *Settings) uint16 {

	if math.IsNaN(value) {
		return 0x0
	}

	a := settings.xc.twoPowerMinExponent
	reversedB := 1.0 - a

	if value >= 0 {
		return encode(value*reversedB+a, settings)
	} else if 0b0 == settings.minus {
		return 0x0
	} else {
		return settings.minus | encode(-value*reversedB+a, settings)
	}
}

func Decode(tf uint16, settings *Settings) float64 {

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

func encode(inner float64, s *Settings) uint16 {

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
	binarySignificand := uint16((normalized - 1.0) * twoPowerM)

	return binarySignificand | binaryExponent
}

func getExponent(v float64, s *Settings) int {
	modulus := math.Abs(v)

	for exp := s.xc.maxExponent; exp > s.xc.minExponent; exp-- {
		if powerOfTwo(exp) <= modulus {
			return exp
		}
	}

	return s.xc.minExponent
}

func powerOfTwo(x int) float64 {
	if x < 0 {
		return 1.0 / float64(int(1)<<-x)
	} else {
		return float64(int(1) << x)
	}
}
