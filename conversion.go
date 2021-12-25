// Package toyfloat provides tiny (12 to 14 bits)
// floating-point number formats for serialization.
package toyfloat

import (
	"github.com/georgy7/toyfloat/internal/impl"
)

func makeDefault() impl.Settings {
	const minus uint16 = 0b1000_0000
	const mMask uint16 = 0b0111_1111
	return impl.MakeSettings(7, 8, minus, mMask, impl.X4())
}

func makeUnsigned() impl.Settings {
	const minus uint16 = 0b0
	const mMask uint16 = 0b1111_1111
	return impl.MakeSettings(8, 8, minus, mMask, impl.X4())
}

func make13() impl.Settings {
	const minus uint16 = 0b1_0000_0000_0000
	const mMask uint16 = 0b0_0000_1111_1111
	return impl.MakeSettings(8, 8, minus, mMask, impl.X4())
}

func make14() impl.Settings {
	const minus uint16 = 0b10_0000_0000
	const mMask uint16 = 0b01_1111_1111
	return impl.MakeSettings(9, 10, minus, mMask, impl.X4())
}

func makeM11X3() impl.Settings {
	const minus uint16 = 0b1000_0000_0000
	const mMask uint16 = 0b0111_1111_1111
	return impl.MakeSettings(11, 12, minus, mMask, impl.X3())
}

func Encode(v float64) uint16 {
	settings := makeDefault()
	return impl.Encode(v, &settings)
}

func Decode(x uint16) float64 {
	settings := makeDefault()
	return impl.Decode(x, &settings)
}

func EncodeUnsigned(v float64) uint16 {
	settings := makeUnsigned()
	return impl.Encode(v, &settings)
}

func DecodeUnsigned(x uint16) float64 {
	settings := makeUnsigned()
	return impl.Decode(x, &settings)
}

func Encode13(v float64) uint16 {
	settings := make13()
	return impl.Encode(v, &settings)
}

func Decode13(x uint16) float64 {
	settings := make13()
	return impl.Decode(x, &settings)
}

func Encode14(v float64) uint16 {
	settings := make14()
	return impl.Encode(v, &settings)
}

func Decode14(x uint16) float64 {
	settings := make14()
	return impl.Decode(x, &settings)
}

func EncodeM11X3(v float64) uint16 {
	settings := makeM11X3()
	return impl.Encode(v, &settings)
}

func DecodeM11X3(x uint16) float64 {
	settings := makeM11X3()
	return impl.Decode(x, &settings)
}
