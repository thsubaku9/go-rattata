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

func (s Sphere) NormalAtPoint(world_point coordinates.Coordinate) coordinates.Coordinate {
	inverse_transformation, _ := s.Transformation().Adj()

	obj_point_mat := matrices.PerformOrderedChainingOps(matrices.CoordinateToMatrix(world_point), inverse_transformation)

	obj_point := matrices.MatrixToCoordinate(obj_point_mat)
	obj_normal := obj_point.Sub(&s.Origin)

	world_normal := matrices.PerformOrderedChainingOps(matrices.CoordinateToMatrix(*obj_normal), inverse_transformation.T())
	world_normal.Set(3, 0, 0)

	res := matrices.MatrixToCoordinate(world_normal)
	return *res.Norm()
}

func NewSphere(origin coordinates.Coordinate, radius float32) Sphere {
	if origin.IsAPoint() {
		return Sphere{Origin: origin, Radius: radius, transformationMat: matrices.NewIdentityMatrix(4)}
	}

	panic("origin is not a point")
}
