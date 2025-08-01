package matrices

import (
	"rattata/coordinates"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoordinateConv(t *testing.T) {
	p := coordinates.CreateCoordinate(1, 2, 3, 0)
	assert.Equal(t, p, MatrixToCoordinate(CoordinateToMatrix(p)))
}

func TestTranslationConversion(t *testing.T) {
	_matexp := NewIdentityMatrix(4)
	_matexp.Set(0, 3, 1)
	_matexp.Set(1, 3, 1)
	_matexp.Set(2, 3, 1)
	_matres := TranslationMatrix(1, 1, 1)

	assert.True(t, _matres.IsEqual(_matexp))

}

func TestPointTranslation(t *testing.T) {
	p := coordinates.CreatePoint(1, 1, 1)
	expected := coordinates.CreatePoint(1, 0, 2)

	res := PeformMatrixTranslation(CoordinateToMatrix(p), 0, -1, 1)

	assert.Equal(t, expected, MatrixToCoordinate(res))
}

func TestVectorTranslation(t *testing.T) {
	v := coordinates.CreateVector(1, 1, 1)

	res := PeformMatrixTranslation(CoordinateToMatrix(v), 0, -1, 1)
	assert.Equal(t, v, MatrixToCoordinate(res))
}

func TestScalingConversion(t *testing.T) {
	_matexp := NewIdentityMatrix(4)
	_matexp.Set(1, 1, 2)
	_matexp.Set(2, 2, 3)
	_matres := ScalingMatrix(1, 2, 3)

	assert.True(t, _matres.IsEqual(_matexp))

}

func TestCoordinateScaling(t *testing.T) {
	v := coordinates.CreateVector(1, 1, 1)
	exp_res := coordinates.CreateVector(-1, 1, 2)
	act_res := MatrixToCoordinate(PeformMatrixScaling(CoordinateToMatrix(v), -1, 1, 2))

	assert.Equal(t, exp_res, act_res)

}
