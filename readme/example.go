package main

import (
	"fmt"
	"github.com/georgy7/toyfloat"
	"math"
	"os"
)

func exitOnError(err error) {
	if err != nil {
		fmt.Println("impossible type")
		os.Exit(1)
	}
}

func printHeader(header string) {
	fmt.Println(header)
	for i := 0; i < len(header); i++ {
		fmt.Print("-")
	}
	fmt.Println()
}

func report(t toyfloat.Type, v float64) {
	tf := t.Encode(v)
	f := t.Decode(tf)
	fmt.Printf("Input:   %f\n", v)
	fmt.Printf("Encoded: 0x%X\n", tf)
	fmt.Printf("Decoded: %f\n", f)
	fmt.Printf("Delta:   %f\n", math.Abs(f-v))
	fmt.Printf("RE:      %f\n\n", math.Abs((f-v)/v))
}

func main() {
	fmt.Println()

	toyfloat12, err12 := toyfloat.NewTypeX4(12, true)
	toyfloat5x3, err5x3 := toyfloat.NewTypeX3(5, true)
	toyfloat3x2u, err3x2u := toyfloat.NewTypeX2(3, false)
	d8x3, errD8x3 := toyfloat.NewType(8, 10, 3, -2, true)

	exitOnError(err12)
	exitOnError(err5x3)
	exitOnError(err3x2u)
	exitOnError(errD8x3)

	const input = 1.567

	printHeader("12-bit signed")
	report(toyfloat12, input)

	printHeader("5-bit signed with 3-bit exponent")
	report(toyfloat5x3, input)

	printHeader("3-bit unsigned with 2-bit exponent")
	report(toyfloat3x2u, input)

	printHeader("8-bit with base 10 exponent (-2..5)")
	report(d8x3, input/10)
	report(d8x3, input)
	report(d8x3, 65536)
	report(d8x3, -7.5e5)
	report(d8x3, 1e6)

	fmt.Println()
	fmt.Println("Delta encoding (12-bit)")
	fmt.Println("-----------------------")
	fmt.Println()

	series := []float64{
		-0.0058, 0.01, -0.0058, 0.01, 0.066, 0.123,
		0.134, 0.132, 0.144, 0.145, 0.140}

	previous := toyfloat12.Encode(series[0])
	cPrevious := toyfloat12.ToComparable(previous)

	pDecoded := toyfloat12.Decode(previous)

	fmt.Printf("  Int. delta    Unsigned    Fp delta    Value\n")
	fmt.Printf("  % 47.6f\n", pDecoded)

	for i := 1; i < len(series); i++ {
		this := toyfloat12.Encode(series[i])
		cThis := toyfloat12.ToComparable(this)

		delta := toyfloat12.GetIntegerDelta(previous, this)
		unsignedDelta := cThis - cPrevious

		x := toyfloat12.Decode(this)
		fpDelta := x - pDecoded
		fmt.Printf("  %+10d    %8d    %+.6f   % .6f\n",
			delta, unsignedDelta, fpDelta, x)

		previous = this
		cPrevious = toyfloat12.ToComparable(previous)

		pDecoded = toyfloat12.Decode(previous)
	}

	fmt.Println()
}
