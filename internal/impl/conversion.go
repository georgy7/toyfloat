package impl

import "math"

const exponentOffset = 8
const a = 1.0 / 256.0

const reversedB = 1.0 - a
const b = 1.0 / reversedB

const minExponent = -exponentOffset
const maxExponent = minExponent + 15

const xMask = 0b1111

func Encode(value float64, mSize int, xShift int,
	minus uint16, mMask uint16) uint16 {

	if math.IsNaN(value) {
		return 0x0
	}

	binaryMaxValue := (xMask << xShift) | mMask

	if value >= 0 {
		return encode(value*reversedB+a, mSize, binaryMaxValue)
	} else {
		return minus | encode(-value*reversedB+a, mSize, binaryMaxValue)
	}
}

func Decode(tf uint16, mSize int, xShift int,
	minus uint16, mMask uint16) float64 {

	x := int((tf>>xShift)&xMask) - exponentOffset

	significand := 1.0 + float64(tf&mMask)/powerOfTwo(mSize)
	characteristic := powerOfTwo(x)

	r := significand * characteristic

	return (r - a) * b * sign(tf, minus)
}

func sign(x uint16, minus uint16) float64 {
	if minus == x&minus {
		return -1
	} else {
		return 1
	}
}

func encode(inner float64, mSize int, maxValue uint16) uint16 {
	internalMaximum := powerOfTwo(mSize+1) - 1
	if inner >= internalMaximum {
		return maxValue
	}

	x := getExponent(inner)
	binaryExponent := uint16(x+exponentOffset) << 8

	characteristic := powerOfTwo(x)
	normalized := inner / characteristic
	binarySignificand := uint16((normalized - 1.0) * powerOfTwo(mSize))

	return binarySignificand | binaryExponent
}

func getExponent(v float64) int {
	modulus := math.Abs(v)

	for exp := maxExponent; exp > minExponent; exp-- {
		if powerOfTwo(exp) <= modulus {
			return exp
		}
	}

	return minExponent
}

func powerOfTwo(x int) float64 {
	if x < 0 {
		return 1.0 / float64(int(1)<<-x)
	} else {
		return float64(int(1) << x)
	}
}
