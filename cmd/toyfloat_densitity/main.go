package main

import (
	"fmt"
	"github.com/georgy7/toyfloat/pkg/base"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
)

func main() {
	binDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	err = os.Chdir(binDir)
	if err != nil {
		log.Fatal(err)
	}

	var counts [510]int

	for head := 0x00; head <= 0xFF; head++ {
		for tail := 0x00; tail <= 0x0F; tail++ {
			value := base.ReadFromEnd(uint8(head), uint8(tail))
			cIndex := toIndex(value)
			if (cIndex >= 0) && (cIndex < len(counts)) {
				counts[cIndex]++
			}
		}
	}

	for i := 0; i < len(counts); i++ {
		// 1 pixel = 2 units.
		// Values per integer = values per pixel / 2.
		counts[i] = int(math.Round(float64(counts[i]) * 0.5))
	}

	lm := 50
	rm := 5
	tm := 20
	bm := 40

	ySize := len(counts) / 4 * 3

	img := image.NewRGBA(image.Rect(0, 0,
		len(counts)+lm+rm,
		ySize+tm+bm))

	draw.Draw(img, img.Rect, image.White, image.Point{}, draw.Src)

	minDensityY := float64(tm + ySize)
	maxDensityY := float64(tm)

	paint := color.RGBA{0x08, 0x5E, 0xC2, 0xFF}
	grey := color.RGBA{0xE9, 0xE9, 0xE9, 0xFF}

	maxCountLog10 := -10000.0
	minCountLog10 := 10000.0
	for _, c := range counts {
		countLog10 := math.Log10(float64(c))

		if countLog10 > maxCountLog10 {
			maxCountLog10 = countLog10
		}

		if countLog10 < minCountLog10 {
			minCountLog10 = countLog10
		}
	}

	countLogDiff := maxCountLog10 - minCountLog10

	uniqueCounts := make(map[int]bool)

	for i, count := range counts {
		countLog10 := math.Log10(float64(count))

		if !uniqueCounts[count] {
			uniqueCounts[count] = true
			fmt.Printf("count: %d, log10: %.2f\n", count, countLog10)
		}

		v := (countLog10 - minCountLog10) / countLogDiff

		y := math.Round(v*maxDensityY + (1.0-v)*minDensityY)
		intY := int(y)
		img.Set(lm+i, intY, paint)

		if intY < int(minDensityY) {
			for y2 := intY + 1; y2 < int(minDensityY); y2++ {
				img.Set(lm+i, y2, grey)
			}
		}
	}

	out, err := os.Create("../images/density.png")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = png.Encode(out, img)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func toIndex(value float64) int {
	u := (value + 510) / 2
	return int(u)
}
