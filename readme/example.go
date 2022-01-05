package main

import (
	"fmt"
	"github.com/georgy7/toyfloat"
	"math"
	"os"
)

func exitOnError(err error) {
	if err != nil {
		println("impossible type")
		os.Exit(1)
	}
}

func main() {
	println()

	toyfloat12, err12 := toyfloat.NewTypeX4(12, true)
	exitOnError(err12)

	toyfloat13, err13 := toyfloat.NewTypeX4(13, true)
	exitOnError(err13)

	toyfloat14, err14 := toyfloat.NewTypeX4(14, true)
	exitOnError(err14)

	toyfloat15x3, err15x3 := toyfloat.NewTypeX3(15, true)
	exitOnError(err15x3)

	toyfloat5x3, err5x3 := toyfloat.NewTypeX3(5, true)
	exitOnError(err5x3)

	toyfloat5x2, err5x2 := toyfloat.NewTypeX2(5, true)
	exitOnError(err5x2)

	const v = 1.567

	fmt.Printf("Input:   %f\n\n", v)

	println("12-bit signed")
	println("-------------")
	tf := toyfloat12.Encode(v)
	fmt.Printf("Encoded: 0x%X\n", tf)
	f := toyfloat12.Decode(tf)
	fmt.Printf("Decoded: %f\n", f)
	fmt.Printf("Delta:   %f\n\n", math.Abs(f-v))

	println("13-bit signed")
	println("-------------")
	tf = toyfloat13.Encode(v)
	fmt.Printf("Encoded: 0x%X\n", tf)
	f = toyfloat13.Decode(tf)
	fmt.Printf("Decoded: %f\n", f)
	fmt.Printf("Delta:   %f\n\n", math.Abs(f-v))

	println("14-bit signed")
	println("-------------")
	tf = toyfloat14.Encode(v)
	fmt.Printf("Encoded: 0x%X\n", tf)
	f = toyfloat14.Decode(tf)
	fmt.Printf("Decoded: %f\n", f)
	fmt.Printf("Delta:   %f\n\n", math.Abs(f-v))

	println("15-bit signed with 3-bit exponent")
	println("---------------------------------")
	tf = toyfloat15x3.Encode(v)
	fmt.Printf("Encoded: 0x%X\n", tf)
	f = toyfloat15x3.Decode(tf)
	fmt.Printf("Decoded: %f\n", f)
	fmt.Printf("Delta:   %f\n\n", math.Abs(f-v))

	println("5-bit signed with 3-bit exponent")
	println("--------------------------------")
	tf = toyfloat5x3.Encode(v)
	fmt.Printf("Encoded: %05b\n", tf)
	f = toyfloat5x3.Decode(tf)
	fmt.Printf("Decoded: %f\n", f)
	fmt.Printf("Delta:   %f\n\n", math.Abs(f-v))

	println("5-bit signed with 2-bit exponent")
	println("--------------------------------")
	tf = toyfloat5x2.Encode(v)
	fmt.Printf("Encoded: %05b\n", tf)
	f = toyfloat5x2.Decode(tf)
	fmt.Printf("Decoded: %f\n", f)
	fmt.Printf("Delta:   %f\n\n", math.Abs(f-v))

	println()
	println("Delta encoding (12-bit)")
	println("-----------------------\n")

	series := []float64{
		-0.0058, 0.01, 0.066, 0.123,
		0.134, 0.132, 0.144, 0.145, 0.140}

	previous := toyfloat12.Encode(series[0])
	pDecoded := toyfloat12.Decode(previous)

	fmt.Printf("  Int. delta    Fp delta    Value\n")
	fmt.Printf("  % 35.6f\n", pDecoded)

	for i := 1; i < len(series); i++ {
		this := toyfloat12.Encode(series[i])
		delta := toyfloat12.GetIntegerDelta(previous, this)

		x := toyfloat12.Decode(this)
		fpDelta := x - pDecoded
		fmt.Printf("  %+10d    %+.6f   % .6f\n", delta, fpDelta, x)

		previous = this
		pDecoded = toyfloat12.Decode(previous)
	}

	println()
}
