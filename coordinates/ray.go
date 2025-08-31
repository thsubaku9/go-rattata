package coordinates

import (
	"math"
)

type Ray struct {
	OriginPoint Coordinate
	Direction   Coordinate
}

func NewRay(origin, direction Coordinate) Ray {
	if origin.IsAPoint() && direction.IsAVector() {
		return Ray{OriginPoint: origin, Direction: direction}
	}

	panic("inputs are wrong")
}

func (r *Ray) PointAtTime(dir float32) *Coordinate {
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

func Hit(Intersections []Intersection) (*Intersection, bool) {

	res := Intersections[0]
	for _, i := range Intersections {
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

	switch casted_shape := shape.(type) {
	case Sphere:
		return IntersectSphere(casted_shape, ray)
	default:
		return []Intersection{}

	}
}

func IntersectSphere(sph Sphere, ray Ray) []Intersection {
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
