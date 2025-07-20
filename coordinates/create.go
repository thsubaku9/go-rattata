package coordinates

import "math"

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

func (c *Coordinate) Mul(f float32) *Coordinate {
	c3 := &Coordinate{tuple: make([]float32, 4, 4)}

	for i := 0; i < 4; i++ {
		c3.tuple[i] = c.tuple[i] * f
	}

	return c3
}

func (c *Coordinate) Div(f float32) *Coordinate {
	c3 := &Coordinate{tuple: make([]float32, 4, 4)}

	for i := 0; i < 4; i++ {
		c3.tuple[i] = c.tuple[i] / f
	}

	return c3
}

func (c *Coordinate) Magnitude() float64 {
	_mag := float64(0)

	for i := 0; i < 4; i++ {
		_mag += math.Pow(float64(c.tuple[i]), 2)
	}

	return math.Sqrt(float64(_mag))
}

func (c *Coordinate) Norm() *Coordinate {
	mag := c.Magnitude()

	return c.Div(float32(mag))
}

func (c1 *Coordinate) DotP(c2 *Coordinate) *Coordinate {
	c3 := &Coordinate{tuple: make([]float32, 4, 4)}

	for i := 0; i < 4; i++ {
		c3.tuple[i] = c1.tuple[i] * c2.tuple[i]
	}
	return c3
}

func (c1 *Coordinate) CrossP(c2 *Coordinate) *Coordinate {
	c3 := &Coordinate{tuple: make([]float32, 4, 4)}

	c3.Set(X, c1.Get(Y)*c2.Get(Z)-c1.Get(Z)*c2.Get(Y))
	c3.Set(Y, c1.Get(Z)*c2.Get(X)-c1.Get(X)*c2.Get(Z))
	c3.Set(Z, c1.Get(X)*c2.Get(Y)-c1.Get(Y)*c2.Get(X))
	return c3
}
