package coordinates

type CoordinateAxis int8

const (
	X CoordinateAxis = iota
	Y
	Z
	W
)

type Coordinate struct {
	tuple []float32
}

func (c *Coordinate) Set(axis CoordinateAxis, val float32) {
	c.tuple[axis] = val
}

func (c *Coordinate) Get(axis CoordinateAxis) float32 {
	return c.tuple[axis]
}

func (c *Coordinate) IsAPoint() bool {
	return c.tuple[W] == 1
}

func (c *Coordinate) IsAVector() bool {
	return c.tuple[W] == 0
}

func createCoordinate(x, y, z, w float32) Coordinate {
	c := Coordinate{tuple: make([]float32, 4, 4)}
	c.Set(X, x)
	c.Set(Y, y)
	c.Set(Z, z)
	c.Set(W, w)
	return c
}

func CreatePoint(x, y, z float32) Coordinate {
	return createCoordinate(x, y, z, 1)
}

func CreateVector(x, y, z float32) Coordinate {
	return createCoordinate(x, y, z, 0)
}

func (c1 *Coordinate) Add(c2 *Coordinate) *Coordinate {
	c3 := &Coordinate{tuple: make([]float32, 4, 4)}

	for i := 0; i < 4; i++ {
		c3.tuple[i] = c1.tuple[i] + c2.tuple[i]
	}

	return c3
}

func (c1 *Coordinate) Sub(c2 *Coordinate) *Coordinate {
	c3 := &Coordinate{tuple: make([]float32, 4, 4)}

	for i := 0; i < 4; i++ {
		c3.tuple[i] = c1.tuple[i] - c2.tuple[i]
	}

	return c3
}

func (c *Coordinate) Negate() *Coordinate {
	c3 := &Coordinate{tuple: make([]float32, 4, 4)}

	for i := 0; i < 4; i++ {
		c3.tuple[i] = -c.tuple[i]
	}

	return c3
}
