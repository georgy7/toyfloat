# Toyfloat

It encodes and decodes floating-point numbers with a width of 4 to 16 bits.

Expected applications:

* file format design,
* lossy compression.

It has:

* exact 0, 1, -1
* no NaN
* values, that are in range about:
  * (-256, +256) for 4-bit exponent
  * (-4, +4) for 3-bit exponent

![Formula](images/formula.png)

```
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

	tf := toyfloat12.Encode(0.345)
	fmt.Printf("0x%X\n", tf)

	f := toyfloat12.Decode(tf)
	fmt.Printf("%f\n\n", f)

	tf = toyfloat13.Encode(0.345)
	fmt.Printf("0x%X\n", tf)

	f = toyfloat13.Decode(tf)
	fmt.Printf("%f\n\n", f)

	tf = toyfloat14.Encode(0.345)
	fmt.Printf("0x%X\n", tf)

	f = toyfloat14.Decode(tf)
	fmt.Printf("%f\n\n", f)

	tf = toyfloat15x3.Encode(0.345)
	fmt.Printf("0x%X\n", tf)

	f = toyfloat15x3.Decode(tf)
	fmt.Printf("%f\n\n", f)

	series := []float64{-0.0058, 0.01, 0.123, 0.134, 0.132, 0.144, 0.145, 0.140}
	previous := toyfloat12.Encode(series[0])
	for i := 1; i < len(series); i++ {
		this := toyfloat12.Encode(series[i])
		delta := toyfloat12.GetIntegerDelta(previous, this)
		fmt.Printf("%d\n", delta)
		previous = this
	}
}
```

```shell
go get -u github.com/georgy7/toyfloat
go run example.go
```

```
0x332
0.345098

0x664
0.345098

0xCC8
0.345098

0x235E
0.344990

387
414
12
-2
12
1
-5
```
