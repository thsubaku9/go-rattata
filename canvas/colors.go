package canvas

import (
	"math"
	"rattata/rays"
)

type ColourMode uint8

const (
	Red ColourMode = iota
	Green
	Blue
)

type Colour [3]uint8

func (c *Colour) GetValue(m ColourMode) uint8 {
	return c[m]
}

func (c *Colour) SetValue(m ColourMode, v uint8) {
	c[m] = v
}

func NewColour() Colour {
	return Colour{0, 0, 0}
}

func NewWhiteColour() Colour {
	return Colour{255, 255, 255}
}

func (c1 *Colour) Add(c2 *Colour) *Colour {
	c3 := NewColour()

	for i := 0; i < 4; i++ {
		_t := int(c1[i]) + int(c2[i])
		c3[i] = uint8(min((_t), 255))
	}

	return &c3
}

func (c1 *Colour) Sub(c2 *Colour) *Colour {
	c3 := NewColour()

	for i := 0; i < 3; i++ {
		_t := int(c1[i]) - int(c2[i])
		c3[i] = uint8(max(_t, 0))
	}
	return &c3
}

func (c1 *Colour) Mul(k uint8) *Colour {
	c3 := NewColour()

	for i := 0; i < 3; i++ {
		_t := int(c1[i]) * int(k)
		c3[i] = uint8(min((_t), 255))
	}
	return &c3
}

func (c1 *Colour) Blend(c2 *Colour) *Colour {
	c3 := NewColour()

	for i := 0; i < 3; i++ {
		_t := int(c1[i]) * int(c2[i])
		c3[i] = uint8(min((_t), 255))
	}
	return &c3
}

func RayColorToCanvasColor(input rays.Colour) Colour {
	return Colour{
		uint8(math.Min(255, max(0, input[0]*255))),
		uint8(math.Min(255, max(0, input[1]*255))),
		uint8(math.Min(255, max(0, input[2]*255)))}
}
