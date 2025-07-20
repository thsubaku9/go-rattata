package coordinates

//go:generate stringer -type=CoordinateAxis
type CoordinateAxis int

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
