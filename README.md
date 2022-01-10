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

![Formula](images/formula.png)

```
Examples:

____ sxxx xmmm mmmm - 12-bit
____ xxxx mmmm mmmm - 12-bit unsigned
___s xxxx mmmm mmmm - 13-bit
__sx xxxm mmmm mmmm - 14-bit
_sxx xmmm mmmm mmmm - 15-bit with 3-bit exponent
```

![Precision graph](images/comparison.png)

Base 3 in 2-bit exponent provides a higher density
of values close to zero and a wider range,
at the cost of reduced precision of values greater than one-third.

![Base 2 and base 3 exponent comparison](images/precision15x2_b2b3.png)

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
		println("impossible type")
		os.Exit(1)
	}
}

func report(header string, tf uint16, f, v float64) {
	println(header)
	for i := 0; i < len(header); i++ {
		print("-")
	}
	println()
	fmt.Printf("Encoded: 0x%X\n", tf)
	fmt.Printf("Decoded: %f\n", f)
	fmt.Printf("Delta:   %f\n\n", math.Abs(f-v))
}

func main() {
	println()

	toyfloat12, err12 := toyfloat.NewTypeX4(12, true)
	toyfloat5x3, err5x3 := toyfloat.NewTypeX3(5, true)
	toyfloat3x2u, err3x2u := toyfloat.NewTypeX2(3, false)

	exitOnError(err12)
	exitOnError(err5x3)
	exitOnError(err3x2u)

	const input = 1.567
	fmt.Printf("Input:   %f\n\n", input)

	tf := toyfloat12.Encode(input)
	f := toyfloat12.Decode(tf)
	report("12-bit signed", tf, f, input)

	tf = toyfloat5x3.Encode(input)
	f = toyfloat5x3.Decode(tf)
	report("5-bit signed with 3-bit exponent", tf, f, input)

	tf = toyfloat3x2u.Encode(input)
	f = toyfloat3x2u.Decode(tf)
	report("3-bit unsigned with 2-bit exponent", tf, f, input)

	println()
	println("Delta encoding (12-bit)")
	println("-----------------------\n")

	series := []float64{
		-0.0058, 0.01, -0.0058, 0.01, 0.066, 0.123,
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
```

```shell
go get -u github.com/georgy7/toyfloat
go run example.go
```

```
Input:   1.567000

12-bit signed
-------------
Encoded: 0x448
Decoded: 1.564706
Delta:   0.002294

5-bit signed with 3-bit exponent
--------------------------------
Encoded: 0xD
Decoded: 1.507937
Delta:   0.059063

3-bit unsigned with 2-bit exponent
----------------------------------
Encoded: 0x7
Decoded: 2.038462
Delta:   0.471462


Delta encoding (12-bit)
-----------------------

  Int. delta    Fp delta    Value
                            -0.005821
        +387    +0.015809    0.009988
        -387    -0.015809   -0.005821
        +387    +0.015809    0.009988
        +300    +0.056189    0.066176
        +114    +0.056373    0.122549
         +12    +0.011765    0.134314
          -2    -0.001961    0.132353
         +12    +0.011765    0.144118
          +1    +0.000980    0.145098
          -5    -0.004902    0.140196
```

## Performance

```
BenchmarkFloat64IncrementAsAReference-8     1000000000           0.274 ns/op
BenchmarkCreateTypeX4-8                     24675729            47.6 ns/op
BenchmarkCreateTypeX3-8                     31378699            36.5 ns/op
BenchmarkCreateTypeX2-8                     37989314            31.8 ns/op
BenchmarkEncode-8                           100000000           11.1 ns/op
BenchmarkDecode-8                           274288977            4.42 ns/op
BenchmarkEncode12X2-8                       100000000           10.6 ns/op
BenchmarkDecode12X2-8                       290536574            4.15 ns/op
BenchmarkGetDelta-8                         304114786            4.21 ns/op
BenchmarkGetDeltaX2-8                       287911833            4.15 ns/op
BenchmarkUseDelta-8                         271670265            4.41 ns/op
BenchmarkUseDeltaX2-8                       246510997            4.92 ns/op
```
