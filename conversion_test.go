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

func TestMinusBitPosition(t *testing.T) {
	tf := Encode(42)
	t.Logf("Encoded: 0b%b", tf)

	a := Decode(tf)
	b := -Decode(tf | 0b1000_0000)

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

func TestM11X3MinusBitPosition(t *testing.T) {
	tf := EncodeM11X3(42)
	t.Logf("Encoded: 0b%b", tf)

	a := DecodeM11X3(tf)
	b := -DecodeM11X3(tf | 0b1000_0000_0000)

	if a != b {
		t.Fatalf("%f != %f", a, b)
	}
}

func TestReadme(t *testing.T) {
	const input = 0.345
	const eps = 1e-6

	{
		tf := Encode(input)
		if tf != 0x631 {
			t.Fatalf("Incorrect encoded: 0x%X (default)\n", tf)
		}

		result := Decode(tf)
		if math.Abs(result-0.343137) > eps {
			t.Fatalf("Incorrect decoded: %f (default)\n", result)
		}
	}

	{
		tf := EncodeUnsigned(input)
		if tf != 0x663 {
			t.Fatalf("Incorrect encoded: 0x%X (unsigned)\n", tf)
		}
	}

	{
		tf := Encode13(input)
		if tf != 0x663 {
			t.Fatalf("Incorrect encoded: 0x%X (13-bit)\n", tf)
		}

		result := Decode13(tf)
		if math.Abs(result-0.344118) > eps {
			t.Fatalf("Incorrect decoded: %f (13-bit)\n", result)
		}
	}

	{
		tf := Encode14(input)
		if tf != 0x18C7 {
			t.Fatalf("Incorrect encoded: 0x%X (14-bit)\n", tf)
		}

		result := Decode14(tf)
		if math.Abs(result-0.344608) > eps {
			t.Fatalf("Incorrect decoded: %f (14-bit)\n", result)
		}
	}

	{
		tf := EncodeM11X3(input)
		if tf != 0x435E {
			t.Fatalf("Incorrect encoded: 0x%X (m11x3)\n", tf)
		}

		result := DecodeM11X3(tf)
		if math.Abs(result-0.34499) > eps {
			t.Fatalf("Incorrect decoded: %f (m11x3)\n", result)
		}
	}
}
