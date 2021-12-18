package base

import (
	impl "github.com/georgy7/toyfloat/internal/base"
)

func WriteHead(v float64) uint8 {
	return impl.GetSignificand(v)
}

func WriteStartOfTail(v float64) uint8 {
	return impl.GetExponentAsANibble(v) << 4
}

func WriteEndOfTail(v float64) uint8 {
	return impl.GetExponentAsANibble(v)
}

func ReadFromStart(head uint8, tail uint8) float64 {
	return impl.Read(head, tail>>4)
}

func ReadFromEnd(head uint8, tail uint8) float64 {
	return impl.Read(head, tail&0x0F)
}
