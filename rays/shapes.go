package rays

import (
	"rattata/coordinates"
	"rattata/matrices"
)

type Shape interface {
	Name() string
	Transformation() matrices.Matrix
}

type Sphere struct {
	Origin            coordinates.Coordinate
	Radius            float32
	transformationMat matrices.Matrix
}

func (s Sphere) Name() string {
	return "Sphere"
}

func (s Sphere) Transformation() matrices.Matrix {
	return s.transformationMat
}

func (s *Sphere) SetTransformation(mt matrices.Matrix) {
	s.transformationMat = mt
}

func NewSphere(origin coordinates.Coordinate, radius float32) Sphere {
	if origin.IsAPoint() {
		return Sphere{Origin: origin, Radius: radius, transformationMat: matrices.NewIdentityMatrix(4)}
	}

	panic("origin is not a point")
}
