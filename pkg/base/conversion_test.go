package base

import (
	"math"
	"testing"
)

func TestExactRepresentationOfIntegers(t *testing.T) {
	const eps = 0.01
	var head, tail uint8

	for intValue := -256; intValue < 256; intValue++ {
		v := float64(intValue)
		head = WriteHead(v)
		tail = WriteStartOfTail(v)
		got := ReadFromStart(head, tail)

		if math.Abs(got-v) > eps {
			t.Fatalf("%.3f != %.3f", got, v)
		}
	}
}
