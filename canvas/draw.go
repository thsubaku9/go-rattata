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
