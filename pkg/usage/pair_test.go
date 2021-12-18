package usage

import (
	"fmt"
	"math"
	"testing"
)

func TestPair(t *testing.T) {
	tests := []struct {
		a, b, eps float64
	}{
		{0, 0, 0.01},
		{0, 1, 0.01},
		{1, 1, 0.01},
		{1, 0, 0.01},
		{234, 128, 0.5},
		{-234, 128, 0.5},
		{234, -128, 0.5},
		{-234, -128, 0.5},
	}

	for ti, tt := range tests {
		testName := fmt.Sprintf("pair #%d", ti)
		t.Run(testName, func(t *testing.T) {
			got, e := NewPair(tt.a, tt.b)
			if e != nil {
				t.Fatal(e)
			}

			if math.Abs(got.A()-tt.a) > tt.eps {
				t.Errorf("A(): %.4f != %.4f", got.A(), tt.a)
			}

			if math.Abs(got.B()-tt.b) > tt.eps {
				t.Errorf("B(): %.4f != %.4f", got.B(), tt.b)
			}
		})
	}
}
