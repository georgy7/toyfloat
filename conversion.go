package toyfloat

import (
	"github.com/georgy7/toyfloat/internal/impl"
)

func Encode(v float64) uint16 {
	return impl.Encode(v)
}

func Decode(x uint16) float64 {
	return impl.Decode(x)
}
