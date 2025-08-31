package coordinates

type Shape interface {
	Name() string
}

type Sphere struct {
	Origin Coordinate
	Radius float32
}

func (s Sphere) Name() string {
	return "Sphere"
}

func NewSphere(origin Coordinate, radius float32) Sphere {
	if origin.IsAPoint() {
		return Sphere{Origin: origin, Radius: radius}
	}

	panic("origin is not a point")
}
