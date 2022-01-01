# Toyfloat

It encodes and decodes floating-point numbers with a width of 12 to 15 bits.

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
____ xxxx smmm mmmm - default
____ xxxx mmmm mmmm - unsigned
___s xxxx mmmm mmmm - 13-bit
__xx xxsm mmmm mmmm - 14-bit
_xxx smmm mmmm mmmm - m11x3

____ sxxx xmmm mmmm - dd, defaultD (D is for delta encoding)
__sx xxxm mmmm mmmm - 14d
_sxx xmmm mmmm mmmm - m11x3d
```

![Precision graph](images/comparison.png)

## Usage

```go
package main

import (
	"fmt"
	"github.com/georgy7/toyfloat"
)

func main() {
	println()

	tf := toyfloat.Encode(0.345)
	fmt.Printf("0x%X\n", tf)

	f := toyfloat.Decode(tf)
	fmt.Printf("%f\n\n", f)

	tf = toyfloat.Encode13(0.345)
	fmt.Printf("0x%X\n", tf)

	f = toyfloat.Decode13(tf)
	fmt.Printf("%f\n\n", f)

	tf = toyfloat.Encode14(0.345)
	fmt.Printf("0x%X\n", tf)

	f = toyfloat.Decode14(tf)
	fmt.Printf("%f\n\n", f)

	tf = toyfloat.EncodeM11X3(0.345)
	fmt.Printf("0x%X\n", tf)

	f = toyfloat.DecodeM11X3(tf)
	fmt.Printf("%f\n\n", f)
}
```

```shell
go get -u github.com/georgy7/toyfloat
go run example.go
```

```
0x632
0.345098

0x664
0.345098

0x18C8
0.345098

0x435E
0.344990
```
