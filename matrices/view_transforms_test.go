package matrices

import (
	"rattata/coordinates"
	"rattata/helpers"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoTransform(t *testing.T) {
	from := coordinates.CreatePoint(0, 0, 0)
	to := coordinates.CreatePoint(0, 0, -1)
	up := coordinates.CreateVector(0, 1, 0)
	trans := View_Transform(from, to, up)

	assert.Equal(t, NewIdentityMatrix(4), trans)
}

func TestTransformPositiveZ(t *testing.T) {
	from := coordinates.CreatePoint(0, 0, 0)
	to := coordinates.CreatePoint(0, 0, 1)
	up := coordinates.CreateVector(0, 1, 0)
	trans := View_Transform(from, to, up)

	assert.Equal(t, ScalingMatrix(-1, 1, -1), trans)
}

func TestTransformTranslational(t *testing.T) {
	from := coordinates.CreatePoint(0, 0, 5)
	to := coordinates.CreatePoint(0, 0, -1)
	up := coordinates.CreateVector(0, 1, 0)
	trans := View_Transform(from, to, up)

	assert.Equal(t, TranslationMatrix(0, 0, -5), trans)
}

func TestArbitraryTransform(t *testing.T) {
	from := coordinates.CreatePoint(1, 3, 2)
	to := coordinates.CreatePoint(4, -2, 8)
	up := coordinates.CreateVector(1, 1, 0)
	trans := View_Transform(from, to, up)

	expected_trans := [][]float64{{-0.50709, 0.50709, 0.67612, -2.36643}, {0.76772, 0.60609, 0.12122, -2.82843}, {-0.35857, 0.59761, -0.71714, 0.00000}, {0.00000, 0.00000, 0.00000, 1.00000}}

	helpers.TestApproxEqualMatrix(t, expected_trans, trans, 0.0001)
}
