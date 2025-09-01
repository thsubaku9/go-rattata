package rays

import (
	"rattata/coordinates"
	"rattata/matrices"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTransformationMatrixInShape(t *testing.T) {

	sph := NewSphere(coordinates.CreatePoint(0, 0, 0), 1)
	assert.Equal(t, matrices.NewIdentityMatrix(3), sph.Transformation())

	rotationMat := matrices.GivensRotationMatrix3D(coordinates.Z, 23)

	(&sph).SetTransformation(rotationMat)

	assert.Equal(t, rotationMat, sph.Transformation())
}
