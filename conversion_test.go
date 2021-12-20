package toyfloat

import (
	"math"
	"testing"
)

func TestZero(t *testing.T) {
	tf := Encode(0)
	t.Logf("Encoded: 0b%b", tf)

	result := Decode(tf)

	if result != 0 {
		t.Fatalf("%f != 0", result)
	}
}

func TestPlusOne(t *testing.T) {
	tf := Encode(1)
	t.Logf("Encoded: 0b%b", tf)

	result := Decode(tf)

	if result != 1 {
		t.Fatalf("%f != 1", result)
	}
}

func TestMinusOne(t *testing.T) {
	tf := Encode(-1)
	t.Logf("Encoded: 0b%b", tf)

	result := Decode(tf)

	if result != -1 {
		t.Fatalf("%f != -1", result)
	}
}

func TestPositiveOverflow(t *testing.T) {
	const expected = 255.99607843137255
	const eps = 0.0001

	for i := 0; i < 1000; i++ {
		v := expected + float64(i)
		tf := Encode(v)

		result := Decode(tf)

		if math.Abs(result-expected) > eps {
			t.Logf("Encoded: 0b%b", tf)
			t.Fatalf("%f != %f (i = %d)", result, expected, i)
		}
	}
}

func TestNegativeOverflow(t *testing.T) {
	const expected = -255.99607843137255
	const eps = 0.0001

	for i := 0; i < 1000; i++ {
		v := expected - float64(i)
		tf := Encode(v)

		result := Decode(tf)

		if math.Abs(result-expected) > eps {
			t.Logf("Encoded: 0b%b", tf)
			t.Fatalf("%f != %f (i = %d)", result, expected, i)
		}
	}
}

func TestNaNConvertedToZero(t *testing.T) {
	tf := Encode(math.NaN())
	t.Logf("Encoded: 0b%b", tf)

	result := Decode(tf)

	if result != 0 {
		t.Fatalf("%f != 0", result)
	}
}

func TestPrecision(t *testing.T) {
	tests := []struct {
		number    float64
		precision float64
	}{
		{0.172, 0.000981},
		{0.345, 0.00197},
		{0.654, 0.00393},
		{0.898, 0.00393},
		{1.015, 0.00785},
		{2.788, 0.01569},
		{41.15, 0.25099},
		{164.6, 1.00393},
	}

	for _, tt := range tests {
		toy := Encode(tt.number)
		result := Decode(toy)

		diff := math.Abs(result - tt.number)
		if diff > tt.precision {
			t.Fatalf("%.4f -> 0b%b, diff: %f", tt.number, toy, diff)
		}
	}
}
