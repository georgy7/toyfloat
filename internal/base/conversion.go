package base

import "math"

const exponentOffset = 7

const minExponent = -exponentOffset
const maxExponent = minExponent + 15

func getExponent(v float64) int {
	modulus := math.Abs(v)

	for exp := maxExponent; exp > minExponent; exp-- {
		if math.Pow(2, float64(exp)) <= modulus {
			return exp
		}
	}

	return minExponent
}

func GetExponentAsANibble(v float64) uint8 {
	if math.IsNaN(v) {
		return 0
	} else {
		return uint8(getExponent(v) + exponentOffset)
	}
}

func Read(head uint8, exponentNibble uint8) float64 {
	exponent := float64(exponentNibble) - exponentOffset

	significand := 1.0 + float64(head&maxSignificand)/128.0
	characteristic := math.Pow(2, exponent)

	result := significand * characteristic

	if minus == head&minus {
		return -result
	} else {
		return result
	}
}

const minus uint8 = 0b1000_0000

const maxSignificand uint8 = 0b0111_1111

// Maximum exponent equals 8 (minimum equals -7).
const maxPositive float64 = (1.0 + 127.0/128.0) * 256.0

func GetSignificand(v float64) uint8 {
	if math.IsNaN(v) {
		// You should filter such values yourself.
		return 0x0
	} else if v >= maxPositive {
		return maxSignificand
	} else if v <= (-maxPositive) {
		return minus | maxSignificand
	}

	characteristic := math.Pow(2, float64(getExponent(v)))
	significand := v / characteristic

	if significand >= 0 {
		return getNumerator(significand)
	} else {
		return minus | getNumerator(-significand)
	}
}

func getNumerator(a float64) uint8 {
	return uint8((a - 1.0) * 128.0)
}
