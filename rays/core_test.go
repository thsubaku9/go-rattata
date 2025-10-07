package rays

import (
	"math"
	"rattata/coordinates"
	"rattata/helpers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test2IntersectionsWithSphereFromOutside(t *testing.T) {
	r := NewRay(coordinates.CreatePoint(0, 0, -5), coordinates.CreateVector(0, 0, 1))
	sph := NewSphere(coordinates.CreatePoint(0, 0, 0), 1)
	xs := Intersect(sph, r)

	assert.Equal(t, 2, len(xs))
	assert.Equal(t, 4.0, xs[0].Tvalue)
	assert.Equal(t, 6.0, xs[1].Tvalue)
}

func Test2IntersectionsWithSphereFromInside(t *testing.T) {
	r := NewRay(coordinates.CreatePoint(0, 0, 0), coordinates.CreateVector(0, 0, 1))
	sph := NewSphere(coordinates.CreatePoint(0, 0, 0), 1)
	xs := Intersect(sph, r)

	assert.Equal(t, 2, len(xs))
	assert.Equal(t, -1.0, xs[0].Tvalue)
	assert.Equal(t, 1.0, xs[1].Tvalue)
}

func Test1IntersectionWithSphere(t *testing.T) {
	r := NewRay(coordinates.CreatePoint(0, 1, -5), coordinates.CreateVector(0, 0, 1))
	sph := NewSphere(coordinates.CreatePoint(0, 0, 0), 1)
	xs := Intersect(sph, r)

	assert.Equal(t, 2, len(xs))
	assert.Equal(t, float64(5.0), xs[0].Tvalue)
	assert.Equal(t, float64(5.0), xs[1].Tvalue)

}

func Test0IntersectionWithSphere(t *testing.T) {

	r := NewRay(coordinates.CreatePoint(0, 2, -5), coordinates.CreateVector(0, 0, 1))
	sph := NewSphere(coordinates.CreatePoint(0, 0, 0), 1)
	xs := Intersect(sph, r)

	assert.Equal(t, 0, len(xs))
}

func Test0IntersectionWithXZPlane(t *testing.T) {

	r := NewRay(coordinates.CreatePoint(0, 1, 0), coordinates.CreateVector(0, 0, 1))
	pl := NewPlane(coordinates.CreatePoint(0, 0, 0))
	xs := Intersect(pl, r)

	assert.Equal(t, 0, len(xs))
}

func Test1IntersectionWithXZPlane(t *testing.T) {

	r := NewRay(coordinates.CreatePoint(0, 5, -1), coordinates.CreateVector(0, -1, 1))
	pl := NewPlane(coordinates.CreatePoint(0, 0, 0))
	xs := Intersect(pl, r)

	assert.Equal(t, 1, len(xs))
	helpers.ApproxEqual(t, 5.0, xs[0].Tvalue, 0.00001)
}

func TestHit(t *testing.T) {
	s := NewSphere(coordinates.CreatePoint(0, 0, 0), 1)
	i1 := NewIntersection(5, s)
	i2 := NewIntersection(7, s)
	i3 := NewIntersection(-3, s)
	i4 := NewIntersection(2, s)
	xs := Intersections(i1, i2, i3, i4)
	i, res := Hit(xs)
	assert.True(t, res)
	assert.Equal(t, i4, *i)

}

func TestStandardReflection(t *testing.T) {
	v := coordinates.CreateVector(1, -1, 0)
	n := coordinates.CreateVector(0, 1, 0)
	r := ReflectVector(v, n)
	helpers.TestApproxEqualCoordinate(t, coordinates.CreateVector(1, 1, 0), r, 0.00001)
}

func TestTiltedReflection(t *testing.T) {
	v := coordinates.CreateVector(0, -1, 0)
	n := coordinates.CreateVector(float64(math.Sqrt(2)/2), float64(math.Sqrt(2)/2), 0)
	r := ReflectVector(v, n)
	helpers.TestApproxEqualCoordinate(t, coordinates.CreateVector(1, 0, 0), r, 0.00001)
}

func TestLightingScenarioEyeBetweenLightNSurface(t *testing.T) {
	sph := NewSphere(coordinates.CreatePoint(0, 0, 0), 1)
	position := coordinates.CreatePoint(0, 0, 0)
	eyeV := coordinates.CreateVector(0, 0, -1)
	normalV := coordinates.CreateVector(0, 0, -1)
	light := NewLightSource(0, 0, -10, NewWhiteLightColour())
	result := Lighting(sph, light, position, eyeV, normalV, false)

	assert.Equal(t, Colour{1.9, 1.9, 1.9}, result)

}

func TestLightingScenarioEyeBetweenLightNSurfaceEyeOffset45Deg(t *testing.T) {
	sph := NewSphere(coordinates.CreatePoint(0, 0, 0), 1)
	position := coordinates.CreatePoint(0, 0, 0)
	eyeV := coordinates.CreateVector(0, math.Sqrt(2)/2, -math.Sqrt(2)/2)
	normalV := coordinates.CreateVector(0, 0, -1)
	light := NewLightSource(0, 0, -10, NewWhiteLightColour())
	result := Lighting(sph, light, position, eyeV, normalV, false)

	assert.Equal(t, Colour{1, 1, 1}, result)
}

func TestLightingScenarioEyeOppositeSurface(t *testing.T) {
	sph := NewSphere(coordinates.CreatePoint(0, 0, 0), 1)
	position := coordinates.CreatePoint(0, 0, 0)
	eyeV := coordinates.CreateVector(0, 0, -1)
	normalV := coordinates.CreateVector(0, 0, -1)
	light := NewLightSource(0, 10, -10, NewWhiteLightColour())
	result := Lighting(sph, light, position, eyeV, normalV, false)

	assert.Equal(t, Colour{0.7363961030678927, 0.7363961030678927, 0.7363961030678927}, result)
}

func TestLightingScenarioEyeInReflectionPath(t *testing.T) {
	sph := NewSphere(coordinates.CreatePoint(0, 0, 0), 1)
	position := coordinates.CreatePoint(0, 0, 0)
	eyeV := coordinates.CreateVector(0, -float64(math.Sqrt(2)/2), -float64(math.Sqrt(2)/2))
	normalV := coordinates.CreateVector(0, 0, -1)
	light := NewLightSource(0, 10, -10, NewWhiteLightColour())
	result := Lighting(sph, light, position, eyeV, normalV, false)

	assert.Equal(t, Colour{1.6363961030678928, 1.6363961030678928, 1.6363961030678928}, result)
}

func TestLightingScenarioLightSourceBehindSurface(t *testing.T) {
	sph := NewSphere(coordinates.CreatePoint(0, 0, 0), 1)
	position := coordinates.CreatePoint(0, 0, 0)
	eyeV := coordinates.CreateVector(0, 0, -1)
	normalV := coordinates.CreateVector(0, 0, -1)
	light := NewLightSource(0, 0, 10, NewWhiteLightColour())
	result := Lighting(sph, light, position, eyeV, normalV, false)

	assert.Equal(t, Colour{0.1, 0.1, 0.1}, result)
}

func TestLightingShadowRegion(t *testing.T) {
	sph := NewSphere(coordinates.CreatePoint(0, 0, 0), 1)
	position := coordinates.CreatePoint(0, 0, 0)
	eyeV := coordinates.CreateVector(0, 0, -1)
	normalV := coordinates.CreateVector(0, 0, -1)
	light := NewLightSource(0, 0, -10, NewWhiteLightColour())
	result := Lighting(sph, light, position, eyeV, normalV, true)

	assert.Equal(t, Colour{0.1, 0.1, 0.1}, result)
}

// todo
func TestStandardRefraction(t *testing.T) {
	v := coordinates.CreateVector(1, -1, 0)
	n := coordinates.CreateVector(0, 1, 0)
	r := ReflectVector(v, n)
	helpers.TestApproxEqualCoordinate(t, coordinates.CreateVector(1, 1, 0), r, 0.00001)
}

func TestTiltedRefraction(t *testing.T) {
	v := coordinates.CreateVector(0, -1, 0)
	n := coordinates.CreateVector(float64(math.Sqrt(2)/2), float64(math.Sqrt(2)/2), 0)
	r := ReflectVector(v, n)
	helpers.TestApproxEqualCoordinate(t, coordinates.CreateVector(1, 0, 0), r, 0.00001)
}

func TestTotalInternalReflection(t *testing.T) {
	v := coordinates.CreateVector(1, -1, 0)
	n := coordinates.CreateVector(0, 1, 0)
	r := ReflectVector(v, n)
	helpers.TestApproxEqualCoordinate(t, coordinates.CreateVector(1, 1, 0), r, 0.00001)
}

func TestSchlickUnderTIR(t *testing.T) {

	eyeVector := coordinates.CreateVector(0, 0, -1)
	normalV := coordinates.CreateVector(0, -1, 0)
	n1 := 1.5
	n2 := 1.0

	schlick_val := SchlickReflectiveScore(eyeVector, normalV, n1, n2)

	assert.Equal(t, 1.0, schlick_val)

}

func TestSchlickForPerpendicularRay(t *testing.T) {
	eyeVector := coordinates.CreateVector(0, 1, 0)
	normalV := coordinates.CreateVector(0, 1, 0)
	n1 := 1.5
	n2 := 1.0

	schlick_val := SchlickReflectiveScore(eyeVector, normalV, n1, n2)

	helpers.ApproxEqual(t, 0.04, schlick_val, 0.001)
}
func TestSchlickFor_N2_GT_N1(t *testing.T) {

	eyeVector := coordinates.CreateVector(0, 0, 1)
	normalV := coordinates.CreateVector(0, 0.95, 0.31)
	n1 := 1.0
	n2 := 1.5

	schlick_val := SchlickReflectiveScore(eyeVector, normalV, n1, n2)

	helpers.ApproxEqual(t, 0.190147, schlick_val, 0.001)

}
