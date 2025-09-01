package coordinates

import "math"

type CoordinateAxis int8

const (
	X CoordinateAxis = iota
	Y
	Z
	W
)

type Coordinate [4]float32

func (c *Coordinate) Set(axis CoordinateAxis, val float32) {
	c[axis] = val
}

func (c *Coordinate) Get(axis CoordinateAxis) float32 {
	return c[axis]
}

func (c *Coordinate) IsAPoint() bool {
	return c[W] == 1
}

func (c *Coordinate) IsAVector() bool {
	return c[W] == 0
}

func CreateCoordinate(x, y, z, w float32) Coordinate {
	c := Coordinate{}
	c.Set(X, x)
	c.Set(Y, y)
	c.Set(Z, z)
	c.Set(W, w)
	return c
}

func CreatePoint(x, y, z float32) Coordinate {
	return CreateCoordinate(x, y, z, 1)
}

func CreateVector(x, y, z float32) Coordinate {
	return CreateCoordinate(x, y, z, 0)
}

func (c1 *Coordinate) Add(c2 *Coordinate) *Coordinate {
	c3 := &Coordinate{}

	for i := 0; i < 4; i++ {
		c3[i] = c1[i] + c2[i]
	}

	return c3
}

func (c1 *Coordinate) Sub(c2 *Coordinate) *Coordinate {
	c3 := &Coordinate{}

	for i := 0; i < 4; i++ {
		c3[i] = c1[i] - c2[i]
	}

	return c3
}

func (c *Coordinate) Negate() *Coordinate {
	c3 := &Coordinate{}

	for i := 0; i < 4; i++ {
		c3[i] = -c[i]
	}

	return c3
}

func (c *Coordinate) Mul(f float32) *Coordinate {
	c3 := &Coordinate{}

	for i := 0; i < 4; i++ {
		c3[i] = c[i] * f
	}

	return c3
}

func (c *Coordinate) Div(f float32) *Coordinate {
	c3 := &Coordinate{}

	for i := 0; i < 4; i++ {
		c3[i] = c[i] / f
	}

	return c3
}

func (c *Coordinate) Magnitude() float32 {
	_mag := float64(0)

	for i := 0; i < 3; i++ {
		_mag += math.Pow(float64(c[i]), 2)
	}

	return float32(math.Sqrt(_mag))
}

func (c *Coordinate) Norm() *Coordinate {
	return c.Div(c.Magnitude())
}

func (c1 *Coordinate) DotP(c2 *Coordinate) float32 {
	c3 := float32(0)

	for i := 0; i < 4; i++ {
		c3 += c1[i] * c2[i]
	}
	return c3
}

func (c1 *Coordinate) CrossP(c2 *Coordinate) *Coordinate {
	c3 := &Coordinate{}

	c3.Set(X, c1.Get(Y)*c2.Get(Z)-c1.Get(Z)*c2.Get(Y))
	c3.Set(Y, c1.Get(Z)*c2.Get(X)-c1.Get(X)*c2.Get(Z))
	c3.Set(Z, c1.Get(X)*c2.Get(Y)-c1.Get(Y)*c2.Get(X))
	return c3
}
