package base

import (
	"fmt"
	"math"
	"testing"
)

type testCase struct {
	arg  float64
	want uint8
}

func TestExponentPositive(t *testing.T) {
	tests := getExponents()
	for ti, tt := range tests {
		a := tt.arg
		testName := fmt.Sprintf("exponent of positive #%d (%f)", ti, a)
		t.Run(testName, func(t *testing.T) {
			got := GetExponentAsANibble(a)
			if got != tt.want {
				t.Fatalf("0b%b != 0b%b", got, tt.want)
			}
		})
	}
}

func TestExponentNegative(t *testing.T) {
	tests := getExponents()
	for ti, tt := range tests {
		a := -tt.arg
		testName := fmt.Sprintf("exponent of negative #%d (%f)", ti, a)
		t.Run(testName, func(t *testing.T) {
			got := GetExponentAsANibble(a)
			if got != tt.want {
				t.Fatalf("0b%b != 0b%b", got, tt.want)
			}
		})
	}
}

func TestSignificandPositive(t *testing.T) {
	tests := getSignificands()
	for ti, tt := range tests {
		a := tt.arg
		testName := fmt.Sprintf("significand of positive #%d (%f)", ti, a)
		t.Run(testName, func(t *testing.T) {
			got := GetSignificand(a)
			if got != tt.want {
				t.Fatalf("0b%b != 0b%b", got, tt.want)
			}
		})
	}
}

func TestSignificandNegative(t *testing.T) {
	tests := getSignificands()
	for ti, tt := range tests {
		a := -tt.arg
		testName := fmt.Sprintf("significand of negative #%d (%f)", ti, a)
		t.Run(testName, func(t *testing.T) {
			got := GetSignificand(a)

			var want uint8
			if math.IsNaN(a) {
				want = tt.want
			} else {
				want = minus | tt.want
			}

			if got != want {
				t.Fatalf("0b%b != 0b%b", got, want)
			}
		})
	}
}

func getExponents() []testCase {
	return []testCase{
		{math.Inf(1), 0b0000_1111},

		{(1.0 + 127.0/128.0) * math.Pow(2, 8), 0b0000_1111},
		{(1.0 + 0.0/128.0) * math.Pow(2, 8), 0b0000_1111},
		{(1.0 + 12.0/128.0) * math.Pow(2, 7), 0b0000_1110},
		{(1.0 + 0.0/128.0) * math.Pow(2, 7), 0b0000_1110},
		{(1.0 + 56.0/128.0) * math.Pow(2, 6), 0b0000_1101},
		{(1.0 + 0.0/128.0) * math.Pow(2, 6), 0b0000_1101},
		{(1.0 + 127.0/128.0) * math.Pow(2, 5), 0b0000_1100},
		{(1.0 + 0.0/128.0) * math.Pow(2, 5), 0b0000_1100},

		{(1.0 + 127.0/128.0) * math.Pow(2, 4), 0b0000_1011},
		{(1.0 + 0.0/128.0) * math.Pow(2, 4), 0b0000_1011},
		{(1.0 + 88.0/128.0) * math.Pow(2, 3), 0b0000_1010},
		{(1.0 + 0.0/128.0) * math.Pow(2, 3), 0b0000_1010},
		{(1.0 + 45.0/128.0) * math.Pow(2, 2), 0b0000_1001},
		{(1.0 + 0.0/128.0) * math.Pow(2, 2), 0b0000_1001},
		{(1.0 + 127.0/128.0) * math.Pow(2, 1), 0b0000_1000},
		{(1.0 + 0.0/128.0) * math.Pow(2, 1), 0b0000_1000},

		{(1.0 + 127.0/128.0) * math.Pow(2, 0), 0b0000_0111},
		{(1.0 + 0.0/128.0) * math.Pow(2, 0), 0b0000_0111},
		{(1.0 + 127.0/128.0) * math.Pow(2, -1), 0b0000_0110},
		{(1.0 + 0.0/128.0) * math.Pow(2, -1), 0b0000_0110},
		{(1.0 + 23.0/128.0) * math.Pow(2, -2), 0b0000_0101},
		{(1.0 + 0.0/128.0) * math.Pow(2, -2), 0b0000_0101},
		{(1.0 + 12.0/128.0) * math.Pow(2, -3), 0b0000_0100},
		{(1.0 + 0.0/128.0) * math.Pow(2, -3), 0b0000_0100},

		{(1.0 + 100.0/128.0) * math.Pow(2, -4), 0b0000_0011},
		{(1.0 + 0.0/128.0) * math.Pow(2, -4), 0b0000_0011},
		{(1.0 + 127.0/128.0) * math.Pow(2, -5), 0b0000_0010},
		{(1.0 + 0.0/128.0) * math.Pow(2, -5), 0b0000_0010},
		{(1.0 + 127.0/128.0) * math.Pow(2, -6), 0b0000_0001},
		{(1.0 + 0.0/128.0) * math.Pow(2, -6), 0b0000_0001},
		{(1.0 + 127.0/128.0) * math.Pow(2, -7), 0b0000_0000},
		{(1.0 + 0.0/128.0) * math.Pow(2, -7), 0b0000_0000},

		{math.NaN(), 0b0000_0000},
	}
}

func getSignificands() []testCase {
	return []testCase{
		{math.Inf(1), 0b0111_1111},

		{(1.0 + 127.0/128.0) * math.Pow(2, 8), 0b0111_1111},
		{(1.0 + 0.0/128.0) * math.Pow(2, 8), 0b0000_0000},
		{(1.0 + 12.0/128.0) * math.Pow(2, 7), 12},
		{(1.0 + 0.0/128.0) * math.Pow(2, 7), 0},
		{(1.0 + 56.0/128.0) * math.Pow(2, 6), 56},
		{(1.0 + 0.0/128.0) * math.Pow(2, 6), 0},
		{(1.0 + 127.0/128.0) * math.Pow(2, 5), 127},
		{(1.0 + 0.0/128.0) * math.Pow(2, 5), 0},

		{(1.0 + 127.0/128.0) * math.Pow(2, 4), 127},
		{(1.0 + 0.0/128.0) * math.Pow(2, 4), 0},
		{(1.0 + 88.0/128.0) * math.Pow(2, 3), 88},
		{(1.0 + 0.0/128.0) * math.Pow(2, 3), 0},
		{(1.0 + 45.0/128.0) * math.Pow(2, 2), 45},
		{(1.0 + 0.0/128.0) * math.Pow(2, 2), 0},
		{(1.0 + 127.0/128.0) * math.Pow(2, 1), 127},
		{(1.0 + 0.0/128.0) * math.Pow(2, 1), 0},

		{(1.0 + 127.0/128.0) * math.Pow(2, 0), 127},
		{(1.0 + 0.0/128.0) * math.Pow(2, 0), 0},
		{(1.0 + 127.0/128.0) * math.Pow(2, -1), 127},
		{(1.0 + 0.0/128.0) * math.Pow(2, -1), 0},
		{(1.0 + 23.0/128.0) * math.Pow(2, -2), 23},
		{(1.0 + 0.0/128.0) * math.Pow(2, -2), 0},
		{(1.0 + 12.0/128.0) * math.Pow(2, -3), 12},
		{(1.0 + 0.0/128.0) * math.Pow(2, -3), 0},

		{(1.0 + 100.0/128.0) * math.Pow(2, -4), 100},
		{(1.0 + 0.0/128.0) * math.Pow(2, -4), 0},
		{(1.0 + 127.0/128.0) * math.Pow(2, -5), 127},
		{(1.0 + 0.0/128.0) * math.Pow(2, -5), 0},
		{(1.0 + 127.0/128.0) * math.Pow(2, -6), 127},
		{(1.0 + 126.0/128.0) * math.Pow(2, -6), 126},
		{(1.0 + 125.0/128.0) * math.Pow(2, -6), 125},
		{(1.0 + 10.0/128.0) * math.Pow(2, -6), 10},
		{(1.0 + 9.0/128.0) * math.Pow(2, -6), 9},
		{(1.0 + 8.0/128.0) * math.Pow(2, -6), 8},
		{(1.0 + 7.0/128.0) * math.Pow(2, -6), 7},
		{(1.0 + 0.0/128.0) * math.Pow(2, -6), 0},
		{(1.0 + 127.0/128.0) * math.Pow(2, -7), 127},
		{(1.0 + 0.0/128.0) * math.Pow(2, -7), 0},

		{math.NaN(), 0b0000_0000},
	}
}
