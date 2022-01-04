// Package toyfloat provides tiny (4 to 16 bits)
// floating-point number formats for serialization.
package toyfloat

import (
	"github.com/georgy7/toyfloat/internal/impl"
)

// Type is a reusable immutable set of encoder settings.
type Type struct {
	settings impl.Settings
}

func NewTypeX3(length int, signed bool) (Type, error) {
	s, e := impl.NewSettings(length, impl.X3(), signed)
	return Type{s}, e
}

func NewTypeX4(length int, signed bool) (Type, error) {
	s, e := impl.NewSettings(length, impl.X4(), signed)
	return Type{s}, e
}

func (t *Type) Encode(v float64) uint16 {
	return impl.Encode(v, &t.settings)
}

func (t *Type) Decode(x uint16) float64 {
	return impl.Decode(x, &t.settings)
}

func (t *Type) GetIntegerDelta(last uint16, x uint16) int {
	return impl.EncodeDelta(last, x, &t.settings)
}

func (t *Type) UseIntegerDelta(last uint16, delta int) uint16 {
	return impl.DecodeDelta(last, delta, &t.settings)
}

// ----------------
// Deprecated API:

// Deprecated: Please use new object-oriented API. It's 4-8 times faster.
func Encode12(v float64) uint16 {
	toyfloat12, _ := NewTypeX4(12, true)
	return toyfloat12.Encode(v)
}

// Deprecated: Please use new object-oriented API. It's 4-8 times faster.
func Decode12(x uint16) float64 {
	toyfloat12, _ := NewTypeX4(12, true)
	return toyfloat12.Decode(x)
}

// Deprecated: Please use new object-oriented API. It's 4-8 times faster.
func Encode12U(v float64) uint16 {
	toyfloat12u, _ := NewTypeX4(12, false)
	return toyfloat12u.Encode(v)
}

// Deprecated: Please use new object-oriented API. It's 4-8 times faster.
func Decode12U(x uint16) float64 {
	toyfloat12u, _ := NewTypeX4(12, false)
	return toyfloat12u.Decode(x)
}

// Deprecated: Please use new object-oriented API. It's 4-8 times faster.
func Encode13(v float64) uint16 {
	toyfloat13, _ := NewTypeX4(13, true)
	return toyfloat13.Encode(v)
}

// Deprecated: Please use new object-oriented API. It's 4-8 times faster.
func Decode13(x uint16) float64 {
	toyfloat13, _ := NewTypeX4(13, true)
	return toyfloat13.Decode(x)
}

// Deprecated: Please use new object-oriented API. It's 4-8 times faster.
func Encode14(v float64) uint16 {
	toyfloat14, _ := NewTypeX4(14, true)
	return toyfloat14.Encode(v)
}

// Deprecated: Please use new object-oriented API. It's 4-8 times faster.
func Decode14(x uint16) float64 {
	toyfloat14, _ := NewTypeX4(14, true)
	return toyfloat14.Decode(x)
}

// Deprecated: Please use new object-oriented API. It's 4-8 times faster.
func Encode15X3(v float64) uint16 {
	toyfloat15X3, _ := NewTypeX3(15, true)
	return toyfloat15X3.Encode(v)
}

// Deprecated: Please use new object-oriented API. It's 4-8 times faster.
func Decode15X3(x uint16) float64 {
	toyfloat15X3, _ := NewTypeX3(15, true)
	return toyfloat15X3.Decode(x)
}

// Deprecated: Please use new object-oriented API. It's 4-8 times faster.
func GetIntegerDelta12(last uint16, x uint16) int {
	toyfloat12, _ := NewTypeX4(12, true)
	return toyfloat12.GetIntegerDelta(last, x)
}

// Deprecated: Please use new object-oriented API. It's 4-8 times faster.
func UseIntegerDelta12(last uint16, delta int) uint16 {
	toyfloat12, _ := NewTypeX4(12, true)
	return toyfloat12.UseIntegerDelta(last, delta)
}

// Deprecated: Please use new object-oriented API. It's 4-8 times faster.
func GetIntegerDelta12U(last uint16, x uint16) int {
	toyfloat12u, _ := NewTypeX4(12, false)
	return toyfloat12u.GetIntegerDelta(last, x)
}

// Deprecated: Please use new object-oriented API. It's 4-8 times faster.
func UseIntegerDelta12U(last uint16, delta int) uint16 {
	toyfloat12u, _ := NewTypeX4(12, false)
	return toyfloat12u.UseIntegerDelta(last, delta)
}

// Deprecated: Please use new object-oriented API. It's 4-8 times faster.
func GetIntegerDelta13(last uint16, x uint16) int {
	toyfloat13, _ := NewTypeX4(13, true)
	return toyfloat13.GetIntegerDelta(last, x)
}

// Deprecated: Please use new object-oriented API. It's 4-8 times faster.
func UseIntegerDelta13(last uint16, delta int) uint16 {
	toyfloat13, _ := NewTypeX4(13, true)
	return toyfloat13.UseIntegerDelta(last, delta)
}

// Deprecated: Please use new object-oriented API. It's 4-8 times faster.
func GetIntegerDelta14(last uint16, x uint16) int {
	toyfloat14, _ := NewTypeX4(14, true)
	return toyfloat14.GetIntegerDelta(last, x)
}

// Deprecated: Please use new object-oriented API. It's 4-8 times faster.
func UseIntegerDelta14(last uint16, delta int) uint16 {
	toyfloat14, _ := NewTypeX4(14, true)
	return toyfloat14.UseIntegerDelta(last, delta)
}

// Deprecated: Please use new object-oriented API. It's 4-8 times faster.
func GetIntegerDelta15X3(last uint16, x uint16) int {
	toyfloat15X3, _ := NewTypeX3(15, true)
	return toyfloat15X3.GetIntegerDelta(last, x)
}

// Deprecated: Please use new object-oriented API. It's 4-8 times faster.
func UseIntegerDelta15X3(last uint16, delta int) uint16 {
	toyfloat15X3, _ := NewTypeX3(15, true)
	return toyfloat15X3.UseIntegerDelta(last, delta)
}
