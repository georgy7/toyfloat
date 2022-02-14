# Toyfloat

It encodes and decodes floating-point numbers with a width of 3 to 16 bits.

Expected applications:

* file format design,
* lossy compression.

It has:

* exact 0, 1, -1
* no NaN, -Inf, +Inf
* values, that are in range about:
  * (-256, +256) for 4-bit exponent
  * (-4, +4) for 3-bit exponent
  * (-3, +3) for 2-bit exponent
  * up to 10^308 (custom settings).

![Formula](images/formula.png)

Base 3 in 2-bit exponent provides a higher density
of values close to zero and a wider range,
at the cost of reduced precision of values greater than one-third.

You can also choose other settings.

```
Examples:

____ sxxx xmmm mmmm - 12-bit
____ xxxx mmmm mmmm - 12-bit unsigned
___s xxxx mmmm mmmm - 13-bit
__sx xxxm mmmm mmmm - 14-bit
_sxx xmmm mmmm mmmm - 15-bit with 3-bit exponent
```

![Precision graph](images/comparison.png)

## Usage

```go
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

func report(t toyfloat.Type, value float64) {
	encoded := t.Encode(value)
	comparable := t.ToComparable(encoded)
	decoded := t.Decode(encoded)
	fmt.Printf("Input:      %f\n", value)
	fmt.Printf("Encoded:    0x%X\n", encoded)
	fmt.Printf("Comparable: 0x%X\n", comparable)
	fmt.Printf("Decoded:    %f\n", decoded)
	fmt.Printf("Error:      %f\n", math.Abs(value-decoded))
	fmt.Printf("Relative:   %f\n\n", math.Abs((value-decoded)/value))
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
```

```shell
go get -u github.com/georgy7/toyfloat
go run example.go
```

```
12-bit signed
-------------
Input:      1.567000
Encoded:    0x448
Comparable: 0xC48
Decoded:    1.564706
Error:      0.002294
Relative:   0.001464

5-bit signed with 3-bit exponent
--------------------------------
Input:      1.567000
Encoded:    0xD
Comparable: 0x1D
Decoded:    1.507937
Error:      0.059063
Relative:   0.037692

3-bit unsigned with 2-bit exponent
----------------------------------
Input:      1.567000
Encoded:    0x7
Comparable: 0x7
Decoded:    2.038462
Error:      0.471462
Relative:   0.300869

8-bit with base 10 exponent (-2..5)
-----------------------------------
Input:      0.156700
Encoded:    0x11
Comparable: 0x91
Decoded:    0.147727
Error:      0.008973
Relative:   0.057261

Input:      1.567000
Encoded:    0x21
Comparable: 0xA1
Decoded:    1.568182
Error:      0.001182
Relative:   0.000754

Input:      65536.000000
Encoded:    0x6A
Comparable: 0xEA
Decoded:    66919.181818
Error:      1383.181818
Relative:   0.021106

Input:      -750000.000000
Encoded:    0xFB
Comparable: 0x4
Decoded:    -726010.090909
Error:      23989.909091
Relative:   0.031987

Input:      1000000.000000
Encoded:    0x7F
Comparable: 0xFF
Decoded:    953282.818182
Error:      46717.181818
Relative:   0.046717


Delta encoding (12-bit)
-----------------------

  Int. delta    Unsigned    Fp delta    Value
                                        -0.005821
        +387         387    +0.015809    0.009988
        -387       65149    -0.015809   -0.005821
        +387         387    +0.015809    0.009988
        +300         300    +0.056189    0.066176
        +114         114    +0.056373    0.122549
         +12          12    +0.011765    0.134314
          -2       65534    -0.001961    0.132353
         +12          12    +0.011765    0.144118
          +1           1    +0.000980    0.145098
          -5       65531    -0.004902    0.140196
```

## Performance

```
BenchmarkFloat64IncrementAsAReference-8     1000000000           1.09 ns/op
BenchmarkCreateTypeX4-8                      3917901           299 ns/op
BenchmarkCreateTypeX3-8                     11907438           103 ns/op
BenchmarkCreateTypeX2-8                     13006556            92.4 ns/op
BenchmarkEncode-8                           90582355            13.2 ns/op
BenchmarkDecode-8                           229913229            5.15 ns/op
BenchmarkEncode12X2-8                       94929739            12.4 ns/op
BenchmarkDecode12X2-8                       236135440            5.40 ns/op
BenchmarkGetDelta-8                         540333522            2.19 ns/op
BenchmarkGetDeltaX2-8                       603912135            1.86 ns/op
BenchmarkUseDelta-8                         263490146            4.53 ns/op
BenchmarkUseDeltaX2-8                       249798333            4.54 ns/op
```
