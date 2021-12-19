package usage

import (
	"fmt"
	"math"
	"testing"
)

func TestNewPair(t *testing.T) {
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
		{-235, -128, 0.5},
		{500, -128, 1.0},
		{501, -128, 1.0},
		{502, -128, 1.0},
	}

	for ti, tt := range tests {
		testName := fmt.Sprintf("NewPair #%d", ti)
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

func TestPairSetA(t *testing.T) {
	pair, e := NewPair(0, 0)
	if e != nil {
		t.Fatal(e)
	}

	tests := []struct {
		v, eps float64
	}{
		{-345, 1.0},
		{-145, 0.5},
		{-1, 0.01},
		{-0.5, 0.01},
		{-0.34, 0.01},
		{0, 0.01},
		{0.34, 0.01},
		{0.5, 0.01},
		{1, 0.01},
		{145, 0.5},
		{345, 1.0},
	}

	for ti, tt := range tests {
		testName := fmt.Sprintf("SetA #%d", ti)
		t.Run(testName, func(t *testing.T) {
			pair.SetA(tt.v)

			if math.Abs(pair.A()-tt.v) > tt.eps {
				t.Errorf("%.4f != %.4f", pair.A(), tt.v)
			}

			if math.Abs(pair.B()) > 0.01 {
				t.Errorf("changed B (%.4f)", pair.B())
			}
		})
	}
}

func TestPairSetB(t *testing.T) {
	pair, e := NewPair(0, 0)
	if e != nil {
		t.Fatal(e)
	}

	tests := []struct {
		v, eps float64
	}{
		{-345, 1.0},
		{-145, 0.5},
		{-1, 0.01},
		{-0.5, 0.01},
		{-0.34, 0.01},
		{0, 0.01},
		{0.34, 0.01},
		{0.5, 0.01},
		{1, 0.01},
		{145, 0.5},
		{345, 1.0},
	}

	for ti, tt := range tests {
		testName := fmt.Sprintf("SetA #%d", ti)
		t.Run(testName, func(t *testing.T) {
			pair.SetB(tt.v)

			if math.Abs(pair.B()-tt.v) > tt.eps {
				t.Errorf("%.4f != %.4f", pair.B(), tt.v)
			}

			if math.Abs(pair.A()) > 0.01 {
				t.Errorf("changed A (%.4f)", pair.A())
			}
		})
	}
}
