package rays

import (
	"fmt"
	"math"
	"rattata/coordinates"
	"rattata/matrices"
)

type Ray struct {
	Origin    coordinates.Coordinate
	Direction coordinates.Coordinate
}

func NewRay(origin, direction coordinates.Coordinate) Ray {
	if origin.IsAPoint() && direction.IsAVector() {
		return Ray{Origin: origin, Direction: direction}
	}

	panic(fmt.Sprintf("Origin coord : %t and Vector coord :%t", origin.IsAPoint(), direction.IsAVector()))
}

func (r *Ray) PointAtTime(dir float64) *coordinates.Coordinate {
	scaled_vector := r.Direction.Mul(dir)

	return r.Origin.Add(scaled_vector)
}

type Intersection struct {
	Tvalue float64
	Obj    Shape
}

func NewIntersection(tval float64, obj Shape) Intersection {
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

	inv_transform, _ := shape.Transformation().Inverse()
	transformed_ray := Transform(ray, inv_transform)

	switch casted_shape := shape.(type) {
	case Sphere:
		return intersectSphere(casted_shape, transformed_ray)
	default:
		return []Intersection{}
	}
}

func intersectSphere(sph Sphere, ray Ray) []Intersection {
	sphere_to_ray := ray.Origin.Sub(&sph.Origin)

	a := ray.Direction.DotP(&ray.Direction)
	b := ray.Direction.DotP(sphere_to_ray) * 2
	c := sphere_to_ray.DotP(sphere_to_ray) - sph.Radius
	discriminant := b*b - 4*a*c

	if discriminant < 0 {
		return []Intersection{}
	}

	t1 := (-b - math.Sqrt(discriminant)) / (2 * a)
	t2 := (-b + math.Sqrt(discriminant)) / (2 * a)
	return Intersections(NewIntersection(t1, sph), NewIntersection(t2, sph))
}

func Transform(ray Ray, matrix matrices.Matrix) Ray {
	ray_origin_matrix, ray_direction_matrix := matrices.CoordinateToMatrix(ray.Origin), matrices.CoordinateToMatrix(ray.Direction)

	_, ray_origin_post_transform := matrix.Multiply(ray_origin_matrix)
	_, ray_direction_post_transform := matrix.Multiply(ray_direction_matrix)

	return NewRay(matrices.MatrixToCoordinate(ray_origin_post_transform), matrices.MatrixToCoordinate(ray_direction_post_transform))
}

func ReflectVector(incidence, normal coordinates.Coordinate) coordinates.Coordinate {
	directionInversionVector := normal.Mul(-2 * normal.DotP(&incidence))
	return *incidence.Add(directionInversionVector)
}

type Light struct {
	Origin coordinates.Coordinate
	Colour Colour
}

type Colour = [3]float64

func NewLightColour(red, green, blue float64) Colour {
	return Colour{red, green, blue}
}

func NewWhiteLightColour() Colour {
	return NewLightColour(1, 1, 1)
}

func NewLightSource(x, y, z float64, colour Colour) Light {
	return Light{Origin: coordinates.CreatePoint(x, y, z), Colour: colour}
}

type Material struct {
	Colour    Colour
	Ambient   float64
	Diffuse   float64
	Specular  float64
	Shininess float64
}

func CreateDefaultMaterial() Material {
	return Material{Colour: Colour{1, 1, 1}, Ambient: 0.1, Diffuse: 0.9, Specular: 0.9, Shininess: 200.0}
}

/*
	Phong reflection model of lighting

Ambient -> Background lighting;
Diffuse -> Light reflected from matte surface;
Specular -> Reflection of light source
*/
func Lighting(m Material, light Light, pos, eyeVector, normalVector coordinates.Coordinate) Colour {
	effectiveColour := Colour{m.Colour[0] * light.Colour[0],
		m.Colour[1] * light.Colour[1],
		m.Colour[2] * light.Colour[2],
	}

	ambient := Colour{
		effectiveColour[0] * m.Ambient,
		effectiveColour[1] * m.Ambient,
		effectiveColour[2] * m.Ambient,
	}
	var diffuse, specular Colour

	lightVector := *light.Origin.Sub(&pos).Norm()
	light_dot_normal := lightVector.DotP(&normalVector)

	if light_dot_normal < 0 {
		diffuse = Colour{}
	} else {
		diffuse = Colour{
			effectiveColour[0] * m.Diffuse * light_dot_normal,
			effectiveColour[1] * m.Diffuse * light_dot_normal,
			effectiveColour[2] * m.Diffuse * light_dot_normal,
		}
	}

	reflectV := ReflectVector(*lightVector.Negate(), normalVector)
	reflect_dot_eye := reflectV.DotP(&eyeVector)

	if reflect_dot_eye <= 0 || light_dot_normal < 0 {
		specular = Colour{}
	} else {
		factor := math.Pow(reflect_dot_eye, m.Shininess)

		specular = Colour{
			light.Colour[0] * m.Specular * factor,
			light.Colour[1] * m.Specular * factor,
			light.Colour[2] * m.Specular * factor,
		}
	}

	return Colour{
		ambient[0] + diffuse[0] + specular[0],
		ambient[1] + diffuse[1] + specular[1],
		ambient[2] + diffuse[2] + specular[2],
	}
}

type PreCompData struct {
	Tvalue         float64
	Object         Shape
	Point          coordinates.Coordinate
	EyeVector      coordinates.Coordinate
	NormalVector   coordinates.Coordinate
	EyeInsideShape bool
}

func PreparePrecompData(intersection Intersection, r Ray) PreCompData {
	_preComp := PreCompData{Tvalue: intersection.Tvalue, Object: intersection.Obj, Point: *r.PointAtTime(intersection.Tvalue),
		EyeVector: *r.Direction.Negate(), NormalVector: intersection.Obj.NormalAtPoint(*r.PointAtTime(intersection.Tvalue))}

	_preComp.EyeInsideShape = _preComp.EyeVector.DotP(&_preComp.NormalVector) < 0
	return _preComp
}

func (pre PreCompData) Shade_Hit(l Light) Colour {
	return Lighting(pre.Object.GetMaterial(), l, pre.Point, pre.EyeVector, pre.NormalVector)
}
