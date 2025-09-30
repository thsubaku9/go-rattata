package rays

import (
	"rattata/coordinates"
	"rattata/matrices"
)

type Shape interface {
	Name() string
	Transformation() matrices.Matrix
	NormalAtPoint(world_point coordinates.Coordinate) coordinates.Coordinate
	GetMaterial() Material
}

type Sphere struct {
	Origin            coordinates.Coordinate
	Radius            float64
	transformationMat matrices.Matrix
	Material          Material
}

func NewSphere(origin coordinates.Coordinate, radius float64) Sphere {
	if !origin.IsAPoint() {
		panic("origin is not a point")
	}

	return Sphere{Origin: origin, Radius: radius, transformationMat: matrices.NewIdentityMatrix(4), Material: CreateDefaultMaterial()}
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

func (s Sphere) GetMaterial() Material {
	return s.Material
}

func (s Sphere) NormalAtPoint(world_point coordinates.Coordinate) coordinates.Coordinate {
	inverse_transformation, _ := s.Transformation().Inverse()

	obj_point_mat := matrices.PerformOrderedChainingOps(matrices.CoordinateToMatrix(world_point), inverse_transformation)

	obj_point := matrices.MatrixToCoordinate(obj_point_mat)
	obj_normal := *obj_point.Sub(&s.Origin)

	world_normal := matrices.PerformOrderedChainingOps(matrices.CoordinateToMatrix(obj_normal), inverse_transformation.T())
	world_normal.Set(3, 0, 0)

	res := matrices.MatrixToCoordinate(world_normal)
	return *res.Norm()
}

type XZPlane struct {
	Origin            coordinates.Coordinate
	transformationMat matrices.Matrix
	Material          Material
}

func NewPlane(origin coordinates.Coordinate) XZPlane {
	if !origin.IsAPoint() {
		panic("origin is not a point")
	}

	return XZPlane{Origin: origin, transformationMat: matrices.NewIdentityMatrix(4), Material: CreateDefaultMaterial()}
}

func (p XZPlane) Name() string {
	return "Plane"
}

func (p XZPlane) Transformation() matrices.Matrix {
	return p.transformationMat
}

func (p *XZPlane) SetTransformation(mt matrices.Matrix) {
	p.transformationMat = mt
}

func (p XZPlane) GetMaterial() Material {
	return p.Material
}

func (p XZPlane) NormalAtPoint(world_point coordinates.Coordinate) coordinates.Coordinate {
	inverse_transformation, _ := p.Transformation().Inverse()
	obj_normal := coordinates.CreateVector(0, 1, 0)

	world_normal := matrices.PerformOrderedChainingOps(matrices.CoordinateToMatrix(obj_normal), inverse_transformation.T())
	world_normal.Set(3, 0, 0)

	res := matrices.MatrixToCoordinate(world_normal)
	return *res.Norm()

}
