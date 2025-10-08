package rays

import (
	"math"
	"rattata/coordinates"
	"rattata/matrices"

	"rattata/helpers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransformationMatrixInShape(t *testing.T) {

	sph := NewCenteredSphere()
	assert.Equal(t, matrices.NewIdentityMatrix(4), sph.Transformation())

	rotationMat := matrices.GivensRotationMatrix3DLeftHanded(coordinates.Z, 23)

	(&sph).SetTransformation(rotationMat)

	assert.Equal(t, rotationMat, sph.Transformation())
}

func TestNormalComputationOnSphere(t *testing.T) {
	sph := NewCenteredSphere()
	point := coordinates.CreatePoint(0, 0, 1)
	normal := sph.NormalAtPoint(point)
	assert.Equal(t, coordinates.CreateVector(0, 0, 1), normal)

	point = coordinates.CreatePoint(float64(math.Sqrt(3))/3, float64(math.Sqrt(3))/3, float64(math.Sqrt(3))/3)
	normal = sph.NormalAtPoint(point)
	assert.Equal(t, coordinates.CreateVector(float64(math.Sqrt(3))/3, float64(math.Sqrt(3))/3, float64(math.Sqrt(3))/3), normal)

}

func TestNormalComputationOnTransformedSphere(t *testing.T) {
	sph := NewCenteredSphere()
	transform_mat, _ := matrices.ScalingMatrix(1, 0.5, 1).Multiply(matrices.GivensRotationMatrix3DLeftHanded(coordinates.Z, math.Pi/5))
	sph.SetTransformation(transform_mat)

	point := coordinates.CreatePoint(0, float64(math.Sqrt(2))/2, -float64(math.Sqrt(2))/2)
	normal := sph.NormalAtPoint(point)
	helpers.TestApproxEqualCoordinate(t, coordinates.CreateVector(0, 0.97014, -0.24254), normal, 0.00001)
}

func TestNormalOnPlane(t *testing.T) {
	pl := NewPlane(coordinates.CreatePoint(0, 0, 0))
	assert.Equal(t, coordinates.CreateVector(0, 1, 0), pl.NormalAtPoint(coordinates.CreatePoint(0, 0, 0)))
}

// ---------------------------------- Intersections ----------------------------------
func Test2IntersectionsWithSphereFromOutside(t *testing.T) {
	r := NewRay(coordinates.CreatePoint(0, 0, -5), coordinates.CreateVector(0, 0, 1))
	sph := NewCenteredSphere()
	xs := Intersect(sph, r)

	assert.Equal(t, 2, len(xs))
	assert.Equal(t, 4.0, xs[0].Tvalue)
	assert.Equal(t, 6.0, xs[1].Tvalue)
}

func Test2IntersectionsWithSphereFromInside(t *testing.T) {
	r := NewRay(coordinates.CreatePoint(0, 0, 0), coordinates.CreateVector(0, 0, 1))
	sph := NewCenteredSphere()
	xs := Intersect(sph, r)

	assert.Equal(t, 2, len(xs))
	assert.Equal(t, -1.0, xs[0].Tvalue)
	assert.Equal(t, 1.0, xs[1].Tvalue)
}

func Test1IntersectionWithSphere(t *testing.T) {
	r := NewRay(coordinates.CreatePoint(0, 1, -5), coordinates.CreateVector(0, 0, 1))
	sph := NewCenteredSphere()
	xs := Intersect(sph, r)

	assert.Equal(t, 2, len(xs))
	assert.Equal(t, float64(5.0), xs[0].Tvalue)
	assert.Equal(t, float64(5.0), xs[1].Tvalue)

}

func Test0IntersectionWithSphere(t *testing.T) {

	r := NewRay(coordinates.CreatePoint(0, 2, -5), coordinates.CreateVector(0, 0, 1))
	sph := NewCenteredSphere()
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

func TestIntersectionsWithCube(t *testing.T) {
	cube := NewCube()

	for _, data := range []struct {
		origin, direction coordinates.Coordinate
		t1, t2            float64
	}{
		{coordinates.CreatePoint(5, 0.5, 0), coordinates.CreateVector(-1, 0, 0), 4, 6},
		{coordinates.CreatePoint(-5, 0.5, 0), coordinates.CreateVector(1, 0, 0), 4, 6},
		{coordinates.CreatePoint(0.5, 5, 0), coordinates.CreateVector(0, -1, 0), 4, 6},
		{coordinates.CreatePoint(0.5, -5, 0), coordinates.CreateVector(0, 1, 0), 4, 6},
		{coordinates.CreatePoint(0.5, 0, 5), coordinates.CreateVector(0, 0, -1), 4, 6},
		{coordinates.CreatePoint(0.5, 0, -5), coordinates.CreateVector(0, 0, 1), 4, 6},
		{coordinates.CreatePoint(0, 0.5, 0), coordinates.CreateVector(0, 0, 1), -1, 1},
	} {
		r := NewRay(data.origin, data.direction)
		xs := cube.IntersectWithRay(r)

		assert.Equal(t, 2, len(xs))
		helpers.ApproxEqual(t, data.t1, xs[0].Tvalue, 0.00001)
		helpers.ApproxEqual(t, data.t2, xs[1].Tvalue, 0.00001)
	}
}

func TestIntersectionsWithCubeMiss(t *testing.T) {
	cube := NewCube()

	for _, data := range []struct {
		origin, direction coordinates.Coordinate
	}{
		{coordinates.CreatePoint(-2, 0, 0), coordinates.CreateVector(0.2673, 0.5345, 0.8018)},
		{coordinates.CreatePoint(0, -2, 0), coordinates.CreateVector(0.8018, 0.2673, 0.5345)},
		{coordinates.CreatePoint(0, 0, -2), coordinates.CreateVector(0.5345, 0.8018, 0.2673)},
		{coordinates.CreatePoint(2, 0, 2), coordinates.CreateVector(0, 0, -1)},
		{coordinates.CreatePoint(0, 2, 2), coordinates.CreateVector(0, -1, 0)},
		{coordinates.CreatePoint(2, 2, 0), coordinates.CreateVector(-1, 0, 0)},
	} {
		r := NewRay(data.origin, data.direction)
		xs := cube.IntersectWithRay(r)

		assert.Equal(t, 0, len(xs))
	}
}

func TestCuberNormal(t *testing.T) {
	cube := NewCube()

	for _, data := range []struct {
		point, normal coordinates.Coordinate
	}{
		{coordinates.CreatePoint(1, 0.5, -0.8), coordinates.CreateVector(1, 0, 0)},
		{coordinates.CreatePoint(-1, -0.2, 0.9), coordinates.CreateVector(-1, 0, 0)},
		{coordinates.CreatePoint(-0.4, 1, -0.1), coordinates.CreateVector(0, 1, 0)},
		{coordinates.CreatePoint(0.3, -1, -0.7), coordinates.CreateVector(0, -1, 0)},
		{coordinates.CreatePoint(-0.6, 0.3, 1), coordinates.CreateVector(0, 0, 1)},
		{coordinates.CreatePoint(0.4, 0.4, -1), coordinates.CreateVector(0, 0, -1)},
		{coordinates.CreatePoint(1, 1, 1), coordinates.CreateVector(1, 0, 0)},
		{coordinates.CreatePoint(-1, -1, -1), coordinates.CreateVector(-1, 0, 0)},
	} {
		n := cube.NormalAtPoint(data.point)
		helpers.TestApproxEqualCoordinate(t, data.normal, n, 0.00001)
	}
}
