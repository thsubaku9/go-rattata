package rays

import (
	"fmt"
	"math"
	"rattata/coordinates"
	"rattata/matrices"
	"sort"

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
	Parent() *Group
}

type IsGroupable interface {
	SetParent(parent *Group)
	GetRefAddress() *Shape
}

// ---------------------------------- Sphere ----------------------------------
type Sphere struct {
	Origin            coordinates.Coordinate
	Radius            float64
	transformationMat matrices.Matrix
	Material          Material
	id                string
	parent            *Group
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

func (s Sphere) Parent() *Group {
	return s.parent
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
	obj_point := world_to_object_orientation(s, world_point)
	obj_normal := *obj_point.Sub(&s.Origin)
	return normal_to_world_orientation(s, obj_normal)
}

func (s *Sphere) SetParent(parent *Group) {
	s.parent = parent
}

func (s *Sphere) GetRefAddress() *Shape {
	var _shape Shape = s
	return &_shape
}

// ---------------------------------- XZPlane ----------------------------------
type XZPlane struct {
	Origin            coordinates.Coordinate
	transformationMat matrices.Matrix
	Material          Material
	id                string
	parent            *Group
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

func (p XZPlane) Parent() *Group {
	return p.parent
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
	obj_normal := coordinates.CreateVector(0, 1, 0)
	return normal_to_world_orientation(p, obj_normal)
}

func (p *XZPlane) SetParent(parent *Group) {
	p.parent = parent
}

func (p *XZPlane) GetRefAddress() *Shape {
	var _shape Shape = p
	return &_shape
}

// ---------------------------------- Cube ----------------------------------

type Cube struct {
	transformationMat matrices.Matrix
	Material          Material
	id                string
	parent            *Group
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

func (c Cube) Parent() *Group {
	return c.parent
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
	obj_point := world_to_object_orientation(c, world_point)

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

	return normal_to_world_orientation(c, normal_v)
}

func (c *Cube) SetParent(parent *Group) {
	c.parent = parent
}

func (c *Cube) GetRefAddress() *Shape {
	var _shape Shape = c
	return &_shape
}

// ---------------------------------- XZCylinder ----------------------------------

type XZCylinder struct {
	transformationMat matrices.Matrix
	Material          Material
	Minimum           float64
	Maximum           float64
	Closed            bool
	id                string
	parent            *Group
}

func NewXZCylinder() XZCylinder {
	new_uuid, _ := uuid.NewV4()
	return XZCylinder{transformationMat: matrices.NewIdentityMatrix(4), Material: CreateDefaultMaterial(), id: new_uuid.String(),
		Minimum: math.Inf(-1), Maximum: math.Inf(1), Closed: false}
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

func (cy XZCylinder) Parent() *Group {
	return cy.parent
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

	if ray_wrt_obj.Origin.Get(coordinates.Y)+t1*ray_wrt_obj.Direction.Get(coordinates.Y) > cy.Minimum &&
		ray_wrt_obj.Origin.Get(coordinates.Y)+t1*ray_wrt_obj.Direction.Get(coordinates.Y) < cy.Maximum {
		res = append(res, NewIntersection(t1, cy))
	}

	if ray_wrt_obj.Origin.Get(coordinates.Y)+t2*ray_wrt_obj.Direction.Get(coordinates.Y) > cy.Minimum &&
		ray_wrt_obj.Origin.Get(coordinates.Y)+t2*ray_wrt_obj.Direction.Get(coordinates.Y) < cy.Maximum {
		res = append(res, NewIntersection(t2, cy))
	}

	return res
}

func (cy XZCylinder) checkCapIntersection(ray_wrt_obj Ray) []Intersection {
	if !cy.Closed || math.Abs(ray_wrt_obj.Direction.Get(coordinates.Y)) < EPSILON {
		return []Intersection{}
	}

	res := make([]Intersection, 0)

	t3 := (cy.Minimum - ray_wrt_obj.Origin.Get(coordinates.Y)) / ray_wrt_obj.Direction.Get(coordinates.Y)
	t4 := (cy.Maximum - ray_wrt_obj.Origin.Get(coordinates.Y)) / ray_wrt_obj.Direction.Get(coordinates.Y)

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
	obj_point := world_to_object_orientation(cy, world_point)

	dist := obj_point.Get(coordinates.X)*obj_point.Get(coordinates.X) + obj_point.Get(coordinates.Z)*obj_point.Get(coordinates.Z)

	var normal_v coordinates.Coordinate
	if dist < 1 && obj_point.Get(coordinates.Y) >= cy.Maximum-EPSILON && cy.Closed {
		normal_v = coordinates.CreateVector(0, 1, 0)
	} else if dist < 1 && obj_point.Get(coordinates.Y) <= cy.Minimum+EPSILON && cy.Closed {
		normal_v = coordinates.CreateVector(0, -1, 0)
	} else {
		normal_v = coordinates.CreateVector(obj_point.Get(coordinates.X), 0, obj_point.Get(coordinates.Z))
	}

	return normal_to_world_orientation(cy, normal_v)

}

func (cy *XZCylinder) SetParent(parent *Group) {
	cy.parent = parent
}

func (cy *XZCylinder) GetRefAddress() *Shape {
	var _shape Shape = cy
	return &_shape
}

// ---------------------------------- Cone ----------------------------------

type Cone struct {
	transformationMat matrices.Matrix
	Material          Material
	Minimum           float64
	Maximum           float64
	Closed            bool
	id                string
	parent            *Group
}

func NewDoubleNappedCone() Cone {
	new_uuid, _ := uuid.NewV4()
	return Cone{transformationMat: matrices.NewIdentityMatrix(4), Material: CreateDefaultMaterial(), id: new_uuid.String(),
		Minimum: math.Inf(-1), Maximum: math.Inf(1), Closed: false}
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

func (co Cone) Parent() *Group {
	return co.parent
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

	if ray_wrt_obj.Origin.Get(coordinates.Y)+t1*ray_wrt_obj.Direction.Get(coordinates.Y) > co.Minimum &&
		ray_wrt_obj.Origin.Get(coordinates.Y)+t1*ray_wrt_obj.Direction.Get(coordinates.Y) < co.Maximum {
		res = append(res, NewIntersection(t1, co))
	}

	if ray_wrt_obj.Origin.Get(coordinates.Y)+t2*ray_wrt_obj.Direction.Get(coordinates.Y) > co.Minimum &&
		ray_wrt_obj.Origin.Get(coordinates.Y)+t2*ray_wrt_obj.Direction.Get(coordinates.Y) < co.Maximum {
		res = append(res, NewIntersection(t2, co))
	}

	return res
}

func (co Cone) checkCapIntersection(ray_wrt_obj Ray) []Intersection {
	if !co.Closed || math.Abs(ray_wrt_obj.Direction.Get(coordinates.Y)) < EPSILON {
		return []Intersection{}
	}

	res := make([]Intersection, 0)

	t3 := (co.Minimum - ray_wrt_obj.Origin.Get(coordinates.Y)) / ray_wrt_obj.Direction.Get(coordinates.Y)
	t4 := (co.Maximum - ray_wrt_obj.Origin.Get(coordinates.Y)) / ray_wrt_obj.Direction.Get(coordinates.Y)
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
	obj_point := world_to_object_orientation(co, world_point)

	dist := obj_point.Get(coordinates.X)*obj_point.Get(coordinates.X) + obj_point.Get(coordinates.Z)*obj_point.Get(coordinates.Z)
	radius_sqr := math.Pow(obj_point.Get(coordinates.Y), 2)

	var normal_v coordinates.Coordinate

	if dist < radius_sqr && obj_point.Get(coordinates.Y) >= co.Maximum-EPSILON && co.Closed {
		normal_v = coordinates.CreateVector(0, 1, 0)
	} else if dist < radius_sqr && obj_point.Get(coordinates.Y) <= co.Minimum+EPSILON && co.Closed {
		normal_v = coordinates.CreateVector(0, -1, 0)
	} else {
		y := math.Sqrt(obj_point.Get(coordinates.X)*obj_point.Get(coordinates.X) + obj_point.Get(coordinates.Z)*obj_point.Get(coordinates.Z))
		if obj_point.Get(coordinates.Y) > 0 {
			y = -y
		}

		normal_v = coordinates.CreateVector(obj_point.Get(coordinates.X), y, obj_point.Get(coordinates.Z))
	}

	return normal_to_world_orientation(co, normal_v)
}

func (co *Cone) SetParent(parent *Group) {
	co.parent = parent
}

func (co *Cone) GetRefAddress() *Shape {
	var _shape Shape = co
	return &_shape
}

// ---------------------------------- Group ----------------------------------

type Group struct {
	containedShapes   []*Shape
	transformationMat matrices.Matrix
	id                string
	parent            *Group
}

func NewGroup() Group {
	new_uuid, _ := uuid.NewV4()
	return Group{containedShapes: make([]*Shape, 0), transformationMat: matrices.NewIdentityMatrix(4), id: new_uuid.String()}
}

func (g Group) Id() string {
	return g.id
}

func (g Group) Name() string {
	return "Group"
}

func (g Group) Transformation() matrices.Matrix {
	return g.transformationMat
}

func (g *Group) SetTransformation(mt matrices.Matrix) {
	g.transformationMat = mt
}

func (g Group) Parent() *Group {
	return g.parent
}

func (g Group) GetMaterial() Material {
	panic("Groups do not have materials")
}

func (g Group) IntersectWithRay(ray_wrt_obj Ray) []Intersection {
	if len(g.containedShapes) == 0 {
		return []Intersection{}
	}

	all_intersections := make([]Intersection, 0)
	for _, shape_ptr := range g.containedShapes {
		all_intersections = append(all_intersections, Intersect(*shape_ptr, ray_wrt_obj)...)
	}

	sort.Slice(all_intersections, func(i, j int) bool {
		return all_intersections[i].Tvalue <= all_intersections[j].Tvalue
	})

	return all_intersections
}

func (g Group) NormalAtPoint(world_point coordinates.Coordinate) coordinates.Coordinate {
	panic("Groups do not have normals")
}

func (g *Group) IndoctrinateShapeToGroup(shape_to_indoctrinate IsGroupable) {
	g.containedShapes = append(g.containedShapes, shape_to_indoctrinate.GetRefAddress())
	shape_to_indoctrinate.SetParent(g)
}

func (g *Group) SetParent(p *Group) {
	g.parent = p
}

func (g *Group) GetRefAddress() *Shape {
	var _shape Shape = g
	return &_shape
}

func world_to_object_orientation(shape Shape, world_coord coordinates.Coordinate) coordinates.Coordinate {
	cur_coord := world_coord
	if shape.Parent() != nil {
		cur_coord = world_to_object_orientation(*shape.Parent(), cur_coord)
	}

	inv_transform, _ := shape.Transformation().Inverse()
	cur_coord_mat := matrices.PerformOrderedChainingOps(matrices.CoordinateToMatrix(cur_coord), inv_transform)
	return matrices.MatrixToCoordinate(cur_coord_mat)
}

func object_to_world_orientation(shape Shape, object_coord coordinates.Coordinate) coordinates.Coordinate {

	cur_coord_mat := matrices.PerformOrderedChainingOps(matrices.CoordinateToMatrix(object_coord), shape.Transformation())
	cur_coord := matrices.MatrixToCoordinate(cur_coord_mat)

	if shape.Parent() != nil {
		cur_coord = object_to_world_orientation(*shape.Parent(), cur_coord)
	}

	return cur_coord

}

func normal_to_world_orientation(shape Shape, object_normal_v coordinates.Coordinate) coordinates.Coordinate {

	inverse_transformation, _ := shape.Transformation().Inverse()
	world_normal := matrices.PerformOrderedChainingOps(matrices.CoordinateToMatrix(object_normal_v), inverse_transformation.T())
	res := matrices.MatrixToCoordinate(world_normal)
	res.Set(coordinates.W, 0)
	res = *res.Norm()

	if shape.Parent() != nil {
		res = normal_to_world_orientation(*shape.Parent(), res)
	}

	return res
}
