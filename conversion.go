// Package toyfloat provides tiny (12 to 15 bits)
// floating-point number formats for serialization.
package toyfloat

import (
	"github.com/georgy7/toyfloat/internal/impl"
)

func make12() impl.Settings {
	const minus uint16 = 0b1000_0000_0000
	const mMask uint16 = 0b0000_0111_1111
	return impl.MakeSettings(7, 7, minus, mMask, impl.X4())
}

func make12Unsigned() impl.Settings {
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
	const minus uint16 = 0b10_0000_0000_0000
	const mMask uint16 = 0b00_0001_1111_1111
	return impl.MakeSettings(9, 9, minus, mMask, impl.X4())
}

func make15X3() impl.Settings {
	const minus uint16 = 0b0100_0000_0000_0000
	const mMask uint16 = 0b0000_0111_1111_1111
	return impl.MakeSettings(11, 11, minus, mMask, impl.X3())
}

// --------------

func Encode12(v float64) uint16 {
	settings := make12()
	return impl.Encode(v, &settings)
}

func Decode12(x uint16) float64 {
	settings := make12()
	return impl.Decode(x, &settings)
}

func Encode12U(v float64) uint16 {
	settings := make12Unsigned()
	return impl.Encode(v, &settings)
}

func Decode12U(x uint16) float64 {
	settings := make12Unsigned()
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

func Encode15X3(v float64) uint16 {
	settings := make15X3()
	return impl.Encode(v, &settings)
}

func Decode15X3(x uint16) float64 {
	settings := make15X3()
	return impl.Decode(x, &settings)
}

// -----------

func GetIntegerDelta12(last uint16, x uint16) int {
	settings := make12()
	return impl.EncodeDelta(last, x, &settings)
}

func UseIntegerDelta12(last uint16, delta int) uint16 {
	settings := make12()
	return impl.DecodeDelta(last, delta, &settings)
}

func GetIntegerDelta12U(last uint16, x uint16) int {
	settings := make12Unsigned()
	return impl.EncodeDelta(last, x, &settings)
}

func UseIntegerDelta12U(last uint16, delta int) uint16 {
	settings := make12Unsigned()
	return impl.DecodeDelta(last, delta, &settings)
}

func GetIntegerDelta13(last uint16, x uint16) int {
	settings := make13()
	return impl.EncodeDelta(last, x, &settings)
}

func UseIntegerDelta13(last uint16, delta int) uint16 {
	settings := make13()
	return impl.DecodeDelta(last, delta, &settings)
}

func GetIntegerDelta14(last uint16, x uint16) int {
	settings := make14()
	return impl.EncodeDelta(last, x, &settings)
}

func UseIntegerDelta14(last uint16, delta int) uint16 {
	settings := make14()
	return impl.DecodeDelta(last, delta, &settings)
}

func GetIntegerDelta15X3(last uint16, x uint16) int {
	settings := make15X3()
	return impl.EncodeDelta(last, x, &settings)
}

func UseIntegerDelta15X3(last uint16, delta int) uint16 {
	settings := make15X3()
	return impl.DecodeDelta(last, delta, &settings)
}
