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

func getToyfloatPositiveSample() []struct {
	number    float64
	precision float64
} {
	return []struct {
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
}

func TestPrecision(t *testing.T) {
	tests := getToyfloatPositiveSample()

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

func TestIgnoringMostSignificantBits(t *testing.T) {
	for f := -255.0; f <= 255.0; f += 0.01 {
		toy := Encode(f)
		original := Decode(toy)

		if 0xF000&toy != 0x0 {
			t.Fatalf("%.4f -> 0b%b (has extra bits)", f, toy)
		}

		for m := 0x1; m <= 0xF; m++ {
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

// ------------------------

func TestDDZero(t *testing.T) {
	tf := EncodeDD(0)
	t.Logf("Encoded: 0b%b", tf)

	result := DecodeDD(tf)

	if result != 0 {
		t.Fatalf("%f != 0", result)
	}
}

func TestDDPlusOne(t *testing.T) {
	tf := EncodeDD(1)
	t.Logf("Encoded: 0b%b", tf)

	result := DecodeDD(tf)

	if result != 1 {
		t.Fatalf("%f != 1", result)
	}
}

func TestDDMinusOne(t *testing.T) {
	tf := EncodeDD(-1)
	t.Logf("Encoded: 0b%b", tf)

	result := DecodeDD(tf)

	if result != -1 {
		t.Fatalf("%f != -1", result)
	}
}

func TestDDPositiveOverflow(t *testing.T) {
	const expected = 255.99607843137255
	const eps = 0.0001

	for i := 0; i < 1000; i++ {
		v := expected + float64(i)
		tf := EncodeDD(v)

		result := DecodeDD(tf)

		if math.Abs(result-expected) > eps {
			t.Logf("Encoded: 0b%b", tf)
			t.Fatalf("%f != %f (i = %d)", result, expected, i)
		}
	}
}

func TestDDNegativeOverflow(t *testing.T) {
	const expected = -255.99607843137255
	const eps = 0.0001

	for i := 0; i < 1000; i++ {
		v := expected - float64(i)
		tf := EncodeDD(v)

		result := DecodeDD(tf)

		if math.Abs(result-expected) > eps {
			t.Logf("Encoded: 0b%b", tf)
			t.Fatalf("%f != %f (i = %d)", result, expected, i)
		}
	}
}

func TestDDPositiveInfinity(t *testing.T) {
	const expected = 255.99607843137255
	const eps = 0.0001

	v := math.Inf(+1)
	tf := EncodeDD(v)

	result := DecodeDD(tf)

	if math.Abs(result-expected) > eps {
		t.Logf("Encoded: 0b%b", tf)
		t.Fatalf("%f != %f", result, expected)
	}
}

func TestDDNegativeInfinity(t *testing.T) {
	const expected = -255.99607843137255
	const eps = 0.0001

	v := math.Inf(-1)
	tf := EncodeDD(v)

	result := DecodeDD(tf)

	if math.Abs(result-expected) > eps {
		t.Logf("Encoded: 0b%b", tf)
		t.Fatalf("%f != %f", result, expected)
	}
}

func TestDDNaNConvertedToZero(t *testing.T) {
	tf := EncodeDD(math.NaN())
	t.Logf("Encoded: 0b%b", tf)

	result := DecodeDD(tf)

	if result != 0 {
		t.Fatalf("%f != 0", result)
	}
}

func TestDDPrecision(t *testing.T) {
	tests := getToyfloatPositiveSample()

	for _, tt := range tests {
		toy := EncodeDD(tt.number)
		result := DecodeDD(toy)

		diff := math.Abs(result - tt.number)
		if diff > tt.precision {
			t.Fatalf("%.4f -> 0b%b, diff: %f", tt.number, toy, diff)
		}
	}

	for _, tt := range tests {
		negative := -tt.number
		toy := EncodeDD(negative)
		result := DecodeDD(toy)

		diff := math.Abs(result - negative)
		if diff > tt.precision {
			t.Fatalf("%.4f -> 0b%b, diff: %f", negative, toy, diff)
		}
	}
}

func TestDDIgnoringMostSignificantBits(t *testing.T) {
	for f := -255.0; f <= 255.0; f += 0.01 {
		toy := EncodeDD(f)
		original := DecodeDD(toy)

		if 0xF000&toy != 0x0 {
			t.Fatalf("%.4f -> 0b%b (has extra bits)", f, toy)
		}

		for m := 0x1; m <= 0xF; m++ {
			modification := uint16(m) << 12
			toyModified := toy | modification
			modified := DecodeDD(toyModified)

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

// ------------------------

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

// ------------------------

func TestUnsignedPrecision(t *testing.T) {
	tests := getToyfloatPositiveSample()

	for _, tt := range tests {
		toy := EncodeUnsigned(tt.number)
		result := DecodeUnsigned(toy)

		diff := math.Abs(result - tt.number)
		if diff > tt.precision*0.5 {
			t.Fatalf("%.4f -> 0b%b, diff: %f", tt.number, toy, diff)
		}
	}
}

func TestUnsignedNegativeInput(t *testing.T) {
	tests := getToyfloatPositiveSample()

	for _, tt := range tests {
		negative := -tt.number
		toy := EncodeUnsigned(negative)
		result := DecodeUnsigned(toy)

		if result != 0 {
			t.Fatalf("%f != 0", result)
		}
	}
}

func TestUnsignedZero(t *testing.T) {
	tf := EncodeUnsigned(0)
	t.Logf("Encoded: 0b%b", tf)

	result := DecodeUnsigned(tf)

	if result != 0 {
		t.Fatalf("%f != 0", result)
	}
}

func TestUnsignedPlusOne(t *testing.T) {
	tf := EncodeUnsigned(1)
	t.Logf("Encoded: 0b%b", tf)

	result := DecodeUnsigned(tf)

	if result != 1 {
		t.Fatalf("%f != 1", result)
	}
}

func TestUnsignedPositiveOverflow(t *testing.T) {
	const expected = 256.4980392156863
	const eps = 0.0001

	for i := 0; i < 1000; i++ {
		v := expected + float64(i)
		tf := EncodeUnsigned(v)

		result := DecodeUnsigned(tf)

		if math.Abs(result-expected) > eps {
			t.Logf("Encoded: 0b%b", tf)
			t.Fatalf("%f != %f (i = %d)", result, expected, i)
		}
	}
}

func TestUnsignedPositiveInfinity(t *testing.T) {
	const expected = 256.4980392156863
	const eps = 0.0001

	v := math.Inf(+1)
	tf := EncodeUnsigned(v)

	result := DecodeUnsigned(tf)

	if math.Abs(result-expected) > eps {
		t.Logf("Encoded: 0b%b", tf)
		t.Fatalf("%f != %f", result, expected)
	}
}

func TestUnsignedNegativeInfinity(t *testing.T) {
	v := math.Inf(-1)
	tf := EncodeUnsigned(v)

	result := DecodeUnsigned(tf)

	if result != 0 {
		t.Fatalf("%f != 0, encoded: 0b%b", result, tf)
	}
}

func TestUnsignedNaNConvertedToZero(t *testing.T) {
	tf := EncodeUnsigned(math.NaN())
	t.Logf("Encoded: 0b%b", tf)

	result := DecodeUnsigned(tf)

	if result != 0 {
		t.Fatalf("%f != 0", result)
	}
}

func TestUnsignedIgnoringMostSignificantByte(t *testing.T) {
	for f := -255.0; f <= 255.0; f += 0.01 {
		toy := EncodeUnsigned(f)
		original := DecodeUnsigned(toy)

		if 0xF000&toy != 0x0 {
			t.Fatalf("%.4f -> 0b%b (has extra bits)", f, toy)
		}

		for m := 0x1; m < 0xF; m++ {
			modification := uint16(m) << 12
			toyModified := toy | modification
			modified := DecodeUnsigned(toyModified)

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

// ------------------------

func Test13Precision(t *testing.T) {
	tests := getToyfloatPositiveSample()

	for _, tt := range tests {
		toy := Encode13(tt.number)
		result := Decode13(toy)

		diff := math.Abs(result - tt.number)
		if diff > tt.precision*0.5 {
			t.Fatalf("%.4f -> 0b%b, diff: %f", tt.number, toy, diff)
		}
	}

	for _, tt := range tests {
		negative := -tt.number
		toy := Encode13(negative)
		result := Decode13(toy)

		diff := math.Abs(result - negative)
		if diff > tt.precision*0.5 {
			t.Fatalf("%.4f -> 0b%b, diff: %f", negative, toy, diff)
		}
	}
}

func Test13Zero(t *testing.T) {
	tf := Encode13(0)
	t.Logf("Encoded: 0b%b", tf)

	result := Decode13(tf)

	if result != 0 {
		t.Fatalf("%f != 0", result)
	}
}

func Test13PlusOne(t *testing.T) {
	tf := Encode13(1)
	t.Logf("Encoded: 0b%b", tf)

	result := Decode13(tf)

	if result != 1 {
		t.Fatalf("%f != 1", result)
	}
}

func Test13MinusOne(t *testing.T) {
	tf := Encode13(-1)
	t.Logf("Encoded: 0b%b", tf)

	result := Decode13(tf)

	if result != -1 {
		t.Fatalf("%f != -1", result)
	}
}

func Test13PositiveOverflow(t *testing.T) {
	const expected = 256.4980392156863
	const eps = 0.0001

	for i := 0; i < 1000; i++ {
		v := expected + float64(i)
		tf := Encode13(v)

		result := Decode13(tf)

		if math.Abs(result-expected) > eps {
			t.Logf("Encoded: 0b%b", tf)
			t.Fatalf("%f != %f (i = %d)", result, expected, i)
		}
	}
}

func Test13NegativeOverflow(t *testing.T) {
	const expected = -256.4980392156863
	const eps = 0.0001

	for i := 0; i < 1000; i++ {
		v := expected - float64(i)
		tf := Encode13(v)

		result := Decode13(tf)

		if math.Abs(result-expected) > eps {
			t.Logf("Encoded: 0b%b", tf)
			t.Fatalf("%f != %f (i = %d)", result, expected, i)
		}
	}
}

func Test13PositiveInfinity(t *testing.T) {
	const expected = 256.4980392156863
	const eps = 0.0001

	v := math.Inf(+1)
	tf := Encode13(v)

	result := Decode13(tf)

	if math.Abs(result-expected) > eps {
		t.Logf("Encoded: 0b%b", tf)
		t.Fatalf("%f != %f", result, expected)
	}
}

func Test13NegativeInfinity(t *testing.T) {
	const expected = -256.4980392156863
	const eps = 0.0001

	v := math.Inf(-1)
	tf := Encode13(v)

	result := Decode13(tf)

	if math.Abs(result-expected) > eps {
		t.Logf("Encoded: 0b%b", tf)
		t.Fatalf("%f != %f", result, expected)
	}
}

func Test13NaNConvertedToZero(t *testing.T) {
	tf := Encode13(math.NaN())
	t.Logf("Encoded: 0b%b", tf)

	result := Decode13(tf)

	if result != 0 {
		t.Fatalf("%f != 0", result)
	}
}

func Test13IgnoringMostSignificantBits(t *testing.T) {
	for f := -255.0; f <= 255.0; f += 0.01 {
		toy := Encode13(f)
		original := Decode13(toy)

		if 0b1110_0000_0000_0000&toy != 0x0 {
			t.Fatalf("%.4f -> 0b%b (has extra bits)", f, toy)
		}

		for m := 0b1; m <= 0b111; m++ {
			modification := uint16(m) << 13
			toyModified := toy | modification
			modified := Decode13(toyModified)

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

// ------------------------

func Test14Precision(t *testing.T) {
	tests := getToyfloatPositiveSample()

	for _, tt := range tests {
		toy := Encode14(tt.number)
		result := Decode14(toy)

		diff := math.Abs(result - tt.number)
		if diff > tt.precision*0.25 {
			t.Fatalf("%.4f -> 0b%b, diff: %f", tt.number, toy, diff)
		}
	}

	for _, tt := range tests {
		negative := -tt.number
		toy := Encode14(negative)
		result := Decode14(toy)

		diff := math.Abs(result - negative)
		if diff > tt.precision*0.25 {
			t.Fatalf("%.4f -> 0b%b, diff: %f", negative, toy, diff)
		}
	}
}

func Test14Zero(t *testing.T) {
	tf := Encode14(0)
	t.Logf("Encoded: 0b%b", tf)

	result := Decode14(tf)

	if result != 0 {
		t.Fatalf("%f != 0", result)
	}
}

func Test14PlusOne(t *testing.T) {
	tf := Encode14(1)
	t.Logf("Encoded: 0b%b", tf)

	result := Decode14(tf)

	if result != 1 {
		t.Fatalf("%f != 1", result)
	}
}

func Test14MinusOne(t *testing.T) {
	tf := Encode14(-1)
	t.Logf("Encoded: 0b%b", tf)

	result := Decode14(tf)

	if result != -1 {
		t.Fatalf("%f != -1", result)
	}
}

func Test14PositiveOverflow(t *testing.T) {
	const expected = 256.74901960784314
	const eps = 0.0001

	for i := 0; i < 1000; i++ {
		v := expected + float64(i)
		tf := Encode14(v)

		result := Decode14(tf)

		if math.Abs(result-expected) > eps {
			t.Logf("Encoded: 0b%b", tf)
			t.Fatalf("%f != %f (i = %d)", result, expected, i)
		}
	}
}

func Test14NegativeOverflow(t *testing.T) {
	const expected = -256.74901960784314
	const eps = 0.0001

	for i := 0; i < 1000; i++ {
		v := expected - float64(i)
		tf := Encode14(v)

		result := Decode14(tf)

		if math.Abs(result-expected) > eps {
			t.Logf("Encoded: 0b%b", tf)
			t.Fatalf("%f != %f (i = %d)", result, expected, i)
		}
	}
}

func Test14PositiveInfinity(t *testing.T) {
	const expected = 256.74901960784314
	const eps = 0.0001

	v := math.Inf(+1)
	tf := Encode14(v)

	result := Decode14(tf)

	if math.Abs(result-expected) > eps {
		t.Logf("Encoded: 0b%b", tf)
		t.Fatalf("%f != %f", result, expected)
	}
}

func Test14NegativeInfinity(t *testing.T) {
	const expected = -256.74901960784314
	const eps = 0.0001

	v := math.Inf(-1)
	tf := Encode14(v)

	result := Decode14(tf)

	if math.Abs(result-expected) > eps {
		t.Logf("Encoded: 0b%b", tf)
		t.Fatalf("%f != %f", result, expected)
	}
}

func Test14NaNConvertedToZero(t *testing.T) {
	tf := Encode14(math.NaN())
	t.Logf("Encoded: 0b%b", tf)

	result := Decode14(tf)

	if result != 0 {
		t.Fatalf("%f != 0", result)
	}
}

func Test14IgnoringMostSignificantBits(t *testing.T) {
	for f := -255.0; f <= 255.0; f += 0.01 {
		toy := Encode14(f)
		original := Decode14(toy)

		if 0b1100_0000_0000_0000&toy != 0x0 {
			t.Fatalf("%.4f -> 0b%b (has extra bits)", f, toy)
		}

		for m := 0b01; m <= 0b11; m++ {
			modification := uint16(m) << 14
			toyModified := toy | modification
			modified := Decode14(toyModified)

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

// ------------------------

func Test14DPrecision(t *testing.T) {
	tests := getToyfloatPositiveSample()

	for _, tt := range tests {
		toy := Encode14D(tt.number)
		result := Decode14D(toy)

		diff := math.Abs(result - tt.number)
		if diff > tt.precision*0.25 {
			t.Fatalf("%.4f -> 0b%b, diff: %f", tt.number, toy, diff)
		}
	}

	for _, tt := range tests {
		negative := -tt.number
		toy := Encode14D(negative)
		result := Decode14D(toy)

		diff := math.Abs(result - negative)
		if diff > tt.precision*0.25 {
			t.Fatalf("%.4f -> 0b%b, diff: %f", negative, toy, diff)
		}
	}
}

func Test14DIgnoringMostSignificantBits(t *testing.T) {
	for f := -255.0; f <= 255.0; f += 0.01 {
		toy := Encode14D(f)
		original := Decode14D(toy)

		if 0b1100_0000_0000_0000&toy != 0x0 {
			t.Fatalf("%.4f -> 0b%b (has extra bits)", f, toy)
		}

		for m := 0b01; m <= 0b11; m++ {
			modification := uint16(m) << 14
			toyModified := toy | modification
			modified := Decode14D(toyModified)

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

// ------------------------

func TestM11X3Precision(t *testing.T) {
	tests := getToyfloatPositiveSample()

	fixPrecision := func(number, precision float64) float64 {
		if number < 0.0158 {
			return 0.000008
		} else {
			return precision * 0.0625
		}
	}

	for _, tt := range tests {
		if tt.number <= 4.0 {
			toy := EncodeM11X3(tt.number)
			result := DecodeM11X3(toy)

			diff := math.Abs(result - tt.number)
			if diff > fixPrecision(tt.number, tt.precision) {
				t.Fatalf("%.6f -> 0b%b, diff: %f", tt.number, toy, diff)
			}
		}
	}

	for _, tt := range tests {
		if tt.number <= 4.0 {
			negative := -tt.number
			toy := EncodeM11X3(negative)
			result := DecodeM11X3(toy)

			diff := math.Abs(result - negative)
			if diff > fixPrecision(tt.number, tt.precision) {
				t.Fatalf("%.6f -> 0b%b, diff: %f", negative, toy, diff)
			}
		}
	}
}

func TestM11X3Zero(t *testing.T) {
	tf := EncodeM11X3(0)
	t.Logf("Encoded: 0b%b", tf)

	result := DecodeM11X3(tf)

	if result != 0 {
		t.Fatalf("%f != 0", result)
	}
}

func TestM11X3PlusOne(t *testing.T) {
	tf := EncodeM11X3(1)
	t.Logf("Encoded: 0b%b", tf)

	result := DecodeM11X3(tf)

	if result != 1 {
		t.Fatalf("%f != 1", result)
	}
}

func TestM11X3MinusOne(t *testing.T) {
	tf := EncodeM11X3(-1)
	t.Logf("Encoded: 0b%b", tf)

	result := DecodeM11X3(tf)

	if result != -1 {
		t.Fatalf("%f != -1", result)
	}
}

func TestM11X3PositiveOverflow(t *testing.T) {
	const expected = 4.046627
	const eps = 0.0001

	for i := 0; i < 1000; i++ {
		v := expected + float64(i)
		tf := EncodeM11X3(v)

		result := DecodeM11X3(tf)

		if math.Abs(result-expected) > eps {
			t.Logf("Encoded: 0b%b", tf)
			t.Fatalf("%f != %f (i = %d)", result, expected, i)
		}
	}
}

func TestM11X3NegativeOverflow(t *testing.T) {
	const expected = -4.046627
	const eps = 0.0001

	for i := 0; i < 1000; i++ {
		v := expected - float64(i)
		tf := EncodeM11X3(v)

		result := DecodeM11X3(tf)

		if math.Abs(result-expected) > eps {
			t.Logf("Encoded: 0b%b", tf)
			t.Fatalf("%f != %f (i = %d)", result, expected, i)
		}
	}
}

func TestM11X3PositiveInfinity(t *testing.T) {
	const expected = 4.046627
	const eps = 0.0001

	v := math.Inf(+1)
	tf := EncodeM11X3(v)

	result := DecodeM11X3(tf)

	if math.Abs(result-expected) > eps {
		t.Logf("Encoded: 0b%b", tf)
		t.Fatalf("%f != %f", result, expected)
	}
}

func TestM11X3NegativeInfinity(t *testing.T) {
	const expected = -4.046627
	const eps = 0.0001

	v := math.Inf(-1)
	tf := EncodeM11X3(v)

	result := DecodeM11X3(tf)

	if math.Abs(result-expected) > eps {
		t.Logf("Encoded: 0b%b", tf)
		t.Fatalf("%f != %f", result, expected)
	}
}

func TestM11X3NaNConvertedToZero(t *testing.T) {
	tf := EncodeM11X3(math.NaN())
	t.Logf("Encoded: 0b%b", tf)

	result := DecodeM11X3(tf)

	if result != 0 {
		t.Fatalf("%f != 0", result)
	}
}

func TestM11X3IgnoringMostSignificantBits(t *testing.T) {
	for f := -255.0; f <= 255.0; f += 0.01 {
		toy := EncodeM11X3(f)
		original := DecodeM11X3(toy)

		if 0b1000_0000_0000_0000&toy != 0x0 {
			t.Fatalf("%.4f -> 0b%b (has extra bits)", f, toy)
		}

		m := 0b1
		modification := uint16(m) << 15
		toyModified := toy | modification
		modified := DecodeM11X3(toyModified)

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

// ------------------------

func TestM11X3DPrecision(t *testing.T) {
	tests := getToyfloatPositiveSample()

	fixPrecision := func(number, precision float64) float64 {
		if number < 0.0158 {
			return 0.000008
		} else {
			return precision * 0.0625
		}
	}

	for _, tt := range tests {
		if tt.number <= 4.0 {
			toy := EncodeM11X3D(tt.number)
			result := DecodeM11X3D(toy)

			diff := math.Abs(result - tt.number)
			if diff > fixPrecision(tt.number, tt.precision) {
				t.Fatalf("%.6f -> 0b%b, diff: %f", tt.number, toy, diff)
			}
		}
	}

	for _, tt := range tests {
		if tt.number <= 4.0 {
			negative := -tt.number
			toy := EncodeM11X3D(negative)
			result := DecodeM11X3D(toy)

			diff := math.Abs(result - negative)
			if diff > fixPrecision(tt.number, tt.precision) {
				t.Fatalf("%.6f -> 0b%b, diff: %f", negative, toy, diff)
			}
		}
	}
}

func TestM11X3DZero(t *testing.T) {
	tf := EncodeM11X3D(0)
	t.Logf("Encoded: 0b%b", tf)

	result := DecodeM11X3D(tf)

	if result != 0 {
		t.Fatalf("%f != 0", result)
	}
}

func TestM11X3DPlusOne(t *testing.T) {
	tf := EncodeM11X3D(1)
	t.Logf("Encoded: 0b%b", tf)

	result := DecodeM11X3D(tf)

	if result != 1 {
		t.Fatalf("%f != 1", result)
	}
}

func TestM11X3DMinusOne(t *testing.T) {
	tf := EncodeM11X3D(-1)
	t.Logf("Encoded: 0b%b", tf)

	result := DecodeM11X3D(tf)

	if result != -1 {
		t.Fatalf("%f != -1", result)
	}
}

func TestM11X3DPositiveOverflow(t *testing.T) {
	const expected = 4.046627
	const eps = 0.0001

	for i := 0; i < 1000; i++ {
		v := expected + float64(i)
		tf := EncodeM11X3D(v)

		result := DecodeM11X3D(tf)

		if math.Abs(result-expected) > eps {
			t.Logf("Encoded: 0b%b", tf)
			t.Fatalf("%f != %f (i = %d)", result, expected, i)
		}
	}
}

func TestM11X3DNegativeOverflow(t *testing.T) {
	const expected = -4.046627
	const eps = 0.0001

	for i := 0; i < 1000; i++ {
		v := expected - float64(i)
		tf := EncodeM11X3D(v)

		result := DecodeM11X3D(tf)

		if math.Abs(result-expected) > eps {
			t.Logf("Encoded: 0b%b", tf)
			t.Fatalf("%f != %f (i = %d)", result, expected, i)
		}
	}
}

func TestM11X3DPositiveInfinity(t *testing.T) {
	const expected = 4.046627
	const eps = 0.0001

	v := math.Inf(+1)
	tf := EncodeM11X3D(v)

	result := DecodeM11X3D(tf)

	if math.Abs(result-expected) > eps {
		t.Logf("Encoded: 0b%b", tf)
		t.Fatalf("%f != %f", result, expected)
	}
}

func TestM11X3DNegativeInfinity(t *testing.T) {
	const expected = -4.046627
	const eps = 0.0001

	v := math.Inf(-1)
	tf := EncodeM11X3D(v)

	result := DecodeM11X3D(tf)

	if math.Abs(result-expected) > eps {
		t.Logf("Encoded: 0b%b", tf)
		t.Fatalf("%f != %f", result, expected)
	}
}

func TestM11X3DNaNConvertedToZero(t *testing.T) {
	tf := EncodeM11X3D(math.NaN())
	t.Logf("Encoded: 0b%b", tf)

	result := DecodeM11X3D(tf)

	if result != 0 {
		t.Fatalf("%f != 0", result)
	}
}

func TestM11X3DIgnoringMostSignificantBits(t *testing.T) {
	for f := -255.0; f <= 255.0; f += 0.01 {
		toy := EncodeM11X3D(f)
		original := DecodeM11X3D(toy)

		if 0b1000_0000_0000_0000&toy != 0x0 {
			t.Fatalf("%.4f -> 0b%b (has extra bits)", f, toy)
		}

		m := 0b1
		modification := uint16(m) << 15
		toyModified := toy | modification
		modified := DecodeM11X3D(toyModified)

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

// ------------------------

func TestEncodeDecodeStability(t *testing.T) {
	tf := EncodeM11X3D(0.6)
	t.Logf("Encoded: 0b%b", tf)

	input := DecodeM11X3D(tf)
	t.Logf("Decoded: %f", input)

	temp := input

	for i := 0; i < 10; i++ {
		if temp != input {
			t.Fatalf("#%d: %f != %f", i, temp, input)
		}

		tfTemp := EncodeM11X3D(temp)
		t.Logf("Encoded: 0b%b", tfTemp)
		temp = DecodeM11X3D(tfTemp)
		t.Logf("Decoded: %f", temp)
	}
}

func TestUseDeltaM11X3D(t *testing.T) {
	const eps060 = 0.00024   // x = -1
	const eps030 = 0.00012   // x = -2
	const eps013 = 0.00006   // x = -3
	const epsMin = 0.0000075 // x = -6
	const epsMax = 0.00096   // x = 1

	a := math.Pow(2, -6)
	b := 1 / (1 - a)
	twoPowerM := math.Pow(2, 11)

	last := DecodeM11X3D(EncodeM11X3D(0.6))
	lastTf := EncodeM11X3D(last)
	t.Logf("Last encoded: 0b%b", lastTf)

	resultTf := UseIntegerDeltaM11X3D(lastTf, 0)
	result := DecodeM11X3D(resultTf)
	expected := last

	if math.Abs(result-expected) > eps060 {
		t.Logf("delta=0 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = UseIntegerDeltaM11X3D(lastTf, 1)
	result = DecodeM11X3D(resultTf)
	expected = last + ((1.0/twoPowerM)*0.5)*b

	if math.Abs(result-expected) > eps060 {
		t.Logf("delta=1 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = UseIntegerDeltaM11X3D(lastTf, -1)
	result = DecodeM11X3D(resultTf)
	expected = last - ((1.0/twoPowerM)*0.5)*b

	if math.Abs(result-expected) > eps060 {
		t.Logf("delta=-1 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = UseIntegerDeltaM11X3D(lastTf, 123)
	result = DecodeM11X3D(resultTf)
	expected = last + ((123.0/twoPowerM)*0.5)*b

	if math.Abs(result-expected) > eps060 {
		t.Logf("delta=123 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = UseIntegerDeltaM11X3D(lastTf, -123)
	result = DecodeM11X3D(resultTf)
	expected = last - ((123.0/twoPowerM)*0.5)*b

	if math.Abs(result-expected) > eps060 {
		t.Logf("delta=-123 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	last = DecodeM11X3D(EncodeM11X3D(0.3))
	lastTf = EncodeM11X3D(last)
	t.Logf("Last encoded: 0b%b", lastTf)

	mantissa := lastTf & 0b11111111111
	if mantissa != 499 {
		t.Fatalf("This test is probably broken: mantissa equals %d.", mantissa)
	}

	result = DecodeM11X3D(UseIntegerDeltaM11X3D(lastTf, 0))
	expected = last

	if math.Abs(result-expected) > eps030 {
		t.Fatalf("%f != %f", result, expected)
	}

	result = DecodeM11X3D(UseIntegerDeltaM11X3D(lastTf, 2047-499))
	expected = last + (((2047.0-499.0)/twoPowerM)*0.25)*b

	if math.Abs(result-expected) > eps030 {
		t.Fatalf("%f != %f", result, expected)
	}

	result = DecodeM11X3D(UseIntegerDeltaM11X3D(lastTf, -499))
	expected = last - ((499.0/twoPowerM)*0.25)*b

	if math.Abs(result-expected) > eps030 {
		t.Fatalf("%f != %f", result, expected)
	}

	result = DecodeM11X3D(UseIntegerDeltaM11X3D(lastTf, 2047-499+1))
	expected = last +
		(((2047.0-499.0)/twoPowerM)*0.25)*b +
		((1.0/twoPowerM)*0.5)*b

	if math.Abs(result-expected) > eps060 {
		t.Fatalf("%f != %f", result, expected)
	}

	result = DecodeM11X3D(UseIntegerDeltaM11X3D(lastTf, -500))
	expected = last -
		((499.0/twoPowerM)*0.25)*b -
		((1.0/twoPowerM)*0.125)*b

	if math.Abs(result-expected) > eps013 {
		t.Fatalf("%f != %f", result, expected)
	}

	lastTf = 0b0
	last = DecodeM11X3D(lastTf)
	t.Logf("Last encoded: 0b%b", lastTf)

	resultTf = UseIntegerDeltaM11X3D(lastTf, 0)
	result = DecodeM11X3D(resultTf)
	expected = last

	if math.Abs(result-expected) > epsMin {
		t.Logf("delta=0 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = UseIntegerDeltaM11X3D(lastTf, 1)
	result = DecodeM11X3D(resultTf)
	expected = last + ((1.0/twoPowerM)*0.015625)*b

	if math.Abs(result-expected) > epsMin {
		t.Logf("delta=1 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = UseIntegerDeltaM11X3D(lastTf, -1)
	result = DecodeM11X3D(resultTf)
	expected = -last // minus zero

	if math.Abs(result-expected) > epsMin {
		t.Logf("delta=-1 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	// and back
	minusBit := EncodeM11X3D(1) ^ EncodeM11X3D(-1)
	t.Logf("Minus bit: 0b%b", minusBit)
	lastTf = minusBit | 0b0
	last = DecodeM11X3D(lastTf)
	t.Logf("Last encoded: 0b%b", lastTf)

	resultTf = UseIntegerDeltaM11X3D(lastTf, 1)
	result = DecodeM11X3D(resultTf)
	expected = -last // plus zero

	if math.Abs(result-expected) > epsMin {
		t.Logf("delta=1 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = UseIntegerDeltaM11X3D(lastTf, -10)
	result = DecodeM11X3D(resultTf)
	expected = last - ((10.0/twoPowerM)*0.015625)*b

	if math.Abs(result-expected) > epsMin {
		t.Logf("delta=-10 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	// big deltas

	resultTf = UseIntegerDeltaM11X3D(lastTf, 234234)
	result = DecodeM11X3D(resultTf)
	expected = ((1.0+(twoPowerM-1)/twoPowerM)*2 - a) * b

	if math.Abs(result-expected) > epsMax {
		t.Logf("delta=234234 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = UseIntegerDeltaM11X3D(lastTf, -16382)
	result = DecodeM11X3D(resultTf)
	expected = -((1.0+(twoPowerM-2)/twoPowerM)*2 - a) * b

	if math.Abs(result-expected) > epsMax {
		t.Logf("delta=-16382 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = UseIntegerDeltaM11X3D(lastTf, -16383)
	result = DecodeM11X3D(resultTf)
	expected = -((1.0+(twoPowerM-1)/twoPowerM)*2 - a) * b

	if math.Abs(result-expected) > epsMax {
		t.Logf("delta=-16383 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = UseIntegerDeltaM11X3D(lastTf, -16384)
	result = DecodeM11X3D(resultTf)
	expected = -((1.0+(twoPowerM-1)/twoPowerM)*2 - a) * b

	if math.Abs(result-expected) > epsMax {
		t.Logf("delta=-16384 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = UseIntegerDeltaM11X3D(lastTf, -26384)
	result = DecodeM11X3D(resultTf)
	expected = -((1.0+(twoPowerM-1)/twoPowerM)*2 - a) * b

	if math.Abs(result-expected) > epsMax {
		t.Logf("delta=-26384 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}
}

func TestGetDeltaM11X3D(t *testing.T) {
	const start = -4.0
	const stop = 4.0
	const step = 0.01

	const eps = 1e-9

	last := EncodeM11X3D(start)

	for x := start + step; x <= stop; x += step {
		xtf := EncodeM11X3D(x)
		expected := DecodeM11X3D(xtf)

		delta := GetIntegerDeltaM11X3D(last, xtf)
		resultTf := UseIntegerDeltaM11X3D(last, delta)
		result := DecodeM11X3D(resultTf)

		diff := math.Abs(result - expected)
		if diff > eps {
			t.Logf("eps = %f", eps)
			t.Logf("delta = %d", delta)
			t.Logf("last = 0b%b", last)
			t.Logf("this = 0b%b", resultTf)
			t.Fatalf("%f != %f, absolute diff=%f", result, expected, diff)
		}

		last = xtf
	}
}

func TestGetDeltaDD(t *testing.T) {
	const start = -256.0
	const stop = 256.0
	const step = 0.01

	const eps = 1e-9

	last := EncodeDD(start)

	for x := start + step; x <= stop; x += step {
		xtf := EncodeDD(x)
		expected := DecodeDD(xtf)

		delta := GetIntegerDeltaDD(last, xtf)
		resultTf := UseIntegerDeltaDD(last, delta)
		result := DecodeDD(resultTf)

		diff := math.Abs(result - expected)
		if diff > eps {
			t.Logf("eps = %f", eps)
			t.Logf("delta = %d", delta)
			t.Logf("last = 0b%b", last)
			t.Logf("this = 0b%b", resultTf)
			t.Fatalf("%f != %f, absolute diff=%f", result, expected, diff)
		}

		last = xtf
	}
}

func TestGetDelta13(t *testing.T) {
	const start = -256.0
	const stop = 256.0
	const step = 0.01

	const eps = 1e-9

	last := Encode13(start)

	for x := start + step; x <= stop; x += step {
		xtf := Encode13(x)
		expected := Decode13(xtf)

		delta := GetIntegerDelta13(last, xtf)
		resultTf := UseIntegerDelta13(last, delta)
		result := Decode13(resultTf)

		diff := math.Abs(result - expected)
		if diff > eps {
			t.Logf("eps = %f", eps)
			t.Logf("delta = %d", delta)
			t.Logf("last = 0b%b", last)
			t.Logf("this = 0b%b", resultTf)
			t.Fatalf("%f != %f, absolute diff=%f", result, expected, diff)
		}

		last = xtf
	}
}

func TestGetDelta14D(t *testing.T) {
	const start = -256.0
	const stop = 256.0
	const step = 0.01

	const eps = 1e-9

	last := Encode14D(start)

	for x := start + step; x <= stop; x += step {
		xtf := Encode14D(x)
		expected := Decode14D(xtf)

		delta := GetIntegerDelta14D(last, xtf)
		resultTf := UseIntegerDelta14D(last, delta)
		result := Decode14D(resultTf)

		diff := math.Abs(result - expected)
		if diff > eps {
			t.Logf("eps = %f", eps)
			t.Logf("delta = %d", delta)
			t.Logf("last = 0b%b", last)
			t.Logf("this = 0b%b", resultTf)
			t.Fatalf("%f != %f, absolute diff=%f", result, expected, diff)
		}

		last = xtf
	}
}

// ------------------------

func TestMinusBitPosition(t *testing.T) {
	tf := Encode(42)
	t.Logf("Encoded: 0b%b", tf)

	a := Decode(tf)
	b := -Decode(tf | 0b1000_0000)

	if a != b {
		t.Fatalf("%f != %f", a, b)
	}
}

func TestDDMinusBitPosition(t *testing.T) {
	tf := EncodeDD(42)
	t.Logf("Encoded: 0b%b", tf)

	a := DecodeDD(tf)
	b := -DecodeDD(tf | 0b1000_0000_0000)

	if a != b {
		t.Fatalf("%f != %f", a, b)
	}
}

func Test13MinusBitPosition(t *testing.T) {
	tf := Encode13(42)
	t.Logf("Encoded: 0b%b", tf)

	a := Decode13(tf)
	b := -Decode13(tf | 0b1_0000_0000_0000)

	if a != b {
		t.Fatalf("%f != %f", a, b)
	}
}

func Test14MinusBitPosition(t *testing.T) {
	tf := Encode14(42)
	t.Logf("Encoded: 0b%b", tf)

	a := Decode14(tf)
	b := -Decode14(tf | 0b10_0000_0000)

	if a != b {
		t.Fatalf("%f != %f", a, b)
	}
}

func Test14DMinusBitPosition(t *testing.T) {
	tf := Encode14D(42)
	t.Logf("Encoded: 0b%b", tf)

	a := Decode14D(tf)
	b := -Decode14D(tf | 0b10_0000_0000_0000)

	if a != b {
		t.Fatalf("%f != %f", a, b)
	}
}

func TestM11X3MinusBitPosition(t *testing.T) {
	tf := EncodeM11X3(42)
	t.Logf("Encoded: 0b%b", tf)

	a := DecodeM11X3(tf)
	b := -DecodeM11X3(tf | 0b1000_0000_0000)

	if a != b {
		t.Fatalf("%f != %f", a, b)
	}
}

func TestM11X3DMinusBitPosition(t *testing.T) {
	tf := EncodeM11X3D(42)
	t.Logf("Encoded: 0b%b", tf)

	a := DecodeM11X3D(tf)
	b := -DecodeM11X3D(tf | 0b0100_0000_0000_0000)

	if a != b {
		t.Fatalf("%f != %f", a, b)
	}
}

func TestReadme(t *testing.T) {
	const input = 0.345
	const eps = 1e-6

	{
		tf := EncodeDD(input)
		if tf != 0x332 {
			t.Fatalf("Incorrect encoded: 0x%X (defaultD)\n", tf)
		}

		result := DecodeDD(tf)
		if math.Abs(result-0.345098) > eps {
			t.Fatalf("Incorrect decoded: %f (defaultD)\n", result)
		}
	}

	{
		tf := EncodeUnsigned(input)
		if tf != 0x664 {
			t.Fatalf("Incorrect encoded: 0x%X (unsigned)\n", tf)
		}
	}

	{
		tf := Encode13(input)
		if tf != 0x664 {
			t.Fatalf("Incorrect encoded: 0x%X (13-bit)\n", tf)
		}

		result := Decode13(tf)
		if math.Abs(result-0.345098) > eps {
			t.Fatalf("Incorrect decoded: %f (13-bit)\n", result)
		}
	}

	{
		tf := Encode14D(input)
		if tf != 0xCC8 {
			t.Fatalf("Incorrect encoded: 0x%X (14d)\n", tf)
		}

		result := Decode14D(tf)
		if math.Abs(result-0.345098) > eps {
			t.Fatalf("Incorrect decoded: %f (14d)\n", result)
		}
	}

	{
		tf := EncodeM11X3D(input)
		if tf != 0x235E {
			t.Fatalf("Incorrect encoded: 0x%X (m11x3d)\n", tf)
		}

		result := DecodeM11X3D(tf)
		if math.Abs(result-0.344990) > eps {
			t.Fatalf("Incorrect decoded: %f (m11x3d)\n", result)
		}
	}

	{
		tf := EncodeM11X3D(input)
		result := DecodeM11X3D(tf)
		if math.Abs(result-0.344990) > eps {
			t.Fatalf("Incorrect decoded: %f (m11x3d)\n", result)
		}
	}

	{
		series := []float64{-0.0058, 0.01, 0.123, 0.134, 0.132, 0.144, 0.145, 0.140}
		expected := []int{387, 414, 12, -2, 12, 1, -5}

		previous := EncodeDD(series[0])
		for i := 1; i < len(series); i++ {
			this := EncodeDD(series[i])

			delta := GetIntegerDeltaDD(previous, this)
			expectedDelta := expected[i-1]

			if delta != expectedDelta {
				t.Fatalf("%d != %d\n", delta, expectedDelta)
			}

			previous = this
		}

	}
}

func TestExtremeCases(t *testing.T) {
	{
		input := 0.9999999999995131
		result := Decode13(Encode13(input))

		if math.Abs(result-input) > 0.001 {
			t.Fatalf("%f != %f (13-bit)\n", result, input)
		}
	}
}
