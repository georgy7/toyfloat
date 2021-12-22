package toyfloat

import (
	"github.com/georgy7/toyfloat/internal/impl"
)

const minus uint16 = 0b1000_0000
const mMask uint16 = 0b0111_1111

func Encode(v float64) uint16 {
	return impl.Encode(v, 7, 8, minus, mMask)
}

func Decode(x uint16) float64 {
	return impl.Decode(x, 7, 8, minus, mMask)
}

const minusUnsigned uint16 = 0b0
const mMaskUnsigned uint16 = 0b1111_1111

func EncodeUnsigned(v float64) uint16 {
	return impl.Encode(v, 8, 8, minusUnsigned, mMaskUnsigned)
}

func DecodeUnsigned(x uint16) float64 {
	return impl.Decode(x, 8, 8, minusUnsigned, mMaskUnsigned)
}

const minus13 uint16 = 0b1_0000_0000_0000
const mMask13 uint16 = 0b0_0000_1111_1111

func Encode13(v float64) uint16 {
	return impl.Encode(v, 8, 8, minus13, mMask13)
}

func Decode13(x uint16) float64 {
	return impl.Decode(x, 8, 8, minus13, mMask13)
}

const minus14 uint16 = 0b10_0000_0000
const mMask14 uint16 = 0b01_1111_1111

func Encode14(v float64) uint16 {
	return impl.Encode(v, 9, 10, minus14, mMask14)
}

func Decode14(x uint16) float64 {
	return impl.Decode(x, 9, 10, minus14, mMask14)
}
