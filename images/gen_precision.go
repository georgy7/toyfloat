package main

import (
	"fmt"
	"github.com/georgy7/toyfloat"
	"math"
	"os"
)

func generate(fn string, encodeDecode func(x float64) float64) {
	output, err := os.Create(fn)
	if err != nil {
		panic(err)
	}

	defer func() {
		err := output.Close()
		if err != nil {
			panic(err)
		}
	}()

	for p := float64(-12); p < 6; p += 0.01 {
		value := math.Pow(10, p)
		v2 := encodeDecode(value)

		if v2 > 0 {
			precision := math.Abs(value - v2)

			_, err := fmt.Fprintf(output, "%E\t%E\n", value, precision)
			if err != nil {
				panic(err)
			}
		}
	}
}

func main() {
	generate("precision.tsv", func(x float64) float64 {
		return toyfloat.Decode(toyfloat.Encode(x))
	})

	generate("precision_unsigned.tsv", func(x float64) float64 {
		return toyfloat.DecodeUnsigned(toyfloat.EncodeUnsigned(x))
	})

	generate("precision13.tsv", func(x float64) float64 {
		return toyfloat.Decode13(toyfloat.Encode13(x))
	})

	generate("precision14.tsv", func(x float64) float64 {
		return toyfloat.Decode14(toyfloat.Encode14(x))
	})

	generate("precision_m11x3.tsv", func(x float64) float64 {
		return toyfloat.DecodeM11X3(toyfloat.EncodeM11X3(x))
	})

	generate("precision_dd.tsv", func(x float64) float64 {
		return toyfloat.DecodeDD(toyfloat.EncodeDD(x))
	})

	generate("precision14d.tsv", func(x float64) float64 {
		return toyfloat.Decode14D(toyfloat.Encode14D(x))
	})

	generate("precision_m11x3d.tsv", func(x float64) float64 {
		return toyfloat.DecodeM11X3D(toyfloat.EncodeM11X3D(x))
	})
}
