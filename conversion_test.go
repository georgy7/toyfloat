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

func checkBadArguments(t *testing.T, i, min, max int, e error, remark string) {
	inLimits := (i >= min) && (i <= max)

	comment := remark
	if len(comment) > 0 {
		comment = " (" + remark + ")"
	}

	if (e != nil) && inLimits {
		t.Fatalf("length=%d must work%s", i, comment)
	}

	if (e == nil) && !inLimits {
		t.Fatalf("length=%d must result in an error%s", i, comment)
	}
}

func TestX2BadArguments(t *testing.T) {
	for i := -10; i <= 20; i++ {
		{
			_, e := NewTypeX2(i, true)
			checkBadArguments(t, i, 4, 16, e, "")
		}
		{
			_, e := NewTypeX2(i, false)
			checkBadArguments(t, i, 3, 16, e, "unsigned")
		}
	}
}

func TestX3BadArguments(t *testing.T) {
	for i := -10; i <= 20; i++ {
		{
			_, e := NewTypeX3(i, true)
			checkBadArguments(t, i, 5, 16, e, "")
		}
		{
			_, e := NewTypeX3(i, false)
			checkBadArguments(t, i, 4, 16, e, "unsigned")
		}
	}
}

func TestX4BadArguments(t *testing.T) {
	for i := -10; i <= 20; i++ {
		{
			_, e := NewTypeX4(i, true)
			checkBadArguments(t, i, 6, 16, e, "")
		}
		{
			_, e := NewTypeX4(i, false)
			checkBadArguments(t, i, 5, 16, e, "unsigned")
		}
	}
}

func makeTypeX2(length int, signed bool, t *testing.T) Type {
	tfType, err := NewTypeX2(length, signed)
	if err != nil {
		t.Fatal(err)
	}

	return tfType
}

func makeTypeX3(length int, signed bool, t *testing.T) Type {
	tfType, err := NewTypeX3(length, signed)
	if err != nil {
		t.Fatal(err)
	}

	return tfType
}

func makeTypeX4(length int, signed bool, t *testing.T) Type {
	tfType, err := NewTypeX4(length, signed)
	if err != nil {
		t.Fatal(err)
	}

	return tfType
}

// ------------------------

func Test12Zero(t *testing.T) {
	toyfloat12 := makeTypeX4(12, true, t)

	tf := toyfloat12.Encode(0)
	t.Logf("Encoded: 0b%b", tf)

	result := toyfloat12.Decode(tf)

	if result != 0 {
		t.Fatalf("%f != 0", result)
	}
}

func Test12PlusOne(t *testing.T) {
	toyfloat12 := makeTypeX4(12, true, t)

	tf := toyfloat12.Encode(1)
	t.Logf("Encoded: 0b%b", tf)

	result := toyfloat12.Decode(tf)

	if result != 1 {
		t.Fatalf("%f != 1", result)
	}
}

func Test12MinusOne(t *testing.T) {
	toyfloat12 := makeTypeX4(12, true, t)

	tf := toyfloat12.Encode(-1)
	t.Logf("Encoded: 0b%b", tf)

	result := toyfloat12.Decode(tf)

	if result != -1 {
		t.Fatalf("%f != -1", result)
	}
}

func Test12Abs(t *testing.T) {
	toyfloat12 := makeTypeX4(12, true, t)
	input := toyfloat12.Encode(12.344)
	result := toyfloat12.Abs(input)
	if result != input {
		t.Fatalf("Abs(positive): %d != %d", result, input)
	}

	input = toyfloat12.Encode(-15.2)
	result = toyfloat12.Abs(input)

	const eps = 0.0001
	want := math.Abs(toyfloat12.Decode(input))
	got := toyfloat12.Decode(result)

	if math.Abs(got-want) > eps {
		t.Fatalf("Abs(negative): %f != %f", got, want)
	}
}

func Test12MaxValue(t *testing.T) {
	toyfloat12 := makeTypeX4(12, true, t)
	const expected = 255.99607843137255
	const eps = 0.0001
	got := toyfloat12.MaxValue()
	if math.Abs(got-expected) > eps {
		t.Fatalf("got %f, expected %f", got, expected)
	}
}

func Test12MinValue(t *testing.T) {
	toyfloat12 := makeTypeX4(12, true, t)
	const expected = -255.99607843137255
	const eps = 0.0001
	got := toyfloat12.MinValue()
	if math.Abs(got-expected) > eps {
		t.Fatalf("got %f, expected %f", got, expected)
	}
}

func Test12PositiveOverflow(t *testing.T) {
	toyfloat12 := makeTypeX4(12, true, t)
	const expected = 255.99607843137255
	const eps = 0.0001

	for i := 0; i < 1000; i++ {
		v := expected + float64(i)
		tf := toyfloat12.Encode(v)

		result := toyfloat12.Decode(tf)

		if math.Abs(result-expected) > eps {
			t.Logf("Encoded: 0b%b", tf)
			t.Fatalf("%f != %f (i = %d)", result, expected, i)
		}
	}
}

func Test12NegativeOverflow(t *testing.T) {
	toyfloat12 := makeTypeX4(12, true, t)
	const expected = -255.99607843137255
	const eps = 0.0001

	for i := 0; i < 1000; i++ {
		v := expected - float64(i)
		tf := toyfloat12.Encode(v)

		result := toyfloat12.Decode(tf)

		if math.Abs(result-expected) > eps {
			t.Logf("Encoded: 0b%b", tf)
			t.Fatalf("%f != %f (i = %d)", result, expected, i)
		}
	}
}

func Test12PositiveInfinity(t *testing.T) {
	toyfloat12 := makeTypeX4(12, true, t)
	const expected = 255.99607843137255
	const eps = 0.0001

	v := math.Inf(+1)
	tf := toyfloat12.Encode(v)

	result := toyfloat12.Decode(tf)

	if math.Abs(result-expected) > eps {
		t.Logf("Encoded: 0b%b", tf)
		t.Fatalf("%f != %f", result, expected)
	}
}

func Test12NegativeInfinity(t *testing.T) {
	toyfloat12 := makeTypeX4(12, true, t)
	const expected = -255.99607843137255
	const eps = 0.0001

	v := math.Inf(-1)
	tf := toyfloat12.Encode(v)

	result := toyfloat12.Decode(tf)

	if math.Abs(result-expected) > eps {
		t.Logf("Encoded: 0b%b", tf)
		t.Fatalf("%f != %f", result, expected)
	}
}

func Test12NaNConvertedToZero(t *testing.T) {
	toyfloat12 := makeTypeX4(12, true, t)

	tf := toyfloat12.Encode(math.NaN())
	t.Logf("Encoded: 0b%b", tf)

	result := toyfloat12.Decode(tf)

	if result != 0 {
		t.Fatalf("%f != 0", result)
	}
}

func Test12Precision(t *testing.T) {
	toyfloat12 := makeTypeX4(12, true, t)

	tests := getToyfloatPositiveSample()

	for _, tt := range tests {
		toy := toyfloat12.Encode(tt.number)
		result := toyfloat12.Decode(toy)

		diff := math.Abs(result - tt.number)
		if diff > tt.precision {
			t.Fatalf("%.4f -> 0b%b, diff: %f", tt.number, toy, diff)
		}
	}

	for _, tt := range tests {
		negative := -tt.number
		toy := toyfloat12.Encode(negative)
		result := toyfloat12.Decode(toy)

		diff := math.Abs(result - negative)
		if diff > tt.precision {
			t.Fatalf("%.4f -> 0b%b, diff: %f", negative, toy, diff)
		}
	}
}

func Test12IgnoringMostSignificantBits(t *testing.T) {
	toyfloat12 := makeTypeX4(12, true, t)

	for f := -255.0; f <= 255.0; f += 0.01 {
		toy := toyfloat12.Encode(f)
		original := toyfloat12.Decode(toy)

		if 0xF000&toy != 0x0 {
			t.Fatalf("%.4f -> 0b%b (has extra bits)", f, toy)
		}

		for m := 0x1; m <= 0xF; m++ {
			modification := uint16(m) << 12
			toyModified := toy | modification
			modified := toyfloat12.Decode(toyModified)

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

func BenchmarkFloat64IncrementAsAReference(b *testing.B) {
	counter := 0.0
	for i := 0; i < b.N; i++ {
		counter++
	}
}

func BenchmarkCreateTypeX4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, e := NewTypeX4(12, true)
		if e != nil {
			b.Fatal(e)
		}
	}
}

func BenchmarkCreateTypeX3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, e := NewTypeX3(12, true)
		if e != nil {
			b.Fatal(e)
		}
	}
}

func BenchmarkCreateTypeX2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, e := NewTypeX2(12, true)
		if e != nil {
			b.Fatal(e)
		}
	}
}

func BenchmarkEncode(b *testing.B) {
	toyfloat12, e := NewTypeX4(12, true)
	if e != nil {
		b.Fatal(e)
	}

	const scale = 256.0 / 10000
	for i := 0; i < b.N; i++ {
		_ = toyfloat12.Encode(scale * float64(i%10000))
	}
}

func BenchmarkDecode(b *testing.B) {
	toyfloat12, e := NewTypeX4(12, true)
	if e != nil {
		b.Fatal(e)
	}

	for i := 0; i < b.N; i++ {
		_ = toyfloat12.Decode(uint16(i))
	}
}

func BenchmarkEncode12X2(b *testing.B) {
	toyfloat12x2, e := NewTypeX2(12, true)
	if e != nil {
		b.Fatal(e)
	}

	const scale = 3.0 / 10000
	for i := 0; i < b.N; i++ {
		_ = toyfloat12x2.Encode(scale * float64(i%10000))
	}
}

func BenchmarkDecode12X2(b *testing.B) {
	toyfloat12x2, e := NewTypeX2(12, true)
	if e != nil {
		b.Fatal(e)
	}

	for i := 0; i < b.N; i++ {
		_ = toyfloat12x2.Decode(uint16(i))
	}
}

func BenchmarkGetDelta(b *testing.B) {
	toyfloat12, e := NewTypeX4(12, true)
	if e != nil {
		b.Fatal(e)
	}

	last := uint16(0)
	for i := 0; i < b.N; i++ {
		kindOfRandom := last + 7*uint16(i)
		_ = toyfloat12.GetIntegerDelta(last, kindOfRandom)
		last = kindOfRandom
	}
}

func BenchmarkGetDeltaX2(b *testing.B) {
	toyfloat12x2, e := NewTypeX2(12, true)
	if e != nil {
		b.Fatal(e)
	}

	last := uint16(0)
	for i := 0; i < b.N; i++ {
		kindOfRandom := last + 7*uint16(i)
		_ = toyfloat12x2.GetIntegerDelta(last, kindOfRandom)
		last = kindOfRandom
	}
}

func BenchmarkUseDelta(b *testing.B) {
	toyfloat12, e := NewTypeX4(12, true)
	if e != nil {
		b.Fatal(e)
	}

	last := uint16(0)
	for i := 0; i < b.N; i++ {
		last = toyfloat12.UseIntegerDelta(last, i%255-128)
	}
}

func BenchmarkUseDeltaX2(b *testing.B) {
	toyfloat12x2, e := NewTypeX2(12, true)
	if e != nil {
		b.Fatal(e)
	}

	last := uint16(0)
	for i := 0; i < b.N; i++ {
		last = toyfloat12x2.UseIntegerDelta(last, i%255-128)
	}
}

// ------------------------

func TestUnsigned12Precision(t *testing.T) {
	toyfloat12u := makeTypeX4(12, false, t)

	tests := getToyfloatPositiveSample()

	for _, tt := range tests {
		toy := toyfloat12u.Encode(tt.number)
		result := toyfloat12u.Decode(toy)

		diff := math.Abs(result - tt.number)
		if diff > tt.precision*0.5 {
			t.Fatalf("%.4f -> 0b%b, diff: %f", tt.number, toy, diff)
		}
	}
}

func TestUnsigned12NegativeInput(t *testing.T) {
	toyfloat12u := makeTypeX4(12, false, t)

	tests := getToyfloatPositiveSample()

	for _, tt := range tests {
		negative := -tt.number
		toy := toyfloat12u.Encode(negative)
		result := toyfloat12u.Decode(toy)

		if result != 0 {
			t.Fatalf("%f != 0", result)
		}
	}
}

func TestUnsigned12Zero(t *testing.T) {
	toyfloat12u := makeTypeX4(12, false, t)

	tf := toyfloat12u.Encode(0)
	t.Logf("Encoded: 0b%b", tf)

	result := toyfloat12u.Decode(tf)

	if result != 0 {
		t.Fatalf("%f != 0", result)
	}
}

func TestUnsigned12PlusOne(t *testing.T) {
	toyfloat12u := makeTypeX4(12, false, t)

	tf := toyfloat12u.Encode(1)
	t.Logf("Encoded: 0b%b", tf)

	result := toyfloat12u.Decode(tf)

	if result != 1 {
		t.Fatalf("%f != 1", result)
	}
}

func TestUnsigned12MaxValue(t *testing.T) {
	toyfloat12u := makeTypeX4(12, false, t)
	const expected = 256.4980392156863
	const eps = 0.0001
	got := toyfloat12u.MaxValue()
	if math.Abs(got-expected) > eps {
		t.Fatalf("got %f, expected %f", got, expected)
	}
}

func TestUnsigned12MinValue(t *testing.T) {
	toyfloat12u := makeTypeX4(12, false, t)
	got := toyfloat12u.MinValue()
	if got != 0 {
		t.Fatalf("got %f, expected %f", got, 0.0)
	}
}

func TestUnsigned12PositiveOverflow(t *testing.T) {
	toyfloat12u := makeTypeX4(12, false, t)
	const expected = 256.4980392156863
	const eps = 0.0001

	for i := 0; i < 1000; i++ {
		v := expected + float64(i)
		tf := toyfloat12u.Encode(v)

		result := toyfloat12u.Decode(tf)

		if math.Abs(result-expected) > eps {
			t.Logf("Encoded: 0b%b", tf)
			t.Fatalf("%f != %f (i = %d)", result, expected, i)
		}
	}
}

func TestUnsigned12PositiveInfinity(t *testing.T) {
	toyfloat12u := makeTypeX4(12, false, t)
	const expected = 256.4980392156863
	const eps = 0.0001

	v := math.Inf(+1)
	tf := toyfloat12u.Encode(v)

	result := toyfloat12u.Decode(tf)

	if math.Abs(result-expected) > eps {
		t.Logf("Encoded: 0b%b", tf)
		t.Fatalf("%f != %f", result, expected)
	}
}

func TestUnsigned12NegativeInfinity(t *testing.T) {
	toyfloat12u := makeTypeX4(12, false, t)

	v := math.Inf(-1)
	tf := toyfloat12u.Encode(v)

	result := toyfloat12u.Decode(tf)

	if result != 0 {
		t.Fatalf("%f != 0, encoded: 0b%b", result, tf)
	}
}

func TestUnsigned12NaNConvertedToZero(t *testing.T) {
	toyfloat12u := makeTypeX4(12, false, t)

	tf := toyfloat12u.Encode(math.NaN())
	t.Logf("Encoded: 0b%b", tf)

	result := toyfloat12u.Decode(tf)

	if result != 0 {
		t.Fatalf("%f != 0", result)
	}
}

func TestUnsigned12IgnoringMostSignificantByte(t *testing.T) {
	toyfloat12u := makeTypeX4(12, false, t)

	for f := -255.0; f <= 255.0; f += 0.01 {
		toy := toyfloat12u.Encode(f)
		original := toyfloat12u.Decode(toy)

		if 0xF000&toy != 0x0 {
			t.Fatalf("%.4f -> 0b%b (has extra bits)", f, toy)
		}

		for m := 0x1; m < 0xF; m++ {
			modification := uint16(m) << 12
			toyModified := toy | modification
			modified := toyfloat12u.Decode(toyModified)

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
	toyfloat13 := makeTypeX4(13, true, t)

	tests := getToyfloatPositiveSample()

	for _, tt := range tests {
		toy := toyfloat13.Encode(tt.number)
		result := toyfloat13.Decode(toy)

		diff := math.Abs(result - tt.number)
		if diff > tt.precision*0.5 {
			t.Fatalf("%.4f -> 0b%b, diff: %f", tt.number, toy, diff)
		}
	}

	for _, tt := range tests {
		negative := -tt.number
		toy := toyfloat13.Encode(negative)
		result := toyfloat13.Decode(toy)

		diff := math.Abs(result - negative)
		if diff > tt.precision*0.5 {
			t.Fatalf("%.4f -> 0b%b, diff: %f", negative, toy, diff)
		}
	}
}

func Test13Zero(t *testing.T) {
	toyfloat13 := makeTypeX4(13, true, t)

	tf := toyfloat13.Encode(0)
	t.Logf("Encoded: 0b%b", tf)

	result := toyfloat13.Decode(tf)

	if result != 0 {
		t.Fatalf("%f != 0", result)
	}
}

func Test13PlusOne(t *testing.T) {
	toyfloat13 := makeTypeX4(13, true, t)

	tf := toyfloat13.Encode(1)
	t.Logf("Encoded: 0b%b", tf)

	result := toyfloat13.Decode(tf)

	if result != 1 {
		t.Fatalf("%f != 1", result)
	}
}

func Test13MinusOne(t *testing.T) {
	toyfloat13 := makeTypeX4(13, true, t)

	tf := toyfloat13.Encode(-1)
	t.Logf("Encoded: 0b%b", tf)

	result := toyfloat13.Decode(tf)

	if result != -1 {
		t.Fatalf("%f != -1", result)
	}
}

func Test13PositiveOverflow(t *testing.T) {
	toyfloat13 := makeTypeX4(13, true, t)

	const expected = 256.4980392156863
	const eps = 0.0001

	for i := 0; i < 1000; i++ {
		v := expected + float64(i)
		tf := toyfloat13.Encode(v)

		result := toyfloat13.Decode(tf)

		if math.Abs(result-expected) > eps {
			t.Logf("Encoded: 0b%b", tf)
			t.Fatalf("%f != %f (i = %d)", result, expected, i)
		}
	}
}

func Test13NegativeOverflow(t *testing.T) {
	toyfloat13 := makeTypeX4(13, true, t)

	const expected = -256.4980392156863
	const eps = 0.0001

	for i := 0; i < 1000; i++ {
		v := expected - float64(i)
		tf := toyfloat13.Encode(v)

		result := toyfloat13.Decode(tf)

		if math.Abs(result-expected) > eps {
			t.Logf("Encoded: 0b%b", tf)
			t.Fatalf("%f != %f (i = %d)", result, expected, i)
		}
	}
}

func Test13PositiveInfinity(t *testing.T) {
	toyfloat13 := makeTypeX4(13, true, t)

	const expected = 256.4980392156863
	const eps = 0.0001

	v := math.Inf(+1)
	tf := toyfloat13.Encode(v)

	result := toyfloat13.Decode(tf)

	if math.Abs(result-expected) > eps {
		t.Logf("Encoded: 0b%b", tf)
		t.Fatalf("%f != %f", result, expected)
	}
}

func Test13NegativeInfinity(t *testing.T) {
	toyfloat13 := makeTypeX4(13, true, t)

	const expected = -256.4980392156863
	const eps = 0.0001

	v := math.Inf(-1)
	tf := toyfloat13.Encode(v)

	result := toyfloat13.Decode(tf)

	if math.Abs(result-expected) > eps {
		t.Logf("Encoded: 0b%b", tf)
		t.Fatalf("%f != %f", result, expected)
	}
}

func Test13NaNConvertedToZero(t *testing.T) {
	toyfloat13 := makeTypeX4(13, true, t)

	tf := toyfloat13.Encode(math.NaN())
	t.Logf("Encoded: 0b%b", tf)

	result := toyfloat13.Decode(tf)

	if result != 0 {
		t.Fatalf("%f != 0", result)
	}
}

func Test13IgnoringMostSignificantBits(t *testing.T) {
	toyfloat13 := makeTypeX4(13, true, t)

	for f := -255.0; f <= 255.0; f += 0.01 {
		toy := toyfloat13.Encode(f)
		original := toyfloat13.Decode(toy)

		if 0b1110_0000_0000_0000&toy != 0x0 {
			t.Fatalf("%.4f -> 0b%b (has extra bits)", f, toy)
		}

		for m := 0b1; m <= 0b111; m++ {
			modification := uint16(m) << 13
			toyModified := toy | modification
			modified := toyfloat13.Decode(toyModified)

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
	toyfloat14 := makeTypeX4(14, true, t)

	tests := getToyfloatPositiveSample()

	for _, tt := range tests {
		toy := toyfloat14.Encode(tt.number)
		result := toyfloat14.Decode(toy)

		diff := math.Abs(result - tt.number)
		if diff > tt.precision*0.25 {
			t.Fatalf("%.4f -> 0b%b, diff: %f", tt.number, toy, diff)
		}
	}

	for _, tt := range tests {
		negative := -tt.number
		toy := toyfloat14.Encode(negative)
		result := toyfloat14.Decode(toy)

		diff := math.Abs(result - negative)
		if diff > tt.precision*0.25 {
			t.Fatalf("%.4f -> 0b%b, diff: %f", negative, toy, diff)
		}
	}
}

func Test14IgnoringMostSignificantBits(t *testing.T) {
	toyfloat14 := makeTypeX4(14, true, t)

	for f := -255.0; f <= 255.0; f += 0.01 {
		toy := toyfloat14.Encode(f)
		original := toyfloat14.Decode(toy)

		if 0b1100_0000_0000_0000&toy != 0x0 {
			t.Fatalf("%.4f -> 0b%b (has extra bits)", f, toy)
		}

		for m := 0b01; m <= 0b11; m++ {
			modification := uint16(m) << 14
			toyModified := toy | modification
			modified := toyfloat14.Decode(toyModified)

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
	toyfloat15x3 := makeTypeX3(15, true, t)

	tests := getToyfloatPositiveSample()

	fixPrecision := func(number, precision float64) float64 {
		if number < 0.0158 {
			return 0.000008
		}
		return precision * 0.0625
	}

	for _, tt := range tests {
		if tt.number <= 4.0 {
			toy := toyfloat15x3.Encode(tt.number)
			result := toyfloat15x3.Decode(toy)

			diff := math.Abs(result - tt.number)
			if diff > fixPrecision(tt.number, tt.precision) {
				t.Fatalf("%.6f -> 0b%b, diff: %f", tt.number, toy, diff)
			}
		}
	}

	for _, tt := range tests {
		if tt.number <= 4.0 {
			negative := -tt.number
			toy := toyfloat15x3.Encode(negative)
			result := toyfloat15x3.Decode(toy)

			diff := math.Abs(result - negative)
			if diff > fixPrecision(tt.number, tt.precision) {
				t.Fatalf("%.6f -> 0b%b, diff: %f", negative, toy, diff)
			}
		}
	}
}

func Test15X3Zero(t *testing.T) {
	toyfloat15x3 := makeTypeX3(15, true, t)

	tf := toyfloat15x3.Encode(0)
	t.Logf("Encoded: 0b%b", tf)

	result := toyfloat15x3.Decode(tf)

	if result != 0 {
		t.Fatalf("%f != 0", result)
	}
}

func Test15X3PlusOne(t *testing.T) {
	toyfloat15x3 := makeTypeX3(15, true, t)

	tf := toyfloat15x3.Encode(1)
	t.Logf("Encoded: 0b%b", tf)

	result := toyfloat15x3.Decode(tf)

	if result != 1 {
		t.Fatalf("%f != 1", result)
	}
}

func Test15X3MinusOne(t *testing.T) {
	toyfloat15x3 := makeTypeX3(15, true, t)

	tf := toyfloat15x3.Encode(-1)
	t.Logf("Encoded: 0b%b", tf)

	result := toyfloat15x3.Decode(tf)

	if result != -1 {
		t.Fatalf("%f != -1", result)
	}
}

func Test15X3MaxValue(t *testing.T) {
	toyfloat15x3 := makeTypeX3(15, true, t)
	const expected = 4.046627
	const eps = 0.0001
	got := toyfloat15x3.MaxValue()
	if math.Abs(got-expected) > eps {
		t.Fatalf("got %f, expected %f", got, expected)
	}
}

func Test15X3MinValue(t *testing.T) {
	toyfloat15x3 := makeTypeX3(15, true, t)
	const expected = -4.046627
	const eps = 0.0001
	got := toyfloat15x3.MinValue()
	if math.Abs(got-expected) > eps {
		t.Fatalf("got %f, expected %f", got, expected)
	}
}

func Test15X3PositiveOverflow(t *testing.T) {
	toyfloat15x3 := makeTypeX3(15, true, t)
	const expected = 4.046627
	const eps = 0.0001

	for i := 0; i < 1000; i++ {
		v := expected + float64(i)
		tf := toyfloat15x3.Encode(v)

		result := toyfloat15x3.Decode(tf)

		if math.Abs(result-expected) > eps {
			t.Logf("Encoded: 0b%b", tf)
			t.Fatalf("%f != %f (i = %d)", result, expected, i)
		}
	}
}

func Test15X3NegativeOverflow(t *testing.T) {
	toyfloat15x3 := makeTypeX3(15, true, t)
	const expected = -4.046627
	const eps = 0.0001

	for i := 0; i < 1000; i++ {
		v := expected - float64(i)
		tf := toyfloat15x3.Encode(v)

		result := toyfloat15x3.Decode(tf)

		if math.Abs(result-expected) > eps {
			t.Logf("Encoded: 0b%b", tf)
			t.Fatalf("%f != %f (i = %d)", result, expected, i)
		}
	}
}

func Test15X3PositiveInfinity(t *testing.T) {
	toyfloat15x3 := makeTypeX3(15, true, t)
	const expected = 4.046627
	const eps = 0.0001

	v := math.Inf(+1)
	tf := toyfloat15x3.Encode(v)

	result := toyfloat15x3.Decode(tf)

	if math.Abs(result-expected) > eps {
		t.Logf("Encoded: 0b%b", tf)
		t.Fatalf("%f != %f", result, expected)
	}
}

func Test15X3NegativeInfinity(t *testing.T) {
	toyfloat15x3 := makeTypeX3(15, true, t)
	const expected = -4.046627
	const eps = 0.0001

	v := math.Inf(-1)
	tf := toyfloat15x3.Encode(v)

	result := toyfloat15x3.Decode(tf)

	if math.Abs(result-expected) > eps {
		t.Logf("Encoded: 0b%b", tf)
		t.Fatalf("%f != %f", result, expected)
	}
}

func Test15X3NaNConvertedToZero(t *testing.T) {
	toyfloat15x3 := makeTypeX3(15, true, t)

	tf := toyfloat15x3.Encode(math.NaN())
	t.Logf("Encoded: 0b%b", tf)

	result := toyfloat15x3.Decode(tf)

	if result != 0 {
		t.Fatalf("%f != 0", result)
	}
}

func Test15X3IgnoringMostSignificantBits(t *testing.T) {
	toyfloat15x3 := makeTypeX3(15, true, t)

	for f := -255.0; f <= 255.0; f += 0.01 {
		toy := toyfloat15x3.Encode(f)
		original := toyfloat15x3.Decode(toy)

		if 0b1000_0000_0000_0000&toy != 0x0 {
			t.Fatalf("%.4f -> 0b%b (has extra bits)", f, toy)
		}

		m := 0b1
		modification := uint16(m) << 15
		toyModified := toy | modification
		modified := toyfloat15x3.Decode(toyModified)

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

func Test4X2Precision(t *testing.T) {
	toyfloat4x2 := makeTypeX2(4, true, t)

	a := math.Pow(3, -3)
	c := 1.0 / (1.0 - a)

	gap := 0.99
	boundary1 := (math.Pow(3, -2) - a) * c
	boundary2 := (math.Pow(3, -1) - a) * c
	boundary3 := (math.Pow(3, 0) - a) * c

	baseEps := (3.0 - 1.0) * (1.0 / math.Pow(2, 1))
	eps1 := (baseEps * math.Pow(3, -3)) * c
	eps2 := (baseEps * math.Pow(3, -2)) * c
	eps3 := (baseEps * math.Pow(3, -1)) * c
	eps4 := (baseEps * math.Pow(3, 0)) * c

	check := func(msg string, result, input, diff, eps float64) {
		if diff > eps {
			t.Fatalf("%s: %f != %f, diff: %.16f > %.16f",
				msg, result, input, diff, eps)
		}
	}

	for input := -3.0; input <= 3.0; input += 0.1 {
		tf := toyfloat4x2.Encode(input)
		result := toyfloat4x2.Decode(tf)

		diff := math.Abs(result - input)
		absInput := math.Abs(input)

		if absInput < boundary1*gap {
			check("< b1", result, input, diff, eps1)
		} else if absInput < boundary2*gap {
			check("< b2", result, input, diff, eps2)
		} else if absInput < boundary3*gap {
			check("< b3", result, input, diff, eps3)
		} else {
			check(">= b3", result, input, diff, eps4)
		}
	}
}

func Test4X2Zero(t *testing.T) {
	toyfloat4x2 := makeTypeX2(4, true, t)

	tf := toyfloat4x2.Encode(0)
	t.Logf("Encoded: 0b%b", tf)

	result := toyfloat4x2.Decode(tf)

	if result != 0 {
		t.Fatalf("%f != 0", result)
	}
}

func Test4X2PlusOne(t *testing.T) {
	toyfloat4x2 := makeTypeX2(4, true, t)

	tf := toyfloat4x2.Encode(1)
	t.Logf("Encoded: 0b%b", tf)

	result := toyfloat4x2.Decode(tf)

	if result != 1 {
		t.Fatalf("%f != 1", result)
	}
}

func Test4X2MinusOne(t *testing.T) {
	toyfloat4x2 := makeTypeX2(4, true, t)

	tf := toyfloat4x2.Encode(-1)
	t.Logf("Encoded: 0b%b", tf)

	result := toyfloat4x2.Decode(tf)

	if result != -1 {
		t.Fatalf("%f != -1", result)
	}
}

func Test4X2MaxValue(t *testing.T) {
	toyfloat4x2 := makeTypeX2(4, true, t)
	const expected = 2.0384615384615383
	const eps = 0.0001
	got := toyfloat4x2.MaxValue()
	if math.Abs(got-expected) > eps {
		t.Fatalf("got %f, expected %f", got, expected)
	}
}

func Test4X2MinValue(t *testing.T) {
	toyfloat4x2 := makeTypeX2(4, true, t)
	const expected = -2.0384615384615383
	const eps = 0.0001
	got := toyfloat4x2.MinValue()
	if math.Abs(got-expected) > eps {
		t.Fatalf("got %f, expected %f", got, expected)
	}
}

func Test4X2PositiveOverflow(t *testing.T) {
	toyfloat4x2 := makeTypeX2(4, true, t)
	const expected = 2.0384615384615383
	const eps = 0.0001

	for i := 0; i < 1000; i++ {
		v := expected + float64(i)
		tf := toyfloat4x2.Encode(v)

		result := toyfloat4x2.Decode(tf)

		if math.Abs(result-expected) > eps {
			t.Logf("Encoded: 0b%b", tf)
			t.Fatalf("%f != %f (i = %d)", result, expected, i)
		}
	}
}

func Test4X2NegativeOverflow(t *testing.T) {
	toyfloat4x2 := makeTypeX2(4, true, t)
	const expected = -2.0384615384615383
	const eps = 0.0001

	for i := 0; i < 1000; i++ {
		v := expected - float64(i)
		tf := toyfloat4x2.Encode(v)

		result := toyfloat4x2.Decode(tf)

		if math.Abs(result-expected) > eps {
			t.Logf("Encoded: 0b%b", tf)
			t.Fatalf("%f != %f (i = %d)", result, expected, i)
		}
	}
}

func Test4X2PositiveInfinity(t *testing.T) {
	toyfloat4x2 := makeTypeX2(4, true, t)
	const expected = 2.0384615384615383
	const eps = 0.0001

	v := math.Inf(+1)
	tf := toyfloat4x2.Encode(v)

	result := toyfloat4x2.Decode(tf)

	if math.Abs(result-expected) > eps {
		t.Logf("Encoded: 0b%b", tf)
		t.Fatalf("%f != %f", result, expected)
	}
}

func Test4X2NegativeInfinity(t *testing.T) {
	toyfloat4x2 := makeTypeX2(4, true, t)
	const expected = -2.0384615384615383
	const eps = 0.0001

	v := math.Inf(-1)
	tf := toyfloat4x2.Encode(v)

	result := toyfloat4x2.Decode(tf)

	if math.Abs(result-expected) > eps {
		t.Logf("Encoded: 0b%b", tf)
		t.Fatalf("%f != %f", result, expected)
	}
}

func Test4X2NaNConvertedToZero(t *testing.T) {
	toyfloat4x2 := makeTypeX2(4, true, t)

	tf := toyfloat4x2.Encode(math.NaN())
	t.Logf("Encoded: 0b%b", tf)

	result := toyfloat4x2.Decode(tf)

	if result != 0 {
		t.Fatalf("%f != 0", result)
	}
}

func Test4X2IgnoringMostSignificantBits(t *testing.T) {
	toyfloat4x2 := makeTypeX2(4, true, t)

	for f := -4.0; f <= 4.0; f += 0.01 {
		toy := toyfloat4x2.Encode(f)
		original := toyfloat4x2.Decode(toy)

		if 0xFFF0&toy != 0x0 {
			t.Fatalf("%.4f -> 0b%b (has extra bits)", f, toy)
		}

		for m := 0x1; m <= 0xFFF; m++ {
			modification := uint16(m) << 4
			toyModified := toy | modification
			modified := toyfloat4x2.Decode(toyModified)

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

func TestEncodeDecodeStability(t *testing.T) {
	toyfloat15x3 := makeTypeX3(15, true, t)

	tf := toyfloat15x3.Encode(0.6)
	t.Logf("Encoded: 0b%b", tf)

	input := toyfloat15x3.Decode(tf)
	t.Logf("Decoded: %f", input)

	temp := input

	for i := 0; i < 10; i++ {
		if temp != input {
			t.Fatalf("#%d: %f != %f", i, temp, input)
		}

		tfTemp := toyfloat15x3.Encode(temp)
		t.Logf("Encoded: 0b%b", tfTemp)
		temp = toyfloat15x3.Decode(tfTemp)
		t.Logf("Decoded: %f", temp)
	}
}

// ------------------------

func TestUseDelta15X3(t *testing.T) {
	toyfloat15x3 := makeTypeX3(15, true, t)

	const eps060 = 0.00024   // x = -1
	const eps030 = 0.00012   // x = -2
	const eps013 = 0.00006   // x = -3
	const epsMin = 0.0000075 // x = -6
	const epsMax = 0.00096   // x = 1

	a := math.Pow(2, -6)
	b := 1 / (1 - a)
	twoPowerM := math.Pow(2, 11)

	last := toyfloat15x3.Decode(toyfloat15x3.Encode(0.6))
	lastTf := toyfloat15x3.Encode(last)
	t.Logf("Last encoded: 0b%b", lastTf)

	resultTf := toyfloat15x3.UseIntegerDelta(lastTf, 0)
	result := toyfloat15x3.Decode(resultTf)
	expected := last

	if math.Abs(result-expected) > eps060 {
		t.Logf("delta=0 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = toyfloat15x3.UseIntegerDelta(lastTf, 1)
	result = toyfloat15x3.Decode(resultTf)
	expected = last + ((1.0/twoPowerM)*0.5)*b

	if math.Abs(result-expected) > eps060 {
		t.Logf("delta=1 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = toyfloat15x3.UseIntegerDelta(lastTf, -1)
	result = toyfloat15x3.Decode(resultTf)
	expected = last - ((1.0/twoPowerM)*0.5)*b

	if math.Abs(result-expected) > eps060 {
		t.Logf("delta=-1 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = toyfloat15x3.UseIntegerDelta(lastTf, 123)
	result = toyfloat15x3.Decode(resultTf)
	expected = last + ((123.0/twoPowerM)*0.5)*b

	if math.Abs(result-expected) > eps060 {
		t.Logf("delta=123 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = toyfloat15x3.UseIntegerDelta(lastTf, -123)
	result = toyfloat15x3.Decode(resultTf)
	expected = last - ((123.0/twoPowerM)*0.5)*b

	if math.Abs(result-expected) > eps060 {
		t.Logf("delta=-123 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	last = toyfloat15x3.Decode(toyfloat15x3.Encode(0.3))
	lastTf = toyfloat15x3.Encode(last)
	t.Logf("Last encoded: 0b%b", lastTf)

	mantissa := lastTf & 0b11111111111
	if mantissa != 499 {
		t.Fatalf("This test is probably broken: mantissa equals %d.", mantissa)
	}

	result = toyfloat15x3.Decode(toyfloat15x3.UseIntegerDelta(lastTf, 0))
	expected = last

	if math.Abs(result-expected) > eps030 {
		t.Fatalf("%f != %f", result, expected)
	}

	result = toyfloat15x3.Decode(toyfloat15x3.UseIntegerDelta(lastTf, 2047-499))
	expected = last + (((2047.0-499.0)/twoPowerM)*0.25)*b

	if math.Abs(result-expected) > eps030 {
		t.Fatalf("%f != %f", result, expected)
	}

	result = toyfloat15x3.Decode(toyfloat15x3.UseIntegerDelta(lastTf, -499))
	expected = last - ((499.0/twoPowerM)*0.25)*b

	if math.Abs(result-expected) > eps030 {
		t.Fatalf("%f != %f", result, expected)
	}

	result = toyfloat15x3.Decode(toyfloat15x3.UseIntegerDelta(lastTf, 2047-499+1))
	expected = last +
		(((2047.0-499.0)/twoPowerM)*0.25)*b +
		((1.0/twoPowerM)*0.5)*b

	if math.Abs(result-expected) > eps060 {
		t.Fatalf("%f != %f", result, expected)
	}

	result = toyfloat15x3.Decode(toyfloat15x3.UseIntegerDelta(lastTf, -500))
	expected = last -
		((499.0/twoPowerM)*0.25)*b -
		((1.0/twoPowerM)*0.125)*b

	if math.Abs(result-expected) > eps013 {
		t.Fatalf("%f != %f", result, expected)
	}

	lastTf = 0b0
	last = toyfloat15x3.Decode(lastTf)
	t.Logf("Last encoded: 0b%b", lastTf)

	resultTf = toyfloat15x3.UseIntegerDelta(lastTf, 0)
	result = toyfloat15x3.Decode(resultTf)
	expected = last

	if math.Abs(result-expected) > epsMin {
		t.Logf("delta=0 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = toyfloat15x3.UseIntegerDelta(lastTf, 1)
	result = toyfloat15x3.Decode(resultTf)
	expected = last + ((1.0/twoPowerM)*0.015625)*b

	if math.Abs(result-expected) > epsMin {
		t.Logf("delta=1 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = toyfloat15x3.UseIntegerDelta(lastTf, -1)
	result = toyfloat15x3.Decode(resultTf)
	expected = -last // minus zero

	if math.Abs(result-expected) > epsMin {
		t.Logf("delta=-1 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	// and back
	minusBit := toyfloat15x3.Encode(1) ^ toyfloat15x3.Encode(-1)
	t.Logf("Minus bit: 0b%b", minusBit)
	lastTf = minusBit | 0b0
	last = toyfloat15x3.Decode(lastTf)
	t.Logf("Last encoded: 0b%b", lastTf)

	resultTf = toyfloat15x3.UseIntegerDelta(lastTf, 1)
	result = toyfloat15x3.Decode(resultTf)
	expected = -last // plus zero

	if math.Abs(result-expected) > epsMin {
		t.Logf("delta=1 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = toyfloat15x3.UseIntegerDelta(lastTf, -10)
	result = toyfloat15x3.Decode(resultTf)
	expected = last - ((10.0/twoPowerM)*0.015625)*b

	if math.Abs(result-expected) > epsMin {
		t.Logf("delta=-10 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	// big deltas

	resultTf = toyfloat15x3.UseIntegerDelta(lastTf, 234234)
	result = toyfloat15x3.Decode(resultTf)
	expected = ((1.0+(twoPowerM-1)/twoPowerM)*2 - a) * b

	if math.Abs(result-expected) > epsMax {
		t.Logf("delta=234234 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = toyfloat15x3.UseIntegerDelta(lastTf, -16382)
	result = toyfloat15x3.Decode(resultTf)
	expected = -((1.0+(twoPowerM-2)/twoPowerM)*2 - a) * b

	if math.Abs(result-expected) > epsMax {
		t.Logf("delta=-16382 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = toyfloat15x3.UseIntegerDelta(lastTf, -16383)
	result = toyfloat15x3.Decode(resultTf)
	expected = -((1.0+(twoPowerM-1)/twoPowerM)*2 - a) * b

	if math.Abs(result-expected) > epsMax {
		t.Logf("delta=-16383 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = toyfloat15x3.UseIntegerDelta(lastTf, -16384)
	result = toyfloat15x3.Decode(resultTf)
	expected = -((1.0+(twoPowerM-1)/twoPowerM)*2 - a) * b

	if math.Abs(result-expected) > epsMax {
		t.Logf("delta=-16384 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = toyfloat15x3.UseIntegerDelta(lastTf, -26384)
	result = toyfloat15x3.Decode(resultTf)
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

	toyfloat15x3 := makeTypeX3(15, true, t)
	last := toyfloat15x3.Encode(start)

	for x := start + step; x <= stop; x += step {
		xtf := toyfloat15x3.Encode(x)
		expected := toyfloat15x3.Decode(xtf)

		delta := toyfloat15x3.GetIntegerDelta(last, xtf)
		resultTf := toyfloat15x3.UseIntegerDelta(last, delta)
		result := toyfloat15x3.Decode(resultTf)

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

	toyfloat12 := makeTypeX4(12, true, t)
	last := toyfloat12.Encode(start)

	for x := start + step; x <= stop; x += step {
		xtf := toyfloat12.Encode(x)
		expected := toyfloat12.Decode(xtf)

		delta := toyfloat12.GetIntegerDelta(last, xtf)
		resultTf := toyfloat12.UseIntegerDelta(last, delta)
		result := toyfloat12.Decode(resultTf)

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
	toyfloat12u := makeTypeX4(12, false, t)

	const eps060 = 0.00195  // x = -1
	const eps030 = 0.00097  // x = -2
	const eps013 = 0.00048  // x = -3
	const epsMin = 0.000015 // x = -8
	const epsMax = 0.49     // x = 7

	a := math.Pow(2, -8)
	b := 1 / (1 - a)
	twoPowerM := math.Pow(2, 8)

	last := toyfloat12u.Decode(toyfloat12u.Encode(0.6))
	lastTf := toyfloat12u.Encode(last)
	t.Logf("Last encoded: 0b%b", lastTf)

	resultTf := toyfloat12u.UseIntegerDelta(lastTf, 0)
	result := toyfloat12u.Decode(resultTf)
	expected := last

	if math.Abs(result-expected) > eps060 {
		t.Logf("delta=0 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = toyfloat12u.UseIntegerDelta(lastTf, 1)
	result = toyfloat12u.Decode(resultTf)
	expected = last + ((1.0/twoPowerM)*0.5)*b

	if math.Abs(result-expected) > eps060 {
		t.Logf("delta=1 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = toyfloat12u.UseIntegerDelta(lastTf, -1)
	result = toyfloat12u.Decode(resultTf)
	expected = last - ((1.0/twoPowerM)*0.5)*b

	if math.Abs(result-expected) > eps060 {
		t.Logf("delta=-1 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = toyfloat12u.UseIntegerDelta(lastTf, 45)
	result = toyfloat12u.Decode(resultTf)
	expected = last + ((45.0/twoPowerM)*0.5)*b

	if math.Abs(result-expected) > eps060 {
		t.Logf("delta=123 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = toyfloat12u.UseIntegerDelta(lastTf, -45)
	result = toyfloat12u.Decode(resultTf)
	expected = last - ((45.0/twoPowerM)*0.5)*b

	if math.Abs(result-expected) > eps060 {
		t.Logf("delta=-123 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	last = toyfloat12u.Decode(toyfloat12u.Encode(0.3))
	lastTf = toyfloat12u.Encode(last)
	t.Logf("Last encoded: 0b%b", lastTf)

	mantissa := lastTf & 0b11111111
	if mantissa != 54 {
		t.Fatalf("This test is probably broken: mantissa equals %d.", mantissa)
	}

	result = toyfloat12u.Decode(toyfloat12u.UseIntegerDelta(lastTf, 0))
	expected = last

	if math.Abs(result-expected) > eps030 {
		t.Fatalf("%f != %f", result, expected)
	}

	result = toyfloat12u.Decode(toyfloat12u.UseIntegerDelta(lastTf, 255-54))
	expected = last + (((255.0-54.0)/twoPowerM)*0.25)*b

	if math.Abs(result-expected) > eps030 {
		t.Fatalf("%f != %f", result, expected)
	}

	result = toyfloat12u.Decode(toyfloat12u.UseIntegerDelta(lastTf, -54))
	expected = last - ((54.0/twoPowerM)*0.25)*b

	if math.Abs(result-expected) > eps030 {
		t.Fatalf("%f != %f", result, expected)
	}

	result = toyfloat12u.Decode(toyfloat12u.UseIntegerDelta(lastTf, 255-54+1))
	expected = last +
		(((255.0-54.0)/twoPowerM)*0.25)*b +
		((1.0/twoPowerM)*0.5)*b

	if math.Abs(result-expected) > eps060 {
		t.Fatalf("%f != %f", result, expected)
	}

	result = toyfloat12u.Decode(toyfloat12u.UseIntegerDelta(lastTf, -55))
	expected = last -
		((54.0/twoPowerM)*0.25)*b -
		((1.0/twoPowerM)*0.125)*b

	if math.Abs(result-expected) > eps013 {
		t.Fatalf("%f != %f", result, expected)
	}

	lastTf = 0b0
	last = toyfloat12u.Decode(lastTf)
	t.Logf("Last encoded: 0b%b", lastTf)

	resultTf = toyfloat12u.UseIntegerDelta(lastTf, 0)
	result = toyfloat12u.Decode(resultTf)
	expected = last

	if math.Abs(result-expected) > epsMin {
		t.Logf("delta=0 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = toyfloat12u.UseIntegerDelta(lastTf, 1)
	result = toyfloat12u.Decode(resultTf)
	expected = last + ((1.0/twoPowerM)*0.00390625)*b

	if math.Abs(result-expected) > epsMin {
		t.Logf("delta=1 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = toyfloat12u.UseIntegerDelta(lastTf, -1)
	result = toyfloat12u.Decode(resultTf)
	expected = last // zero is the minimum value

	if result != expected {
		t.Logf("delta=-1 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = toyfloat12u.UseIntegerDelta(lastTf, -10)
	result = toyfloat12u.Decode(resultTf)
	expected = last // zero is the minimum value

	if result != expected {
		t.Logf("delta=-10 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	// big deltas

	resultTf = toyfloat12u.UseIntegerDelta(lastTf, 234234)
	result = toyfloat12u.Decode(resultTf)
	expected = ((1.0+(twoPowerM-1)/twoPowerM)*128 - a) * b

	if math.Abs(result-expected) > epsMax {
		t.Logf("delta=234234 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	resultTf = toyfloat12u.UseIntegerDelta(lastTf, -16382)
	result = toyfloat12u.Decode(resultTf)
	expected = 0.0

	if math.Abs(result-expected) > epsMin {
		t.Logf("delta=-16382 encoded: 0b%b", resultTf)
		t.Fatalf("%f != %f", result, expected)
	}

	last = toyfloat12u.Decode(toyfloat12u.Encode(0.234))
	lastTf = toyfloat12u.Encode(last)

	resultTf = toyfloat12u.UseIntegerDelta(lastTf, -26384)
	result = toyfloat12u.Decode(resultTf)
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

	toyfloat12u := makeTypeX4(12, false, t)
	last := toyfloat12u.Encode(start)

	for x := start + step; x <= stop; x += step {
		xtf := toyfloat12u.Encode(x)
		expected := toyfloat12u.Decode(xtf)

		delta := toyfloat12u.GetIntegerDelta(last, xtf)
		resultTf := toyfloat12u.UseIntegerDelta(last, delta)
		result := toyfloat12u.Decode(resultTf)

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

	toyfloat13 := makeTypeX4(13, true, t)
	last := toyfloat13.Encode(start)

	for x := start + step; x <= stop; x += step {
		xtf := toyfloat13.Encode(x)
		expected := toyfloat13.Decode(xtf)

		delta := toyfloat13.GetIntegerDelta(last, xtf)
		resultTf := toyfloat13.UseIntegerDelta(last, delta)
		result := toyfloat13.Decode(resultTf)

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

	toyfloat14 := makeTypeX4(14, true, t)
	last := toyfloat14.Encode(start)

	for x := start + step; x <= stop; x += step {
		xtf := toyfloat14.Encode(x)
		expected := toyfloat14.Decode(xtf)

		delta := toyfloat14.GetIntegerDelta(last, xtf)
		resultTf := toyfloat14.UseIntegerDelta(last, delta)
		result := toyfloat14.Decode(resultTf)

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
	toyfloat12 := makeTypeX4(12, true, t)

	tf := toyfloat12.Encode(42)
	t.Logf("Encoded: 0b%b", tf)

	a := toyfloat12.Decode(tf)
	b := -toyfloat12.Decode(tf | 0b1000_0000_0000)

	if a != b {
		t.Fatalf("%f != %f", a, b)
	}
}

func Test13MinusBitPosition(t *testing.T) {
	toyfloat13 := makeTypeX4(13, true, t)

	tf := toyfloat13.Encode(42)
	t.Logf("Encoded: 0b%b", tf)

	a := toyfloat13.Decode(tf)
	b := -toyfloat13.Decode(tf | 0b1_0000_0000_0000)

	if a != b {
		t.Fatalf("%f != %f", a, b)
	}
}

func Test14MinusBitPosition(t *testing.T) {
	toyfloat14 := makeTypeX4(14, true, t)

	tf := toyfloat14.Encode(42)
	t.Logf("Encoded: 0b%b", tf)

	a := toyfloat14.Decode(tf)
	b := -toyfloat14.Decode(tf | 0b10_0000_0000_0000)

	if a != b {
		t.Fatalf("%f != %f", a, b)
	}
}

func Test15X3MinusBitPosition(t *testing.T) {
	toyfloat15x3, err15x3 := NewTypeX3(15, true)
	if err15x3 != nil {
		t.Fatal(err15x3)
	}

	tf := toyfloat15x3.Encode(42)
	t.Logf("Encoded: 0b%b", tf)

	a := toyfloat15x3.Decode(tf)
	b := -toyfloat15x3.Decode(tf | 0b0100_0000_0000_0000)

	if a != b {
		t.Fatalf("%f != %f", a, b)
	}
}

func Test4X2MinusBitPosition(t *testing.T) {
	toyfloat4x2 := makeTypeX2(4, true, t)

	tf := toyfloat4x2.Encode(0.5)
	t.Logf("Encoded: 0b%b", tf)

	a := toyfloat4x2.Decode(tf)
	b := -toyfloat4x2.Decode(tf | 0b1000)

	if a != b {
		t.Fatalf("%f != %f", a, b)
	}
}

func TestReadme(t *testing.T) {
	const input = 1.567
	const eps = 1e-6

	toyfloat12 := makeTypeX4(12, true, t)
	toyfloat12u := makeTypeX4(12, false, t)
	toyfloat13 := makeTypeX4(13, true, t)
	toyfloat14 := makeTypeX4(14, true, t)

	toyfloat15x3 := makeTypeX3(15, true, t)
	toyfloat5x3 := makeTypeX3(5, true, t)
	toyfloat5x2 := makeTypeX2(5, true, t)
	toyfloat3x2u := makeTypeX2(3, false, t)

	{
		tf := toyfloat12.Encode(input)
		if tf != 0x448 {
			t.Fatalf("Incorrect encoded: 0x%X (12-bit)\n", tf)
		}

		result := toyfloat12.Decode(tf)
		if math.Abs(result-1.564706) > eps {
			t.Fatalf("Incorrect decoded: %f (12-bit)\n", result)
		}
	}

	{
		tf := toyfloat12u.Encode(input)
		if tf != 0x891 {
			t.Fatalf("Incorrect encoded: 0x%X (12-bit unsigned)\n", tf)
		}
	}

	{
		tf := toyfloat13.Encode(input)
		if tf != 0x891 {
			t.Fatalf("Incorrect encoded: 0x%X (13-bit)\n", tf)
		}

		result := toyfloat13.Decode(tf)
		if math.Abs(result-1.568627) > eps {
			t.Fatalf("Incorrect decoded: %f (13-bit)\n", result)
		}
	}

	{
		tf := toyfloat14.Encode(input)
		if tf != 0x1121 {
			t.Fatalf("Incorrect encoded: 0x%X (14-bit)\n", tf)
		}

		result := toyfloat14.Decode(tf)
		if math.Abs(result-1.566667) > eps {
			t.Fatalf("Incorrect decoded: %f (14-bit)\n", result)
		}
	}

	{
		tf := toyfloat15x3.Encode(input)
		if tf != 0x3477 {
			t.Fatalf("Incorrect encoded: 0x%X (15x3)\n", tf)
		}

		result := toyfloat15x3.Decode(tf)
		if math.Abs(result-1.566964) > eps {
			t.Fatalf("Incorrect decoded: %f (15x3)\n", result)
		}
	}

	{
		tf := toyfloat5x3.Encode(input)
		if tf != 0b01101 {
			t.Fatalf("Incorrect encoded: %05b (5x3)\n", tf)
		}

		result := toyfloat5x3.Decode(tf)
		if math.Abs(result-1.507937) > eps {
			t.Fatalf("Incorrect decoded: %f (5x3)\n", result)
		}
	}

	{
		tf := toyfloat5x2.Encode(input)
		if tf != 0b01101 {
			t.Fatalf("Incorrect encoded: %05b (5x2)\n", tf)
		}

		result := toyfloat5x2.Decode(tf)
		if math.Abs(result-1.519231) > eps {
			t.Fatalf("Incorrect decoded: %f (5x2)\n", result)
		}
	}

	{
		tf := toyfloat3x2u.Encode(input)
		if tf != 0x7 {
			t.Fatalf("Incorrect encoded: 0x%X (3x2u)\n", tf)
		}

		result := toyfloat3x2u.Decode(tf)
		if math.Abs(result-2.038462) > eps {
			t.Fatalf("Incorrect decoded: %f (3x2u)\n", result)
		}
	}

	{
		series := []float64{
			-0.0058, 0.01, -0.0058, 0.01, 0.066, 0.123,
			0.134, 0.132, 0.144, 0.145, 0.140}

		expected := []int{
			387, -387, 387, 300, 114,
			12, -2, 12, 1, -5}

		previous := toyfloat12.Encode(series[0])
		for i := 1; i < len(series); i++ {
			this := toyfloat12.Encode(series[i])

			delta := toyfloat12.GetIntegerDelta(previous, this)
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
		toyfloat13 := makeTypeX4(13, true, t)

		input := 0.9999999999995131
		result := toyfloat13.Decode(toyfloat13.Encode(input))

		if math.Abs(result-input) > 0.001 {
			t.Fatalf("%f != %f (13-bit)\n", result, input)
		}
	}

	{
		toyfloat4x2 := makeTypeX2(4, true, t)

		input := 0.9999999999995131
		result := toyfloat4x2.Decode(toyfloat4x2.Encode(input))

		if math.Abs(result-input) > 0.001 {
			t.Fatalf("%f != %f (4x2)\n", result, input)
		}
	}
}
