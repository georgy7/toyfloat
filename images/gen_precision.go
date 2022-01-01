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
	generate("precision12.tsv", func(x float64) float64 {
		return toyfloat.Decode12(toyfloat.Encode12(x))
	})

	generate("precision12u.tsv", func(x float64) float64 {
		return toyfloat.Decode12U(toyfloat.Encode12U(x))
	})

	generate("precision13.tsv", func(x float64) float64 {
		return toyfloat.Decode13(toyfloat.Encode13(x))
	})

	generate("precision14.tsv", func(x float64) float64 {
		return toyfloat.Decode14(toyfloat.Encode14(x))
	})

	generate("precision15x3.tsv", func(x float64) float64 {
		return toyfloat.Decode15X3(toyfloat.Encode15X3(x))
	})
}
