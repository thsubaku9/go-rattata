package rays

import (
	"fmt"
	"math"
	"rattata/coordinates"
	"rattata/matrices"

	"github.com/gofrs/uuid"
)

// ---------------------------------- Shapes ----------------------------------
type Shape interface {
	Name() string
	Transformation() matrices.Matrix
	IntersectWithRay(ray_wrt_obj Ray) []Intersection
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

func (s Sphere) IntersectWithRay(ray_wrt_obj Ray) []Intersection {
	sphere_to_ray := ray_wrt_obj.Origin.Sub(&s.Origin)

	a := ray_wrt_obj.Direction.DotP(&ray_wrt_obj.Direction)
	b := ray_wrt_obj.Direction.DotP(sphere_to_ray) * 2
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

func (p XZPlane) IntersectWithRay(ray_wrt_obj Ray) []Intersection {
	if math.Abs(ray_wrt_obj.Direction.Get(coordinates.Y)) < EPSILON {
		return []Intersection{}
	}

	t := -ray_wrt_obj.Origin.Get(coordinates.Y) / ray_wrt_obj.Direction.Get(coordinates.Y)
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

func (c Cube) Name() string {
	return "Cube"
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

func (c Cube) IntersectWithRay(ray_wrt_obj Ray) []Intersection {

	y_tmin, y_tmax := c.axis_intersection_points(ray_wrt_obj.Origin.Get(coordinates.Y), ray_wrt_obj.Direction.Get(coordinates.Y))
	x_tmin, x_tmax := c.axis_intersection_points(ray_wrt_obj.Origin.Get(coordinates.X), ray_wrt_obj.Direction.Get(coordinates.X))
	z_tmin, z_tmax := c.axis_intersection_points(ray_wrt_obj.Origin.Get(coordinates.Z), ray_wrt_obj.Direction.Get(coordinates.Z))

	tmin := max(x_tmin, y_tmin, z_tmin)
	tmax := min(x_tmax, y_tmax, z_tmax)

	if tmin > tmax {
		return []Intersection{}
	}
	return []Intersection{{Tvalue: tmin, Obj: c}, {Tvalue: tmax, Obj: c}}
}

func (c Cube) axis_intersection_points(origin, direction float64) (float64, float64) {
	tmin_numerator := -1 - origin
	tmax_numerator := 1 - origin

	var tmin, tmax float64
	if math.Abs(direction) >= EPSILON {
		tmin = tmin_numerator / direction
		tmax = tmax_numerator / direction
	} else {
		tmin = tmin_numerator * math.Inf(1)
		tmax = tmax_numerator * math.Inf(1)
	}

	if tmin > tmax {
		return tmax, tmin
	}
	return tmin, tmax
}

func (c Cube) NormalAtPoint(world_point coordinates.Coordinate) coordinates.Coordinate {
	inverse_transformation, _ := c.Transformation().Inverse()

	obj_point_mat := matrices.PerformOrderedChainingOps(matrices.CoordinateToMatrix(world_point), inverse_transformation)
	obj_point := matrices.MatrixToCoordinate(obj_point_mat)

	maxc := math.Max(math.Abs(obj_point.Get(coordinates.X)), math.Max(math.Abs(obj_point.Get(coordinates.Y)), math.Abs(obj_point.Get(coordinates.Z))))

	var normal_v coordinates.Coordinate

	switch {
	case maxc == math.Abs(obj_point.Get(coordinates.X)):
		normal_v = coordinates.CreateVector(obj_point.Get(coordinates.X), 0, 0)
	case maxc == math.Abs(obj_point.Get(coordinates.Y)):
		normal_v = coordinates.CreateVector(0, obj_point.Get(coordinates.Y), 0)
	default:
		normal_v = coordinates.CreateVector(0, 0, obj_point.Get(coordinates.Z))
	}

	world_normal := matrices.PerformOrderedChainingOps(matrices.CoordinateToMatrix(normal_v), inverse_transformation.T())

	res := matrices.MatrixToCoordinate(world_normal)
	return *res.Norm()
}

// ---------------------------------- XZCylinder ----------------------------------

type XZCylinder struct {
	transformationMat matrices.Matrix
	Material          Material
	minimum           float64
	maximum           float64
	closed            bool
	id                string
}

func NewXZCylinder() XZCylinder {
	new_uuid, _ := uuid.NewV4()
	return XZCylinder{transformationMat: matrices.NewIdentityMatrix(4), Material: CreateDefaultMaterial(), id: new_uuid.String(),
		minimum: math.Inf(-1), maximum: math.Inf(1), closed: false}
}

func (cy XZCylinder) Id() string {
	return cy.id
}

func (cy XZCylinder) Name() string {
	return "XZCylinder"
}

func (cy XZCylinder) Transformation() matrices.Matrix {
	return cy.transformationMat
}

func (cy *XZCylinder) SetTransformation(mt matrices.Matrix) {
	cy.transformationMat = mt
}

func (cy XZCylinder) GetMaterial() Material {
	return cy.Material
}

func (cy XZCylinder) IntersectWithRay(ray_wrt_obj Ray) []Intersection {
	res := make([]Intersection, 0)

	res = append(res, cy.checkCircularIntersection(ray_wrt_obj)...)
	res = append(res, cy.checkCapIntersection(ray_wrt_obj)...)
	return res
}

func (cy XZCylinder) checkCircularIntersection(ray_wrt_obj Ray) []Intersection {
	a := ray_wrt_obj.Direction.Get(coordinates.X)*ray_wrt_obj.Direction.Get(coordinates.X) + ray_wrt_obj.Direction.Get(coordinates.Z)*ray_wrt_obj.Direction.Get(coordinates.Z)
	if math.Abs(a) < EPSILON {
		return []Intersection{}
	}

	b := 2*ray_wrt_obj.Origin.Get(coordinates.X)*ray_wrt_obj.Direction.Get(coordinates.X) + 2*ray_wrt_obj.Origin.Get(coordinates.Z)*ray_wrt_obj.Direction.Get(coordinates.Z)
	c := ray_wrt_obj.Origin.Get(coordinates.X)*ray_wrt_obj.Origin.Get(coordinates.X) + ray_wrt_obj.Origin.Get(coordinates.Z)*ray_wrt_obj.Origin.Get(coordinates.Z) - 1

	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return []Intersection{}
	}

	t1 := (-b - math.Sqrt(discriminant)) / (2 * a)
	t2 := (-b + math.Sqrt(discriminant)) / (2 * a)

	res := make([]Intersection, 0)

	if ray_wrt_obj.Origin.Get(coordinates.Y)+t1*ray_wrt_obj.Direction.Get(coordinates.Y) > cy.minimum &&
		ray_wrt_obj.Origin.Get(coordinates.Y)+t1*ray_wrt_obj.Direction.Get(coordinates.Y) < cy.maximum {
		res = append(res, NewIntersection(t1, cy))
	}

	if ray_wrt_obj.Origin.Get(coordinates.Y)+t2*ray_wrt_obj.Direction.Get(coordinates.Y) > cy.minimum &&
		ray_wrt_obj.Origin.Get(coordinates.Y)+t2*ray_wrt_obj.Direction.Get(coordinates.Y) < cy.maximum {
		res = append(res, NewIntersection(t2, cy))
	}

	return res
}

func (cy XZCylinder) checkCapIntersection(ray_wrt_obj Ray) []Intersection {
	if !cy.closed || math.Abs(ray_wrt_obj.Direction.Get(coordinates.Y)) < EPSILON {
		return []Intersection{}
	}

	res := make([]Intersection, 0)

	t3 := (cy.minimum - ray_wrt_obj.Origin.Get(coordinates.Y)) / ray_wrt_obj.Direction.Get(coordinates.Y)
	t4 := (cy.maximum - ray_wrt_obj.Origin.Get(coordinates.Y)) / ray_wrt_obj.Direction.Get(coordinates.Y)

	if cy.withinBoundingRadius(ray_wrt_obj, t3) {
		res = append(res, NewIntersection(t3, cy))
	}
	if cy.withinBoundingRadius(ray_wrt_obj, t4) {
		res = append(res, NewIntersection(t4, cy))
	}

	return res

}

func (cy XZCylinder) withinBoundingRadius(ray_wrt_obj Ray, t float64) bool {
	x := ray_wrt_obj.Origin.Get(coordinates.X) + t*ray_wrt_obj.Direction.Get(coordinates.X)
	z := ray_wrt_obj.Origin.Get(coordinates.Z) + t*ray_wrt_obj.Direction.Get(coordinates.Z)

	return (x*x + z*z) <= 1
}

func (cy XZCylinder) NormalAtPoint(world_point coordinates.Coordinate) coordinates.Coordinate {
	inverse_transformation, _ := cy.Transformation().Inverse()
	obj_point_matrix := matrices.PerformOrderedChainingOps(matrices.CoordinateToMatrix(world_point), inverse_transformation)
	obj_point := matrices.MatrixToCoordinate(obj_point_matrix)

	dist := obj_point.Get(coordinates.X)*obj_point.Get(coordinates.X) + obj_point.Get(coordinates.Z)*obj_point.Get(coordinates.Z)

	var normal_v coordinates.Coordinate
	if dist < 1 && obj_point.Get(coordinates.Y) >= cy.maximum-EPSILON && cy.closed {
		normal_v = coordinates.CreateVector(0, 1, 0)
	} else if dist < 1 && obj_point.Get(coordinates.Y) <= cy.minimum+EPSILON && cy.closed {
		normal_v = coordinates.CreateVector(0, -1, 0)
	} else {
		normal_v = coordinates.CreateVector(obj_point.Get(coordinates.X), 0, obj_point.Get(coordinates.Z))
	}

	world_normal := matrices.PerformOrderedChainingOps(matrices.CoordinateToMatrix(normal_v), inverse_transformation.T())
	res := matrices.MatrixToCoordinate(world_normal)
	return *res.Norm()
}

// ---------------------------------- Cone ----------------------------------

type Cone struct {
	transformationMat matrices.Matrix
	Material          Material
	minimum           float64
	maximum           float64
	closed            bool
	id                string
}

func NewDoubleNappedCone() Cone {
	new_uuid, _ := uuid.NewV4()
	return Cone{transformationMat: matrices.NewIdentityMatrix(4), Material: CreateDefaultMaterial(), id: new_uuid.String(),
		minimum: math.Inf(-1), maximum: math.Inf(1), closed: false}
}

func (co Cone) Id() string {
	return co.id
}

func (co Cone) Name() string {
	return "Cone"
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

func (co Cone) IntersectWithRay(ray_wrt_obj Ray) []Intersection {
	res := make([]Intersection, 0)

	res = append(res, co.checkCircularIntersection(ray_wrt_obj)...)
	res = append(res, co.checkCapIntersection(ray_wrt_obj)...)
	return res
}

func (co Cone) checkCircularIntersection(ray_wrt_obj Ray) []Intersection {
	a := ray_wrt_obj.Direction.Get(coordinates.X)*ray_wrt_obj.Direction.Get(coordinates.X) -
		ray_wrt_obj.Direction.Get(coordinates.Y)*ray_wrt_obj.Direction.Get(coordinates.Y) +
		ray_wrt_obj.Direction.Get(coordinates.Z)*ray_wrt_obj.Direction.Get(coordinates.Z)

	b := 2*ray_wrt_obj.Origin.Get(coordinates.X)*ray_wrt_obj.Direction.Get(coordinates.X) -
		2*ray_wrt_obj.Origin.Get(coordinates.Y)*ray_wrt_obj.Direction.Get(coordinates.Y) +
		2*ray_wrt_obj.Origin.Get(coordinates.Z)*ray_wrt_obj.Direction.Get(coordinates.Z)

	c := ray_wrt_obj.Origin.Get(coordinates.X)*ray_wrt_obj.Origin.Get(coordinates.X) -
		ray_wrt_obj.Origin.Get(coordinates.Y)*ray_wrt_obj.Origin.Get(coordinates.Y) +
		ray_wrt_obj.Origin.Get(coordinates.Z)*ray_wrt_obj.Origin.Get(coordinates.Z)

	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return []Intersection{}
	}

	if math.Abs(a) < EPSILON {
		if math.Abs(b) < EPSILON {
			return []Intersection{}
		}
		t := -c / (2 * b)
		return []Intersection{NewIntersection(t, co)}
	}

	t1 := (-b - math.Sqrt(discriminant)) / (2 * a)
	t2 := (-b + math.Sqrt(discriminant)) / (2 * a)

	fmt.Println(discriminant)
	fmt.Println(t1, t2)

	res := make([]Intersection, 0)

	if ray_wrt_obj.Origin.Get(coordinates.Y)+t1*ray_wrt_obj.Direction.Get(coordinates.Y) > co.minimum &&
		ray_wrt_obj.Origin.Get(coordinates.Y)+t1*ray_wrt_obj.Direction.Get(coordinates.Y) < co.maximum {
		res = append(res, NewIntersection(t1, co))
	}

	if ray_wrt_obj.Origin.Get(coordinates.Y)+t2*ray_wrt_obj.Direction.Get(coordinates.Y) > co.minimum &&
		ray_wrt_obj.Origin.Get(coordinates.Y)+t2*ray_wrt_obj.Direction.Get(coordinates.Y) < co.maximum {
		res = append(res, NewIntersection(t2, co))
	}

	return res
}

func (co Cone) checkCapIntersection(ray_wrt_obj Ray) []Intersection {
	if !co.closed || math.Abs(ray_wrt_obj.Direction.Get(coordinates.Y)) < EPSILON {
		return []Intersection{}
	}

	res := make([]Intersection, 0)

	t3 := (co.minimum - ray_wrt_obj.Origin.Get(coordinates.Y)) / ray_wrt_obj.Direction.Get(coordinates.Y)
	t4 := (co.maximum - ray_wrt_obj.Origin.Get(coordinates.Y)) / ray_wrt_obj.Direction.Get(coordinates.Y)
	if co.withinBoundingRadius(ray_wrt_obj, t3) {
		res = append(res, NewIntersection(t3, co))
	}
	if co.withinBoundingRadius(ray_wrt_obj, t4) {
		res = append(res, NewIntersection(t4, co))
	}

	return res
}

func (co Cone) withinBoundingRadius(ray_wrt_obj Ray, t float64) bool {
	x := ray_wrt_obj.Origin.Get(coordinates.X) + t*ray_wrt_obj.Direction.Get(coordinates.X)
	y := ray_wrt_obj.Origin.Get(coordinates.Y) + t*ray_wrt_obj.Direction.Get(coordinates.Y)
	z := ray_wrt_obj.Origin.Get(coordinates.Z) + t*ray_wrt_obj.Direction.Get(coordinates.Z)

	return (x*x + z*z) <= y*y
}

func (co Cone) NormalAtPoint(world_point coordinates.Coordinate) coordinates.Coordinate {
	inverse_transformation, _ := co.Transformation().Inverse()

	obj_point_matrix := matrices.PerformOrderedChainingOps(matrices.CoordinateToMatrix(world_point), inverse_transformation)
	obj_point := matrices.MatrixToCoordinate(obj_point_matrix)

	dist := obj_point.Get(coordinates.X)*obj_point.Get(coordinates.X) + obj_point.Get(coordinates.Z)*obj_point.Get(coordinates.Z)
	radius_sqr := math.Pow(obj_point.Get(coordinates.Y), 2)

	var normal_v coordinates.Coordinate

	if dist < radius_sqr && obj_point.Get(coordinates.Y) >= co.maximum-EPSILON && co.closed {
		normal_v = coordinates.CreateVector(0, 1, 0)
	} else if dist < radius_sqr && obj_point.Get(coordinates.Y) <= co.minimum+EPSILON && co.closed {
		normal_v = coordinates.CreateVector(0, -1, 0)
	} else {
		y := math.Sqrt(obj_point.Get(coordinates.X)*obj_point.Get(coordinates.X) + obj_point.Get(coordinates.Z)*obj_point.Get(coordinates.Z))
		if obj_point.Get(coordinates.Y) > 0 {
			y = -y
		}

		normal_v = coordinates.CreateVector(obj_point.Get(coordinates.X), y, obj_point.Get(coordinates.Z))
		fmt.Println("FREEDOM ", normal_v)
	}

	world_normal := matrices.PerformOrderedChainingOps(matrices.CoordinateToMatrix(normal_v), inverse_transformation.T())
	res := matrices.MatrixToCoordinate(world_normal)
	return res
}
