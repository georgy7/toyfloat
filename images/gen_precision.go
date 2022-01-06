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

func exitOnError(err error) {
	if err != nil {
		println("impossible type")
		os.Exit(1)
	}
}

func main() {
	toyfloat12, err12 := toyfloat.NewTypeX4(12, true)
	exitOnError(err12)

	toyfloat12u, err12u := toyfloat.NewTypeX4(12, false)
	exitOnError(err12u)

	toyfloat13, err13 := toyfloat.NewTypeX4(13, true)
	exitOnError(err13)

	toyfloat14, err14 := toyfloat.NewTypeX4(14, true)
	exitOnError(err14)

	toyfloat15x3, err15x3 := toyfloat.NewTypeX3(15, true)
	exitOnError(err15x3)

	toyfloat15x2, err15x2 := toyfloat.NewTypeX2(15, true)
	exitOnError(err15x2)

	toyfloat8x3, err8x3 := toyfloat.NewTypeX3(8, true)
	exitOnError(err8x3)

	toyfloat4x3u, err4x3u := toyfloat.NewTypeX3(4, false)
	exitOnError(err4x3u)

	toyfloat3x2u, err3x2u := toyfloat.NewTypeX2(3, false)
	exitOnError(err3x2u)

	toyfloat16, err16 := toyfloat.NewTypeX4(16, true)
	exitOnError(err16)

	toyfloat16x3, err16x3 := toyfloat.NewTypeX3(16, true)
	exitOnError(err16x3)

	toyfloat16x2, err16x2 := toyfloat.NewTypeX2(16, true)
	exitOnError(err16x2)

	toyfloat16u, err16u := toyfloat.NewTypeX4(16, false)
	exitOnError(err16u)

	toyfloat16x3u, err16x3u := toyfloat.NewTypeX3(16, false)
	exitOnError(err16x3u)

	generate("precision12.tsv", func(x float64) float64 {
		return toyfloat12.Decode(toyfloat12.Encode(x))
	})

	generate("precision12u.tsv", func(x float64) float64 {
		return toyfloat12u.Decode(toyfloat12u.Encode(x))
	})

	generate("precision13.tsv", func(x float64) float64 {
		return toyfloat13.Decode(toyfloat13.Encode(x))
	})

	generate("precision14.tsv", func(x float64) float64 {
		return toyfloat14.Decode(toyfloat14.Encode(x))
	})

	generate("precision15x3.tsv", func(x float64) float64 {
		return toyfloat15x3.Decode(toyfloat15x3.Encode(x))
	})

	generate("precision15x2.tsv", func(x float64) float64 {
		return toyfloat15x2.Decode(toyfloat15x2.Encode(x))
	})

	generate("precision8x3.tsv", func(x float64) float64 {
		return toyfloat8x3.Decode(toyfloat8x3.Encode(x))
	})

	generate("precision4x3u.tsv", func(x float64) float64 {
		return toyfloat4x3u.Decode(toyfloat4x3u.Encode(x))
	})

	generate("precision3x2u.tsv", func(x float64) float64 {
		return toyfloat3x2u.Decode(toyfloat3x2u.Encode(x))
	})

	generate("precision16.tsv", func(x float64) float64 {
		return toyfloat16.Decode(toyfloat16.Encode(x))
	})

	generate("precision16x3.tsv", func(x float64) float64 {
		return toyfloat16x3.Decode(toyfloat16x3.Encode(x))
	})

	generate("precision16x2.tsv", func(x float64) float64 {
		return toyfloat16x2.Decode(toyfloat16x2.Encode(x))
	})

	generate("precision16u.tsv", func(x float64) float64 {
		return toyfloat16u.Decode(toyfloat16u.Encode(x))
	})

	generate("precision16x3u.tsv", func(x float64) float64 {
		return toyfloat16x3u.Decode(toyfloat16x3u.Encode(x))
	})
}
