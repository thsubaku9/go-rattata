package canvas

type ColourMode uint8

const (
	Red ColourMode = iota
	Green
	Blue
	Alpha
)

type Colour [4]uint8

func (c *Colour) GetValue(m ColourMode) uint8 {
	return c[m]
}

func (c *Colour) SetValue(m ColourMode, v uint8) {
	c[m] = v
}

func NewColour() Colour {
	return Colour{}
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

	c3[Alpha] = c1[Alpha]
	return &c3
}

func (c1 *Colour) Mul(k uint8) *Colour {
	c3 := NewColour()

	for i := 0; i < 3; i++ {
		_t := int(c1[i]) * int(k)
		c3[i] = uint8(min((_t), 255))
	}
	c3[Alpha] = c1[Alpha]
	return &c3
}

func (c1 *Colour) Blend(c2 *Colour) *Colour {
	c3 := NewColour()

	for i := 0; i < 3; i++ {
		_t := int(c1[i]) * int(c2[i])
		c3[i] = uint8(min((_t), 255))
	}
	c3[Alpha] = c1[Alpha]
	return &c3
}
