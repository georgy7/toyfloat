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

func TestPositiveInfinity(t *testing.T) {
	const expected = 255.99607843137255
	const eps = 0.0001

	v := math.Inf(+1)
	tf := Encode(v)

	result := Decode(tf)

	if math.Abs(result-expected) > eps {
		t.Logf("Encoded: 0b%b", tf)
		t.Fatalf("%f != %f", result, expected)
	}
}

func TestNegativeInfinity(t *testing.T) {
	const expected = -255.99607843137255
	const eps = 0.0001

	v := math.Inf(-1)
	tf := Encode(v)

	result := Decode(tf)

	if math.Abs(result-expected) > eps {
		t.Logf("Encoded: 0b%b", tf)
		t.Fatalf("%f != %f", result, expected)
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
		{0.000015, 0.000031},
		{0.000055, 0.000031},
		{0.000565, 0.000031},
		{0.000665, 0.000031},
		{0.000765, 0.000031},
		{0.000865, 0.000031},
		{0.000965, 0.000031},
		{0.001165, 0.000031},
		{0.001265, 0.000031},
		{0.001365, 0.000031},
		{0.001465, 0.000031},
		{0.002065, 0.000031},
		{0.002165, 0.000031},
		{0.003065, 0.000031},
		{0.003165, 0.000031},

		{0.009621, 0.000062},
		{0.010107, 0.000062},

		{0.177964, 0.000981},
		{0.342423, 0.00197},
		{0.659234, 0.00393},
		{0.898094, 0.00393},
		{1.015633, 0.00785},
		{2.788122, 0.01569},

		{41.15423, 0.25099},
		{164.6678, 1.00393},
	}

	for _, tt := range tests {
		toy := Encode(tt.number)
		result := Decode(toy)

		diff := math.Abs(result - tt.number)
		if diff > tt.precision {
			t.Fatalf("%.4f -> 0b%b, diff: %f", tt.number, toy, diff)
		}
	}

	for _, tt := range tests {
		negative := -tt.number
		toy := Encode(negative)
		result := Decode(toy)

		diff := math.Abs(result - negative)
		if diff > tt.precision {
			t.Fatalf("%.4f -> 0b%b, diff: %f", negative, toy, diff)
		}
	}
}

func TestIgnoringMostSignificantByte(t *testing.T) {
	for f := -255.0; f <= 255.0; f += 0.01 {
		toy := Encode(f)
		original := Decode(toy)

		if 0xF000&toy != 0x0 {
			t.Fatalf("%.4f -> 0b%b (has extra bits)", f, toy)
		}

		for m := 0x1; m < 0xF; m++ {
			modification := uint16(m) << 12
			toyModified := toy | modification
			modified := Decode(toyModified)

			if toy == toyModified {
				t.Fatalf("This test is broken. "+
					"Toy: 0b%b. Modification: 0b%b.",
					toy, modification)
			}

			if modified != original {
				t.Fatalf("%.4f != %.4f, modification: 0b%b",
					modified, original, modification)
			}
		}
	}
}

func BenchmarkFloat64Increment(b *testing.B) {
	for i := 0; i < b.N; i++ {
		one := 1.0
		counter := 0.0
		for x := 0; x < 100000; x++ {
			counter += one
		}
	}
}

func BenchmarkDecodeEncodeIncrement(b *testing.B) {
	for i := 0; i < b.N; i++ {
		one := 1.0
		counter := Encode(0.0)
		for x := 0; x < 100000; x++ {
			counter = Encode(Decode(counter) + one)
		}
	}
}
