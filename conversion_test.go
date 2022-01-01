package toyfloat

import (
	"math"
	"testing"
)

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

// ------------------------

func Test12Zero(t *testing.T) {
	tf := Encode12(0)
	t.Logf("Encoded: 0b%b", tf)

	result := Decode12(tf)

	if result != 0 {
		t.Fatalf("%f != 0", result)
	}
}

func Test12PlusOne(t *testing.T) {
	tf := Encode12(1)
	t.Logf("Encoded: 0b%b", tf)

	result := Decode12(tf)

	if result != 1 {
		t.Fatalf("%f != 1", result)
	}
}

func Test12MinusOne(t *testing.T) {
	tf := Encode12(-1)
	t.Logf("Encoded: 0b%b", tf)

	result := Decode12(tf)

	if result != -1 {
		t.Fatalf("%f != -1", result)
	}
}

func Test12PositiveOverflow(t *testing.T) {
	const expected = 255.99607843137255
	const eps = 0.0001

	for i := 0; i < 1000; i++ {
		v := expected + float64(i)
		tf := Encode12(v)

		result := Decode12(tf)

		if math.Abs(result-expected) > eps {
			t.Logf("Encoded: 0b%b", tf)
			t.Fatalf("%f != %f (i = %d)", result, expected, i)
		}
	}
}

func Test12NegativeOverflow(t *testing.T) {
	const expected = -255.99607843137255
	const eps = 0.0001

	for i := 0; i < 1000; i++ {
		v := expected - float64(i)
		tf := Encode12(v)

		result := Decode12(tf)

		if math.Abs(result-expected) > eps {
			t.Logf("Encoded: 0b%b", tf)
			t.Fatalf("%f != %f (i = %d)", result, expected, i)
		}
	}
}

func Test12PositiveInfinity(t *testing.T) {
	const expected = 255.99607843137255
	const eps = 0.0001

	v := math.Inf(+1)
	tf := Encode12(v)

	result := Decode12(tf)

	if math.Abs(result-expected) > eps {
		t.Logf("Encoded: 0b%b", tf)
		t.Fatalf("%f != %f", result, expected)
	}
}

func Test12NegativeInfinity(t *testing.T) {
	const expected = -255.99607843137255
	const eps = 0.0001

	v := math.Inf(-1)
	tf := Encode12(v)

	result := Decode12(tf)

	if math.Abs(result-expected) > eps {
		t.Logf("Encoded: 0b%b", tf)
		t.Fatalf("%f != %f", result, expected)
	}
}

func Test12NaNConvertedToZero(t *testing.T) {
	tf := Encode12(math.NaN())
	t.Logf("Encoded: 0b%b", tf)

	result := Decode12(tf)

	if result != 0 {
		t.Fatalf("%f != 0", result)
	}
}

func Test12Precision(t *testing.T) {
	tests := getToyfloatPositiveSample()

	for _, tt := range tests {
		toy := Encode12(tt.number)
		result := Decode12(toy)

		diff := math.Abs(result - tt.number)
		if diff > tt.precision {
			t.Fatalf("%.4f -> 0b%b, diff: %f", tt.number, toy, diff)
		}
	}

	for _, tt := range tests {
		negative := -tt.number
		toy := Encode12(negative)
		result := Decode12(toy)

		diff := math.Abs(result - negative)
		if diff > tt.precision {
			t.Fatalf("%.4f -> 0b%b, diff: %f", negative, toy, diff)
		}
	}
}

func Test12IgnoringMostSignificantBits(t *testing.T) {
	for f := -255.0; f <= 255.0; f += 0.01 {
		toy := Encode12(f)
		original := Decode12(toy)

		if 0xF000&toy != 0x0 {
			t.Fatalf("%.4f -> 0b%b (has extra bits)", f, toy)
		}

		for m := 0x1; m <= 0xF; m++ {
			modification := uint16(m) << 12
			toyModified := toy | modification
			modified := Decode12(toyModified)

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
		counter := Encode12(0.0)
		for x := 0; x < 100000; x++ {
			counter = Encode12(Decode12(counter) + one)
		}
	}
}

// ------------------------

func TestUnsigned12Precision(t *testing.T) {
	tests := getToyfloatPositiveSample()

	for _, tt := range tests {
		toy := Encode12U(tt.number)
		result := Decode12U(toy)

		diff := math.Abs(result - tt.number)
		if diff > tt.precision*0.5 {
			t.Fatalf("%.4f -> 0b%b, diff: %f", tt.number, toy, diff)
		}
	}
}

func TestUnsigned12NegativeInput(t *testing.T) {
	tests := getToyfloatPositiveSample()

	for _, tt := range tests {
		negative := -tt.number
		toy := Encode12U(negative)
		result := Decode12U(toy)

		if result != 0 {
			t.Fatalf("%f != 0", result)
		}
	}
}

func TestUnsigned12Zero(t *testing.T) {
	tf := Encode12U(0)
	t.Logf("Encoded: 0b%b", tf)

	result := Decode12U(tf)

	if result != 0 {
		t.Fatalf("%f != 0", result)
	}
}

func TestUnsigned12PlusOne(t *testing.T) {
	tf := Encode12U(1)
	t.Logf("Encoded: 0b%b", tf)

	result := Decode12U(tf)

	if result != 1 {
		t.Fatalf("%f != 1", result)
	}
}

func TestUnsigned12PositiveOverflow(t *testing.T) {
	const expected = 256.4980392156863
	const eps = 0.0001

	for i := 0; i < 1000; i++ {
		v := expected + float64(i)
		tf := Encode12U(v)

		result := Decode12U(tf)

		if math.Abs(result-expected) > eps {
			t.Logf("Encoded: 0b%b", tf)
			t.Fatalf("%f != %f (i = %d)", result, expected, i)
		}
	}
}

func TestUnsigned12PositiveInfinity(t *testing.T) {
	const expected = 256.4980392156863
	const eps = 0.0001

	v := math.Inf(+1)
	tf := Encode12U(v)

	result := Decode12U(tf)

	if math.Abs(result-expected) > eps {
		t.Logf("Encoded: 0b%b", tf)
		t.Fatalf("%f != %f", result, expected)
	}
}

func TestUnsigned12NegativeInfinity(t *testing.T) {
	v := math.Inf(-1)
	tf := Encode12U(v)

	result := Decode12U(tf)

	if result != 0 {
		t.Fatalf("%f != 0, encoded: 0b%b", result, tf)
	}
}

func TestUnsigned12NaNConvertedToZero(t *testing.T) {
	tf := Encode12U(math.NaN())
	t.Logf("Encoded: 0b%b", tf)

	result := Decode12U(tf)

	if result != 0 {
		t.Fatalf("%f != 0", result)
	}
}

func TestUnsigned12IgnoringMostSignificantByte(t *testing.T) {
	for f := -255.0; f <= 255.0; f += 0.01 {
		toy := Encode12U(f)
		original := Decode12U(toy)

		if 0xF000&toy != 0x0 {
			t.Fatalf("%.4f -> 0b%b (has extra bits)", f, toy)
		}

		for m := 0x1; m < 0xF; m++ {
			modification := uint16(m) << 12
			toyModified := toy | modification
			modified := Decode12U(toyModified)

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

func Test15X3Precision(t *testing.T) {
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
			toy := Encode15X3(tt.number)
			result := Decode15X3(toy)

			diff := math.Abs(result - tt.number)
			if diff > fixPrecision(tt.number, tt.precision) {
				t.Fatalf("%.6f -> 0b%b, diff: %f", tt.number, toy, diff)
			}
		}
	}

	for _, tt := range tests {
		if tt.number <= 4.0 {
			negative := -tt.number
			toy := Encode15X3(negative)
			result := Decode15X3(toy)

			diff := math.Abs(result - negative)
			if diff > fixPrecision(tt.number, tt.precision) {
				t.Fatalf("%.6f -> 0b%b, diff: %f", negative, toy, diff)
			}
		}
	}
}

func Test15X3Zero(t *testing.T) {
	tf := Encode15X3(0)
	t.Logf("Encoded: 0b%b", tf)

	result := Decode15X3(tf)

	if result != 0 {
		t.Fatalf("%f != 0", result)
	}
}

func Test15X3PlusOne(t *testing.T) {
	tf := Encode15X3(1)
	t.Logf("Encoded: 0b%b", tf)

	result := Decode15X3(tf)

	if result != 1 {
		t.Fatalf("%f != 1", result)
	}
}

func Test15X3MinusOne(t *testing.T) {
	tf := Encode15X3(-1)
	t.Logf("Encoded: 0b%b", tf)

	result := Decode15X3(tf)

	if result != -1 {
		t.Fatalf("%f != -1", result)
	}
}

func Test15X3PositiveOverflow(t *testing.T) {
	const expected = 4.046627
	const eps = 0.0001

	for i := 0; i < 1000; i++ {
		v := expected + float64(i)
		tf := Encode15X3(v)

		result := Decode15X3(tf)

		if math.Abs(result-expected) > eps {
			t.Logf("Encoded: 0b%b", tf)
			t.Fatalf("%f != %f (i = %d)", result, expected, i)
		}
	}
}

func Test15X3NegativeOverflow(t *testing.T) {
	const expected = -4.046627
	const eps = 0.0001

	for i := 0; i < 1000; i++ {
		v := expected - float64(i)
		tf := Encode15X3(v)

		result := Decode15X3(tf)

		if math.Abs(result-expected) > eps {
			t.Logf("Encoded: 0b%b", tf)
			t.Fatalf("%f != %f (i = %d)", result, expected, i)
		}
	}
}

func Test15X3PositiveInfinity(t *testing.T) {
	const expected = 4.046627
	const eps = 0.0001

	v := math.Inf(+1)
	tf := Encode15X3(v)

	result := Decode15X3(tf)

	if math.Abs(result-expected) > eps {
		t.Logf("Encoded: 0b%b", tf)
		t.Fatalf("%f != %f", result, expected)
	}
}

func Test15X3NegativeInfinity(t *testing.T) {
	const expected = -4.046627
	const eps = 0.0001

	v := math.Inf(-1)
	tf := Encode15X3(v)

	result := Decode15X3(tf)

	if math.Abs(result-expected) > eps {
		t.Logf("Encoded: 0b%b", tf)
		t.Fatalf("%f != %f", result, expected)
	}
}

func Test15X3NaNConvertedToZero(t *testing.T) {
	tf := Encode15X3(math.NaN())
	t.Logf("Encoded: 0b%b", tf)

	result := Decode15X3(tf)

	if result != 0 {
		t.Fatalf("%f != 0", result)
	}
}

func Test15X3IgnoringMostSignificantBits(t *testing.T) {
	for f := -255.0; f <= 255.0; f += 0.01 {
		toy := Encode15X3(f)
		original := Decode15X3(toy)

		if 0b1000_0000_0000_0000&toy != 0x0 {
			t.Fatalf("%.4f -> 0b%b (has extra bits)", f, toy)
		}

		m := 0b1
		modification := uint16(m) << 15
		toyModified := toy | modification
		modified := Decode15X3(toyModified)

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
	tf := Encode15X3(0.6)
	t.Logf("Encoded: 0b%b", tf)

	input := Decode15X3(tf)
	t.Logf("Decoded: %f", input)

	temp := input

	for i := 0; i < 10; i++ {
		if temp != input {
			t.Fatalf("#%d: %f != %f", i, temp, input)
		}

		tfTemp := Encode15X3(temp)
		t.Logf("Encoded: 0b%b", tfTemp)
		temp = Decode15X3(tfTemp)
		t.Logf("Decoded: %f", temp)
	}
}

// ------------------------

func TestUseDelta15X3(t *testing.T) {
	const eps060 = 0.00024   // x = -1
	const eps030 = 0.00012   // x = -2
	const eps013 = 0.00006   // x = -3
	const epsMin = 0.0000075 // x = -6
	const epsMax = 0.00096   // x = 1

	a := math.Pow(2, -6)
	b := 1 / (1 - a)
	twoPowerM := math.Pow(2, 11)

	last := Decode15X3(Encode15X3(0.6))
	lastTf := Encode15X3(last)
	t.Logf("Last encoded: 0b%b", lastTf)

	resultTf := UseIntegerDelta15X3(lastTf, 0)
	result := Decode15X3(resultTf)
	expected := last

	if math.Abs(result-expected) > eps060 {
		t.Logf("delta=0 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = UseIntegerDelta15X3(lastTf, 1)
	result = Decode15X3(resultTf)
	expected = last + ((1.0/twoPowerM)*0.5)*b

	if math.Abs(result-expected) > eps060 {
		t.Logf("delta=1 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = UseIntegerDelta15X3(lastTf, -1)
	result = Decode15X3(resultTf)
	expected = last - ((1.0/twoPowerM)*0.5)*b

	if math.Abs(result-expected) > eps060 {
		t.Logf("delta=-1 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = UseIntegerDelta15X3(lastTf, 123)
	result = Decode15X3(resultTf)
	expected = last + ((123.0/twoPowerM)*0.5)*b

	if math.Abs(result-expected) > eps060 {
		t.Logf("delta=123 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = UseIntegerDelta15X3(lastTf, -123)
	result = Decode15X3(resultTf)
	expected = last - ((123.0/twoPowerM)*0.5)*b

	if math.Abs(result-expected) > eps060 {
		t.Logf("delta=-123 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	last = Decode15X3(Encode15X3(0.3))
	lastTf = Encode15X3(last)
	t.Logf("Last encoded: 0b%b", lastTf)

	mantissa := lastTf & 0b11111111111
	if mantissa != 499 {
		t.Fatalf("This test is probably broken: mantissa equals %d.", mantissa)
	}

	result = Decode15X3(UseIntegerDelta15X3(lastTf, 0))
	expected = last

	if math.Abs(result-expected) > eps030 {
		t.Fatalf("%f != %f", result, expected)
	}

	result = Decode15X3(UseIntegerDelta15X3(lastTf, 2047-499))
	expected = last + (((2047.0-499.0)/twoPowerM)*0.25)*b

	if math.Abs(result-expected) > eps030 {
		t.Fatalf("%f != %f", result, expected)
	}

	result = Decode15X3(UseIntegerDelta15X3(lastTf, -499))
	expected = last - ((499.0/twoPowerM)*0.25)*b

	if math.Abs(result-expected) > eps030 {
		t.Fatalf("%f != %f", result, expected)
	}

	result = Decode15X3(UseIntegerDelta15X3(lastTf, 2047-499+1))
	expected = last +
		(((2047.0-499.0)/twoPowerM)*0.25)*b +
		((1.0/twoPowerM)*0.5)*b

	if math.Abs(result-expected) > eps060 {
		t.Fatalf("%f != %f", result, expected)
	}

	result = Decode15X3(UseIntegerDelta15X3(lastTf, -500))
	expected = last -
		((499.0/twoPowerM)*0.25)*b -
		((1.0/twoPowerM)*0.125)*b

	if math.Abs(result-expected) > eps013 {
		t.Fatalf("%f != %f", result, expected)
	}

	lastTf = 0b0
	last = Decode15X3(lastTf)
	t.Logf("Last encoded: 0b%b", lastTf)

	resultTf = UseIntegerDelta15X3(lastTf, 0)
	result = Decode15X3(resultTf)
	expected = last

	if math.Abs(result-expected) > epsMin {
		t.Logf("delta=0 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = UseIntegerDelta15X3(lastTf, 1)
	result = Decode15X3(resultTf)
	expected = last + ((1.0/twoPowerM)*0.015625)*b

	if math.Abs(result-expected) > epsMin {
		t.Logf("delta=1 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = UseIntegerDelta15X3(lastTf, -1)
	result = Decode15X3(resultTf)
	expected = -last // minus zero

	if math.Abs(result-expected) > epsMin {
		t.Logf("delta=-1 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	// and back
	minusBit := Encode15X3(1) ^ Encode15X3(-1)
	t.Logf("Minus bit: 0b%b", minusBit)
	lastTf = minusBit | 0b0
	last = Decode15X3(lastTf)
	t.Logf("Last encoded: 0b%b", lastTf)

	resultTf = UseIntegerDelta15X3(lastTf, 1)
	result = Decode15X3(resultTf)
	expected = -last // plus zero

	if math.Abs(result-expected) > epsMin {
		t.Logf("delta=1 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = UseIntegerDelta15X3(lastTf, -10)
	result = Decode15X3(resultTf)
	expected = last - ((10.0/twoPowerM)*0.015625)*b

	if math.Abs(result-expected) > epsMin {
		t.Logf("delta=-10 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	// big deltas

	resultTf = UseIntegerDelta15X3(lastTf, 234234)
	result = Decode15X3(resultTf)
	expected = ((1.0+(twoPowerM-1)/twoPowerM)*2 - a) * b

	if math.Abs(result-expected) > epsMax {
		t.Logf("delta=234234 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = UseIntegerDelta15X3(lastTf, -16382)
	result = Decode15X3(resultTf)
	expected = -((1.0+(twoPowerM-2)/twoPowerM)*2 - a) * b

	if math.Abs(result-expected) > epsMax {
		t.Logf("delta=-16382 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = UseIntegerDelta15X3(lastTf, -16383)
	result = Decode15X3(resultTf)
	expected = -((1.0+(twoPowerM-1)/twoPowerM)*2 - a) * b

	if math.Abs(result-expected) > epsMax {
		t.Logf("delta=-16383 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = UseIntegerDelta15X3(lastTf, -16384)
	result = Decode15X3(resultTf)
	expected = -((1.0+(twoPowerM-1)/twoPowerM)*2 - a) * b

	if math.Abs(result-expected) > epsMax {
		t.Logf("delta=-16384 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = UseIntegerDelta15X3(lastTf, -26384)
	result = Decode15X3(resultTf)
	expected = -((1.0+(twoPowerM-1)/twoPowerM)*2 - a) * b

	if math.Abs(result-expected) > epsMax {
		t.Logf("delta=-26384 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}
}

func TestGetDelta15X3(t *testing.T) {
	const start = -4.0
	const stop = 4.0
	const step = 0.01

	const eps = 1e-9

	last := Encode15X3(start)

	for x := start + step; x <= stop; x += step {
		xtf := Encode15X3(x)
		expected := Decode15X3(xtf)

		delta := GetIntegerDelta15X3(last, xtf)
		resultTf := UseIntegerDelta15X3(last, delta)
		result := Decode15X3(resultTf)

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

func TestGetDelta12(t *testing.T) {
	const start = -256.0
	const stop = 256.0
	const step = 0.01

	const eps = 1e-9

	last := Encode12(start)

	for x := start + step; x <= stop; x += step {
		xtf := Encode12(x)
		expected := Decode12(xtf)

		delta := GetIntegerDelta12(last, xtf)
		resultTf := UseIntegerDelta12(last, delta)
		result := Decode12(resultTf)

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

func TestUseDeltaUnsigned(t *testing.T) {
	const eps060 = 0.00195  // x = -1
	const eps030 = 0.00097  // x = -2
	const eps013 = 0.00048  // x = -3
	const epsMin = 0.000015 // x = -8
	const epsMax = 0.49     // x = 7

	a := math.Pow(2, -8)
	b := 1 / (1 - a)
	twoPowerM := math.Pow(2, 8)

	last := Decode12U(Encode12U(0.6))
	lastTf := Encode12U(last)
	t.Logf("Last encoded: 0b%b", lastTf)

	resultTf := UseIntegerDelta12U(lastTf, 0)
	result := Decode12U(resultTf)
	expected := last

	if math.Abs(result-expected) > eps060 {
		t.Logf("delta=0 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = UseIntegerDelta12U(lastTf, 1)
	result = Decode12U(resultTf)
	expected = last + ((1.0/twoPowerM)*0.5)*b

	if math.Abs(result-expected) > eps060 {
		t.Logf("delta=1 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = UseIntegerDelta12U(lastTf, -1)
	result = Decode12U(resultTf)
	expected = last - ((1.0/twoPowerM)*0.5)*b

	if math.Abs(result-expected) > eps060 {
		t.Logf("delta=-1 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = UseIntegerDelta12U(lastTf, 45)
	result = Decode12U(resultTf)
	expected = last + ((45.0/twoPowerM)*0.5)*b

	if math.Abs(result-expected) > eps060 {
		t.Logf("delta=123 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = UseIntegerDelta12U(lastTf, -45)
	result = Decode12U(resultTf)
	expected = last - ((45.0/twoPowerM)*0.5)*b

	if math.Abs(result-expected) > eps060 {
		t.Logf("delta=-123 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	last = Decode12U(Encode12U(0.3))
	lastTf = Encode12U(last)
	t.Logf("Last encoded: 0b%b", lastTf)

	mantissa := lastTf & 0b11111111
	if mantissa != 54 {
		t.Fatalf("This test is probably broken: mantissa equals %d.", mantissa)
	}

	result = Decode12U(UseIntegerDelta12U(lastTf, 0))
	expected = last

	if math.Abs(result-expected) > eps030 {
		t.Fatalf("%f != %f", result, expected)
	}

	result = Decode12U(UseIntegerDelta12U(lastTf, 255-54))
	expected = last + (((255.0-54.0)/twoPowerM)*0.25)*b

	if math.Abs(result-expected) > eps030 {
		t.Fatalf("%f != %f", result, expected)
	}

	result = Decode12U(UseIntegerDelta12U(lastTf, -54))
	expected = last - ((54.0/twoPowerM)*0.25)*b

	if math.Abs(result-expected) > eps030 {
		t.Fatalf("%f != %f", result, expected)
	}

	result = Decode12U(UseIntegerDelta12U(lastTf, 255-54+1))
	expected = last +
		(((255.0-54.0)/twoPowerM)*0.25)*b +
		((1.0/twoPowerM)*0.5)*b

	if math.Abs(result-expected) > eps060 {
		t.Fatalf("%f != %f", result, expected)
	}

	result = Decode12U(UseIntegerDelta12U(lastTf, -55))
	expected = last -
		((54.0/twoPowerM)*0.25)*b -
		((1.0/twoPowerM)*0.125)*b

	if math.Abs(result-expected) > eps013 {
		t.Fatalf("%f != %f", result, expected)
	}

	lastTf = 0b0
	last = Decode12U(lastTf)
	t.Logf("Last encoded: 0b%b", lastTf)

	resultTf = UseIntegerDelta12U(lastTf, 0)
	result = Decode12U(resultTf)
	expected = last

	if math.Abs(result-expected) > epsMin {
		t.Logf("delta=0 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = UseIntegerDelta12U(lastTf, 1)
	result = Decode12U(resultTf)
	expected = last + ((1.0/twoPowerM)*0.00390625)*b

	if math.Abs(result-expected) > epsMin {
		t.Logf("delta=1 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = UseIntegerDelta12U(lastTf, -1)
	result = Decode12U(resultTf)
	expected = last // zero is the minimum value

	if result != expected {
		t.Logf("delta=-1 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = UseIntegerDelta12U(lastTf, -10)
	result = Decode12U(resultTf)
	expected = last // zero is the minimum value

	if result != expected {
		t.Logf("delta=-10 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	// big deltas

	resultTf = UseIntegerDelta12U(lastTf, 234234)
	result = Decode12U(resultTf)
	expected = ((1.0+(twoPowerM-1)/twoPowerM)*128 - a) * b

	if math.Abs(result-expected) > epsMax {
		t.Logf("delta=234234 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = UseIntegerDelta12U(lastTf, -16382)
	result = Decode12U(resultTf)
	expected = 0.0

	if math.Abs(result-expected) > epsMin {
		t.Logf("delta=-16382 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	last = Decode12U(Encode12U(0.234))
	lastTf = Encode12U(last)

	resultTf = UseIntegerDelta12U(lastTf, -26384)
	result = Decode12U(resultTf)
	expected = 0.0

	if math.Abs(result-expected) > epsMin {
		t.Logf("delta=-26384 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}
}

func TestGetDeltaUnsigned(t *testing.T) {
	const start = 0.0
	const stop = 256.0
	const step = 0.01

	const eps = 1e-9

	last := Encode12U(start)

	for x := start + step; x <= stop; x += step {
		xtf := Encode12U(x)
		expected := Decode12U(xtf)

		delta := GetIntegerDelta12U(last, xtf)
		resultTf := UseIntegerDelta12U(last, delta)
		result := Decode12U(resultTf)

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

func TestGetDelta14(t *testing.T) {
	const start = -256.0
	const stop = 256.0
	const step = 0.01

	const eps = 1e-9

	last := Encode14(start)

	for x := start + step; x <= stop; x += step {
		xtf := Encode14(x)
		expected := Decode14(xtf)

		delta := GetIntegerDelta14(last, xtf)
		resultTf := UseIntegerDelta14(last, delta)
		result := Decode14(resultTf)

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

func Test12MinusBitPosition(t *testing.T) {
	tf := Encode12(42)
	t.Logf("Encoded: 0b%b", tf)

	a := Decode12(tf)
	b := -Decode12(tf | 0b1000_0000_0000)

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
	b := -Decode14(tf | 0b10_0000_0000_0000)

	if a != b {
		t.Fatalf("%f != %f", a, b)
	}
}

func Test15X3MinusBitPosition(t *testing.T) {
	tf := Encode15X3(42)
	t.Logf("Encoded: 0b%b", tf)

	a := Decode15X3(tf)
	b := -Decode15X3(tf | 0b0100_0000_0000_0000)

	if a != b {
		t.Fatalf("%f != %f", a, b)
	}
}

func TestReadme(t *testing.T) {
	const input = 0.345
	const eps = 1e-6

	{
		tf := Encode12(input)
		if tf != 0x332 {
			t.Fatalf("Incorrect encoded: 0x%X (12-bit)\n", tf)
		}

		result := Decode12(tf)
		if math.Abs(result-0.345098) > eps {
			t.Fatalf("Incorrect decoded: %f (12-bit)\n", result)
		}
	}

	{
		tf := Encode12U(input)
		if tf != 0x664 {
			t.Fatalf("Incorrect encoded: 0x%X (12-bit unsigned)\n", tf)
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
		tf := Encode14(input)
		if tf != 0xCC8 {
			t.Fatalf("Incorrect encoded: 0x%X (14-bit)\n", tf)
		}

		result := Decode14(tf)
		if math.Abs(result-0.345098) > eps {
			t.Fatalf("Incorrect decoded: %f (14-bit)\n", result)
		}
	}

	{
		tf := Encode15X3(input)
		if tf != 0x235E {
			t.Fatalf("Incorrect encoded: 0x%X (15x3)\n", tf)
		}

		result := Decode15X3(tf)
		if math.Abs(result-0.344990) > eps {
			t.Fatalf("Incorrect decoded: %f (15x3)\n", result)
		}
	}

	{
		tf := Encode15X3(input)
		result := Decode15X3(tf)
		if math.Abs(result-0.344990) > eps {
			t.Fatalf("Incorrect decoded: %f (15x3)\n", result)
		}
	}

	{
		series := []float64{-0.0058, 0.01, 0.123, 0.134, 0.132, 0.144, 0.145, 0.140}
		expected := []int{387, 414, 12, -2, 12, 1, -5}

		previous := Encode12(series[0])
		for i := 1; i < len(series); i++ {
			this := Encode12(series[i])

			delta := GetIntegerDelta12(previous, this)
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
