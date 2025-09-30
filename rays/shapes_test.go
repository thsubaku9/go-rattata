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

	sph := NewSphere(coordinates.CreatePoint(0, 0, 0), 1)
	assert.Equal(t, matrices.NewIdentityMatrix(4), sph.Transformation())

	rotationMat := matrices.GivensRotationMatrix3DLeftHanded(coordinates.Z, 23)

	(&sph).SetTransformation(rotationMat)

	assert.Equal(t, rotationMat, sph.Transformation())
}

func TestNormalComputationOnSphere(t *testing.T) {
	sph := NewSphere(coordinates.CreatePoint(0, 0, 0), 1)
	point := coordinates.CreatePoint(0, 0, 1)
	normal := sph.NormalAtPoint(point)
	assert.Equal(t, coordinates.CreateVector(0, 0, 1), normal)

	point = coordinates.CreatePoint(float64(math.Sqrt(3))/3, float64(math.Sqrt(3))/3, float64(math.Sqrt(3))/3)
	normal = sph.NormalAtPoint(point)
	assert.Equal(t, coordinates.CreateVector(float64(math.Sqrt(3))/3, float64(math.Sqrt(3))/3, float64(math.Sqrt(3))/3), normal)

}

func TestNormalComputationOnTransformedSphere(t *testing.T) {
	sph := NewSphere(coordinates.CreatePoint(0, 0, 0), 1)
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
