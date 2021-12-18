package usage

import (
	"errors"
	"github.com/georgy7/toyfloat/pkg/base"
	"math"
)

type Pair [3]uint8

func NewPair(a float64, b float64) (Pair, error) {
	var result Pair
	if math.IsNaN(a) || math.IsNaN(b) {
		return result, errors.New("NaN is not supported")
	}

	result[0] = base.WriteHead(a)
	result[1] = base.WriteHead(b)
	result[2] = base.WriteStartOfTail(a) | base.WriteEndOfTail(b)

	return result, nil
}

func (p Pair) A() float64 {
	return base.ReadFromStart(p[0], p[2])
}

func (p Pair) B() float64 {
	return base.ReadFromEnd(p[1], p[2])
}

// Using setters doesn't make much sense here.
// The struct is tiny and easily re-created.
// But as an example...

func (p *Pair) SetA(a float64) {
	p[0] = base.WriteHead(a)
	tail := p[2] & 0x0F
	p[2] = base.WriteStartOfTail(a) | tail
}

func (p *Pair) SetB(b float64) {
	p[1] = base.WriteHead(b)
	tail := p[2] & 0xF0
	p[2] = tail | base.WriteEndOfTail(b)
}
