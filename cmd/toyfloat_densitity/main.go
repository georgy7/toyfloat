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

	var counts [1020]int

	for head := 0x00; head <= 0xFF; head++ {
		for tail := 0x00; tail <= 0x0F; tail++ {
			value := base.ReadFromEnd(uint8(head), uint8(tail))
			cIndex := toIndex(value)
			if (cIndex >= 0) && (cIndex < len(counts)) {
				counts[cIndex]++
			}
		}
	}

	lm := 50
	rm := 5
	tm := 20
	bm := 40

	xSize := len(counts) / 2
	ySize := xSize / 4 * 3

	img := image.NewRGBA(image.Rect(0, 0,
		xSize+lm+rm,
		ySize+tm+bm))

	draw.Draw(img, img.Rect, image.White, image.Point{}, draw.Src)

	minDensityY := float64(tm + ySize)
	maxDensityY := float64(tm)

	paint := color.RGBA{0x08, 0x5E, 0xC2, 0xFF}
	grey := color.RGBA{0xE9, 0xE9, 0xE9, 0xFF}

	maxCountLog10 := -10000.0
	minCountLog10 := 10000.0

	for i := 0; i < len(counts); i += 2 {
		cAverage := (float64(counts[i]) + float64(counts[i+1])) / 2
		countLog10 := math.Log10(cAverage)

		if countLog10 > maxCountLog10 {
			maxCountLog10 = countLog10
		}

		if countLog10 < minCountLog10 {
			minCountLog10 = countLog10
		}
	}

	countLogDiff := maxCountLog10 - minCountLog10

	uniqueCounts := make(map[float64]bool)

	for i := 0; i < len(counts); i += 2 {
		count := (float64(counts[i]) + float64(counts[i+1])) / 2
		countLog10 := math.Log10(float64(count))

		if !uniqueCounts[count] {
			uniqueCounts[count] = true
			fmt.Printf("count: %.1f, log10: %.2f\n", count, countLog10)
		}

		v := (countLog10 - minCountLog10) / countLogDiff

		y := math.Round(v*maxDensityY + (1.0-v)*minDensityY)
		intY := int(y)
		x := lm + i/2
		img.Set(x, intY, paint)

		bGrey := bm / 3
		if intY < int(minDensityY)+bGrey {
			for y2 := intY + 1; y2 < int(minDensityY)+bGrey; y2++ {
				img.Set(x, y2, grey)
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
	u := value + 510
	return int(u)
}
