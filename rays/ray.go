package rays

import (
	"math"
	"rattata/coordinates"
	"rattata/matrices"
)

type Ray struct {
	OriginPoint coordinates.Coordinate
	Direction   coordinates.Coordinate
}

func NewRay(origin, direction coordinates.Coordinate) Ray {
	if origin.IsAPoint() && direction.IsAVector() {
		return Ray{OriginPoint: origin, Direction: direction}
	}

	panic("inputs are wrong")
}

func (r *Ray) PointAtTime(dir float32) *coordinates.Coordinate {
	scaled_vector := r.Direction.Mul(dir)

	return r.OriginPoint.Add(scaled_vector)
}

type Intersection struct {
	Tvalue float32
	Obj    Shape
}

func NewIntersection(tval float32, obj Shape) Intersection {
	return Intersection{Tvalue: tval, Obj: obj}
}

func Intersections(isections ...Intersection) []Intersection {
	return isections
}

func Hit(intersections []Intersection) (*Intersection, bool) {

	if len(intersections) == 0 {
		return nil, false
	}

	res := intersections[0]
	for _, i := range intersections {
		if res.Tvalue < 0 {
			res = i
		} else if i.Tvalue > 0 && i.Tvalue < res.Tvalue {
			res = i
		}
	}

	if res.Tvalue < 0 {
		return nil, false
	}
	return &res, true
}

func Intersect(shape Shape, ray Ray) []Intersection {

	transformed_ray := Transform(ray, shape.Transformation())

	switch casted_shape := shape.(type) {
	case Sphere:
		return intersectSphere(casted_shape, transformed_ray)
	default:
		return []Intersection{}
	}
}

func intersectSphere(sph Sphere, ray Ray) []Intersection {
	sphere_to_ray := ray.OriginPoint.Sub(&sph.Origin)

	a := ray.Direction.DotP(&ray.Direction)
	b := ray.Direction.DotP(sphere_to_ray) * 2
	c := sphere_to_ray.DotP(sphere_to_ray) - sph.Radius
	discriminant := b*b - 4*a*c

	if discriminant < 0 {
		return []Intersection{}
	}

	t1 := float32((float64(-b) - math.Sqrt(float64(discriminant))) / (2 * float64(a)))
	t2 := float32((float64(-b) + math.Sqrt(float64(discriminant))) / (2 * float64(a)))
	return Intersections(NewIntersection(t1, sph), NewIntersection(t2, sph))
}

func Transform(ray Ray, matrix matrices.Matrix) Ray {
	ray_origin_matrix, ray_direction_matrix := matrices.CoordinateToMatrix(ray.OriginPoint), matrices.CoordinateToMatrix(ray.Direction)

	_, ray_origin_post_transform := matrix.Multiply(ray_origin_matrix)
	_, ray_direction_post_transform := matrix.Multiply(ray_direction_matrix)

	return NewRay(matrices.MatrixToCoordinate(ray_origin_post_transform), matrices.MatrixToCoordinate(ray_direction_post_transform))
}

func ReflectVector(incidence, normal coordinates.Coordinate) coordinates.Coordinate {
	directionInversionVector := normal.Mul(-2 * normal.DotP(&incidence))
	return *incidence.Add(directionInversionVector)
}
