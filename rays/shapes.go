package rays

import "rattata/coordinates"

type Shape interface {
	Name() string
}

type Sphere struct {
	Origin coordinates.Coordinate
	Radius float32
}

func (s Sphere) Name() string {
	return "Sphere"
}

func NewSphere(origin coordinates.Coordinate, radius float32) Sphere {
	if origin.IsAPoint() {
		return Sphere{Origin: origin, Radius: radius}
	}

	panic("origin is not a point")
}
