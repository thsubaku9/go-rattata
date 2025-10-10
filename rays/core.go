package rays

import (
	"math"
	"rattata/coordinates"
	"rattata/matrices"
)

var EPSILON float64 = 0.0001
var REC_LIMIT uint = 3

// ------------------------------------- Light and Color struct ------------------------------------
type Colour = [3]float64

type Light struct {
	Origin coordinates.Coordinate
	Colour Colour
}

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

// ------------------------------------- Rays and Intersections ------------------------------------
type Ray struct {
	Origin    coordinates.Coordinate
	Direction coordinates.Coordinate
}

func NewRay(origin, direction coordinates.Coordinate) Ray {
	if origin.IsAPoint() && direction.IsAVector() {
		return Ray{Origin: origin, Direction: direction}
	}

	org_cpy := origin.Copy()
	dir_cpy := direction.Copy()
	org_cpy.Set(coordinates.W, 1)
	dir_cpy.Set(coordinates.W, 0)

	return Ray{Origin: org_cpy, Direction: dir_cpy}
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

	return shape.IntersectWithRay(transformed_ray)
}

func Transform(ray Ray, matrix matrices.Matrix) Ray {
	ray_origin_matrix, ray_direction_matrix := matrices.CoordinateToMatrix(ray.Origin), matrices.CoordinateToMatrix(ray.Direction)

	ray_origin_post_transform, _ := matrix.Multiply(ray_origin_matrix)
	ray_direction_post_transform, _ := matrix.Multiply(ray_direction_matrix)

	return NewRay(matrices.MatrixToCoordinate(ray_origin_post_transform), matrices.MatrixToCoordinate(ray_direction_post_transform))
}

// ------------------------------------- Utility Functions  ------------------------------------

func ReflectVector(incidence, normal coordinates.Coordinate) coordinates.Coordinate {
	directionInversionVector := normal.Mul(-2 * normal.DotP(&incidence))
	return *incidence.Add(directionInversionVector)
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

func RefractiveVector(normal_vector, incidence_vector coordinates.Coordinate, inbound_ri, outbound_ri float64) (*coordinates.Coordinate, bool) {

	//n_inci * sin(incidence) = n_refrac * sin(refraction)

	n_inc_over_n_refrac := inbound_ri / outbound_ri
	cos_incidence := -incidence_vector.DotP(&normal_vector)
	sin_incidence_squared := 1.0 - cos_incidence*cos_incidence
	sin_refraction_squared := n_inc_over_n_refrac * n_inc_over_n_refrac * sin_incidence_squared

	if sin_refraction_squared > 1.0 {
		return nil, false // Total internal reflection
	}

	cos_refraction := math.Sqrt(1.0 - sin_refraction_squared)

	// https://physics.stackexchange.com/questions/435512/snells-law-in-vector-form
	// https://en.wikipedia.org/wiki/Snell%27s_law#Vector_form
	refracted_vector := normal_vector.Mul(n_inc_over_n_refrac*cos_incidence - cos_refraction).
		Add(incidence_vector.Mul(n_inc_over_n_refrac))

	return refracted_vector, true
}

func SchlickReflectiveScore(eye_vector, normal_vector coordinates.Coordinate, inbound_ri, outbound_ri float64) float64 {
	// https://graphics.stanford.edu/courses/cs148-10-summer/docs/2006--degreve--reflection_refraction.pdf
	cos_theta_i := eye_vector.DotP(&normal_vector)

	if inbound_ri > outbound_ri {
		n := inbound_ri / outbound_ri
		sin2_theta_t := n * n * (1.0 - cos_theta_i*cos_theta_i)
		if sin2_theta_t > 1.0 {
			return 1.0 // Total internal reflection
		}

		cos_theta_t := math.Sqrt(1.0 - sin2_theta_t)
		cos_theta_i = cos_theta_t
	}

	r0 := math.Pow(((inbound_ri - outbound_ri) / (inbound_ri + outbound_ri)), 2)
	return r0 + (1-r0)*math.Pow((1-cos_theta_i), 5)
}

// ------------------------------------- Material  ------------------------------------
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
