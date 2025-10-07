package rays

import (
	"math"
	"rattata/coordinates"
	"rattata/matrices"

	"github.com/gofrs/uuid"
)

// ---------------------------------- Shapes ----------------------------------
type Shape interface {
	Name() string
	Transformation() matrices.Matrix
	IntersectWithRay(ray Ray) []Intersection
	/*
		Returns the normalized vector perpendicular to the shape at the given world point
	*/
	NormalAtPoint(world_point coordinates.Coordinate) coordinates.Coordinate
	GetMaterial() Material
	Id() string
}

// ---------------------------------- Sphere ----------------------------------
type Sphere struct {
	Origin            coordinates.Coordinate
	Radius            float64
	transformationMat matrices.Matrix
	Material          Material
	id                string
}

func NewSphere(origin coordinates.Coordinate, radius float64) Sphere {
	if !origin.IsAPoint() {
		panic("origin is not a point")
	}

	new_uuid, _ := uuid.NewV4()
	return Sphere{Origin: origin, Radius: radius, transformationMat: matrices.NewIdentityMatrix(4), Material: CreateDefaultMaterial(), id: new_uuid.String()}
}

func NewCenteredSphere() Sphere {
	return NewSphere(coordinates.CreatePoint(0, 0, 0), 1.0)
}

func NewGlassSphere() Sphere {
	s := NewCenteredSphere()
	s.Material.Transparency = 1.0
	s.Material.RefractiveIndex = 1.5
	return s
}

func (s Sphere) Name() string {
	return "Sphere"
}

func (s Sphere) Id() string {
	return s.id
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

func (s Sphere) IntersectWithRay(ray Ray) []Intersection {
	sphere_to_ray := ray.Origin.Sub(&s.Origin)

	a := ray.Direction.DotP(&ray.Direction)
	b := ray.Direction.DotP(sphere_to_ray) * 2
	c := sphere_to_ray.DotP(sphere_to_ray) - s.Radius
	discriminant := b*b - 4*a*c

	if discriminant < 0 {
		return []Intersection{}
	}

	t1 := (-b - math.Sqrt(discriminant)) / (2 * a)
	t2 := (-b + math.Sqrt(discriminant)) / (2 * a)
	return Intersections(NewIntersection(t1, s), NewIntersection(t2, s))
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

// ---------------------------------- XZPlane ----------------------------------
type XZPlane struct {
	Origin            coordinates.Coordinate
	transformationMat matrices.Matrix
	Material          Material
	id                string
}

func NewPlane(origin coordinates.Coordinate) XZPlane {
	if !origin.IsAPoint() {
		panic("origin is not a point")
	}

	new_uuid, _ := uuid.NewV4()

	return XZPlane{Origin: origin, transformationMat: matrices.NewIdentityMatrix(4), Material: CreateDefaultMaterial(), id: new_uuid.String()}
}

func (p XZPlane) Name() string {
	return "Plane"
}

func (p XZPlane) Id() string {
	return p.id
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

func (p XZPlane) IntersectWithRay(ray Ray) []Intersection {
	if math.Abs(ray.Direction[coordinates.Y]) < EPSILON {
		return []Intersection{}
	}

	t := -ray.Origin[coordinates.Y] / ray.Direction[coordinates.Y]
	return Intersections(NewIntersection(t, p))
}

func (p XZPlane) NormalAtPoint(world_point coordinates.Coordinate) coordinates.Coordinate {
	inverse_transformation, _ := p.Transformation().Inverse()
	obj_normal := coordinates.CreateVector(0, 1, 0)

	world_normal := matrices.PerformOrderedChainingOps(matrices.CoordinateToMatrix(obj_normal), inverse_transformation.T())
	world_normal.Set(3, 0, 0)

	res := matrices.MatrixToCoordinate(world_normal)
	return *res.Norm()

}

// ---------------------------------- Cube ----------------------------------

type Cube struct {
	transformationMat matrices.Matrix
	Material          Material
	id                string
}

func NewCube() Cube {
	new_uuid, _ := uuid.NewV4()
	return Cube{transformationMat: matrices.NewIdentityMatrix(4), Material: CreateDefaultMaterial(), id: new_uuid.String()}
}

func (c Cube) Id() string {
	return c.id
}

func (c Cube) Transformation() matrices.Matrix {
	return c.transformationMat
}

func (c *Cube) SetTransformation(mt matrices.Matrix) {
	c.transformationMat = mt
}

func (c Cube) GetMaterial() Material {
	return c.Material
}

func (c Cube) IntersectWithRay(ray Ray) []Intersection {
	return []Intersection{}
}

func (c Cube) NormalAtPoint(world_point coordinates.Coordinate) coordinates.Coordinate {
	return coordinates.CreateVector(0, 0, 0)
}

// ---------------------------------- Cylinder ----------------------------------

type Cylinder struct {
	transformationMat matrices.Matrix
	Material          Material
	id                string
}

func (cy Cylinder) Id() string {
	return cy.id
}

func (cy Cylinder) Transformation() matrices.Matrix {
	return cy.transformationMat
}

func (cy *Cylinder) SetTransformation(mt matrices.Matrix) {
	cy.transformationMat = mt
}

func (cy Cylinder) GetMaterial() Material {
	return cy.Material
}

func (cy Cylinder) IntersectWithRay(ray Ray) []Intersection {
	return []Intersection{}
}

func (cy Cylinder) NormalAtPoint(world_point coordinates.Coordinate) coordinates.Coordinate {
	return coordinates.CreateVector(0, 0, 0)
}

// ---------------------------------- Cone ----------------------------------

type Cone struct {
	transformationMat matrices.Matrix
	Material          Material
	id                string
}

func (co Cone) Id() string {
	return co.id
}

func (co Cone) Transformation() matrices.Matrix {
	return co.transformationMat
}

func (co *Cone) SetTransformation(mt matrices.Matrix) {
	co.transformationMat = mt
}

func (co Cone) GetMaterial() Material {
	return co.Material
}

func (co Cone) IntersectWithRay(ray Ray) []Intersection {
	return []Intersection{}
}

func (co Cone) NormalAtPoint(world_point coordinates.Coordinate) coordinates.Coordinate {
	return coordinates.CreateVector(0, 0, 0)
}
