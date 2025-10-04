package rays

import (
	"fmt"
	"math"
	"rattata/coordinates"
	"rattata/matrices"
)

var EPSILON float64 = 0.0001
var REC_LIMIT uint = 3

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
	case XZPlane:
		return intersectPlane(casted_shape, transformed_ray)
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

func intersectPlane(pl XZPlane, ray Ray) []Intersection {
	if math.Abs(ray.Direction[coordinates.Y]) < EPSILON {
		return []Intersection{}
	}

	t := -ray.Origin[coordinates.Y] / ray.Direction[coordinates.Y]
	return Intersections(NewIntersection(t, pl))
}

func Transform(ray Ray, matrix matrices.Matrix) Ray {
	ray_origin_matrix, ray_direction_matrix := matrices.CoordinateToMatrix(ray.Origin), matrices.CoordinateToMatrix(ray.Direction)

	ray_origin_post_transform, _ := matrix.Multiply(ray_origin_matrix)
	ray_direction_post_transform, _ := matrix.Multiply(ray_direction_matrix)

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

func AddColour(c1, c2 Colour) Colour {
	c3 := Colour{}

	for i := 0; i < 3; i++ {
		c3[i] = c1[i] + c2[i]
	}

	return c3
}

func SubColour(c1, c2 Colour) Colour {
	c3 := Colour{}

	for i := 0; i < 3; i++ {
		c3[i] = c1[i] - c2[i]
	}
	return c3
}

func MulColour(c1 Colour, k float64) Colour {
	c3 := Colour{}

	for i := 0; i < 3; i++ {
		_t := c1[i] * k
		c3[i] = float64(min((_t), 255))
	}
	return c3
}

func NewLightSource(x, y, z float64, colour Colour) Light {
	return Light{Origin: coordinates.CreatePoint(x, y, z), Colour: colour}
}

type Material struct {
	Ambient         float64
	Diffuse         float64
	Specular        float64
	Shininess       float64
	Pattern         Pattern
	Reflective      float64
	Transparency    float64
	RefractiveIndex float64
}

func CreateDefaultMaterial() Material {
	return Material{Pattern: PlainPattern{Colour{1, 1, 1}}, Ambient: 0.1, Diffuse: 0.9, Specular: 0.9, Shininess: 200.0,
		Reflective: 0.0, RefractiveIndex: 1.0, Transparency: 0.0}
}

func PatternAtPoint(world_point coordinates.Coordinate, objectTransformation matrices.Matrix, pattern Pattern) Colour {
	objectTransformationInverse, _ := objectTransformation.Inverse()
	object_point := matrices.PerformOrderedChainingOps(matrices.CoordinateToMatrix(world_point), objectTransformationInverse)
	return pattern.PatternAt(matrices.MatrixToCoordinate(object_point))
}

/*
	Phong reflection model of lighting

Ambient -> Background lighting;
Diffuse -> Light reflected from matte surface;
Specular -> Reflection of light source
*/
func Lighting(shp Shape, light Light, pos, eyeVector, normalVector coordinates.Coordinate, isInShadow bool) Colour {

	m := shp.GetMaterial()

	_point_color := PatternAtPoint(pos, shp.Transformation(), m.Pattern)
	effectiveColour := Colour{_point_color[0] * light.Colour[0],
		_point_color[1] * light.Colour[1],
		_point_color[2] * light.Colour[2],
	}

	ambient := Colour{
		effectiveColour[0] * m.Ambient,
		effectiveColour[1] * m.Ambient,
		effectiveColour[2] * m.Ambient,
	}
	var diffuse, specular Colour = Colour{}, Colour{}

	lightVector := *light.Origin.Sub(&pos).Norm()
	light_dot_normal := lightVector.DotP(&normalVector)

	if light_dot_normal >= 0 && !isInShadow {
		diffuse = Colour{
			effectiveColour[0] * m.Diffuse * light_dot_normal,
			effectiveColour[1] * m.Diffuse * light_dot_normal,
			effectiveColour[2] * m.Diffuse * light_dot_normal,
		}
	}

	reflectV := ReflectVector(*lightVector.Negate(), normalVector)
	reflect_dot_eye := reflectV.DotP(&eyeVector)

	if reflect_dot_eye > 0 && light_dot_normal >= 0 && !isInShadow {
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
