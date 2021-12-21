package toyfloat

import (
	"github.com/georgy7/toyfloat/internal/impl"
)

func Encode(v float64) uint16 {
	const minus uint16 = 0b1000_0000
	const mMask uint16 = 0b0111_1111
	return impl.Encode(v, 7, 8, minus, mMask)
}

func Decode(x uint16) float64 {
	const minus uint16 = 0b1000_0000
	const mMask uint16 = 0b0111_1111
	return impl.Decode(x, 7, 8, minus, mMask)
}
