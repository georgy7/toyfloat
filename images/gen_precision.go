package main

import (
	"fmt"
	"github.com/georgy7/toyfloat"
	"math"
	"os"
)

func main() {
	output, err := os.Create("precision.tsv")
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
		v2 := toyfloat.Decode(toyfloat.Encode(value))

		if v2 > 0 {
			precision := math.Abs(value - v2)

			_, err := fmt.Fprintf(output, "%E\t%E\n", value, precision)
			if err != nil {
				panic(err)
			}
		}
	}
}
