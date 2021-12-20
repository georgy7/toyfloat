package impl

import "math"

const exponentOffset = 8

const minExponent = -exponentOffset
const maxExponent = minExponent + 15

const a = 1.0 / 256.0
const reversedB = 1.0 - a
const b = 1.0 / reversedB

const minus uint16 = 0b1000_0000
const binaryMaxSignificand uint16 = 0b0111_1111
const binaryMaxExponent uint16 = 0b1111_0000_0000
const binaryMaxValue = binaryMaxExponent | binaryMaxSignificand

func Encode(v float64) uint16 {
	if math.IsNaN(v) {
		return 0x0
	} else if v >= 0 {
		return encode(v*reversedB + a)
	} else {
		return minus | encode(-v*reversedB+a)
	}
}

func Decode(x uint16) float64 {
	return (decode(x) - a) * b * sign(x)
}

func sign(x uint16) float64 {
	if minus == x&minus {
		return -1
	} else {
		return 1
	}
}

func encode(inner float64) uint16 {
	const internalMaximum float64 = 255.0
	if inner >= internalMaximum {
		return binaryMaxValue
	} else if inner <= (-internalMaximum) {
		return binaryMaxValue | minus
	}

	x := getExponent(inner)
	binaryExponent := uint16(x+exponentOffset) << 8

	characteristic := math.Pow(2, float64(x))
	normalized := inner / characteristic

	if normalized >= 0 {
		return toBinarySignificand(normalized) | binaryExponent
	} else {
		return minus | toBinarySignificand(-normalized) | binaryExponent
	}
}

func toBinarySignificand(normalizedSignificand float64) uint16 {
	return uint16((normalizedSignificand - 1.0) * 128.0)
}

func decode(tf uint16) float64 {
	x := float64((tf&binaryMaxExponent)>>8) - exponentOffset

	significand := 1.0 + float64(tf&binaryMaxSignificand)/128.0
	characteristic := math.Pow(2, x)

	return significand * characteristic
}

func getExponent(v float64) int {
	modulus := math.Abs(v)

	for exp := maxExponent; exp > minExponent; exp-- {
		if math.Pow(2, float64(exp)) <= modulus {
			return exp
		}
	}

	return minExponent
}
