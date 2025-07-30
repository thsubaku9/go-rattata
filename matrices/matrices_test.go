package matrices

import (
	"context"
	"errors"
	"strconv"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

type matrixHolderKey struct{}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../features/matrices"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Given(`^the following 4x4 matrix M:$`, theFollowingXMatrixM)
	ctx.Step(`^M\[(\d+),(\d+)\] = (-?\d+\.?\d*)$`, m)
}

func m(ctx context.Context, row, col int, val float32) error {
	mat, ok := ctx.Value(matrixHolderKey{}).(Matrix)

	if !ok {
		return errors.New("Value not found")
	}

	if mat[row][col] != val {
		return errors.New("Value mismatch")
	}
	return nil
}

func theFollowingXMatrixM(ctx context.Context, table *godog.Table) context.Context {

	row_size, col_size := len(table.Rows), len(table.Rows[0].Cells)

	_matrix := NewMatrix(row_size, col_size)

	for i, row := range table.Rows {
		for j, cell_val := range row.Cells {
			val, _ := strconv.ParseFloat(cell_val.Value, 32)
			_matrix.Set(i, j, float32(val))
		}
	}

	return context.WithValue(ctx, matrixHolderKey{}, _matrix)
}

func TestMatrixEquality(t *testing.T) {
	M_A, M_B := NewMatrix(3, 3), NewMatrix(3, 3)

	M_A.Set(1, 1, 5.0)
	M_A.Set(2, 2, -1.1)

	M_B.Set(1, 1, 5.0)
	M_B.Set(2, 2, -1.1)

	assert.True(t, M_A.IsEqual(M_B))
}

func TestMatrixMult(t *testing.T) {
	Matrix_a := Matrix{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 8, 7, 6}, {5, 4, 3, 2}}
	Matrix_b := Matrix{{-2, 1, 2, 3}, {3, 2, 1, -1}, {4, 3, 6, 5}, {1, 2, 7, 8}}
	Matrix_req_res := Matrix{{20, 22, 50, 48}, {44, 54, 114, 108}, {40, 58, 110, 102}, {16, 26, 46, 42}}

	_, Matrix_act_res := Matrix_a.Multiply(Matrix_b)

	assert.True(t, Matrix_req_res.IsEqual(Matrix_act_res))
}

func TestMatrixIdentityMult(t *testing.T) {
	Matrix_a := Matrix{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 8, 7, 6}, {5, 4, 3, 2}}
	Matrix_id4 := NewIdentityMatrix(4)
	_, Matrix_act_res := Matrix_a.Multiply(Matrix_id4)

	assert.True(t, Matrix_a.IsEqual(Matrix_act_res))
}

func TestMatrixTranspose(t *testing.T) {
	Matrix_a := Matrix{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 8, 7, 6}, {5, 4, 3, 2}}
	assert.True(t, Matrix_a.IsEqual(Matrix_a.T().T()))
}

func TestMatrixDeterminant(t *testing.T) {
	Matrix_a := Matrix{{15}}
	err, d := Matrix_a.Determinant()
	assert.Nil(t, err)
	assert.Equal(t, float32(15), d)

	Matrix_a = Matrix{{1, 2}, {3, 4}}
	err, d = Matrix_a.Determinant()
	assert.Nil(t, err)
	assert.Equal(t, float32(-2), d)

	Matrix_a = Matrix{{-5, 0, 1}, {1, -2, 3}, {6, -2, 1}}
	err, d = Matrix_a.Determinant()
	assert.Nil(t, err)
	assert.Equal(t, float32(-10), d)

}

func TestMatrixSubMatrix(t *testing.T) {
	Matrix_a := Matrix{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 8, 7, 6}, {5, 4, 3, 2}}
	Matrix_sub := Matrix_a.SubMatrix(1, 3)

	Matrix_res := Matrix{{1, 2, 3}, {9, 8, 7}, {5, 4, 3}}
	assert.True(t, Matrix_res.IsEqual(Matrix_sub))
}

func TestMatrixMinorAndCofactor3x3(t *testing.T) {
	Matrix_a := Matrix{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 8, 7, 6}, {5, 4, 3, 2}}
	m22 := Matrix_a.Minor(2, 2)
	m31 := Matrix_a.Minor(3, 1)
	c00 := Matrix_a.Cofactor(0, 1)
	c11 := Matrix_a.Cofactor(1, 1)
	assert.Equal(t, float32(0), m22)
	assert.Equal(t, float32(0), m31)
	assert.Equal(t, float32(0), c00)
	assert.Equal(t, float32(0), c11)
}

func TestMatrixIsInvertable(t *testing.T) {
	Matrix_a := Matrix{{9, 6}, {12, 8}}
	d, err := Matrix_a.Determinant()
	assert.Nil(t, err)
	assert.False(t, IsMatrixInvertableBasedOnDeterminant(d))

	Matrix_a = Matrix{{-5, 0, 1}, {1, -2, 3}, {6, -2, 1}}
	d, err = Matrix_a.Determinant()
	assert.Nil(t, err)
	assert.True(t, IsMatrixInvertableBasedOnDeterminant(d))

}

func TestMatrixInverseCalculation(t *testing.T) {
	Matrix_a := Matrix{{-5, 0, 1}, {1, -2, 3}, {6, -2, 1}}
	Matrix_a_inv, _ := Matrix_a.Adj()

	_, id3x3 := Matrix_a.Multiply(Matrix_a_inv)

	assert.True(t, id3x3.IsEqual(NewIdentityMatrix(3)))
}
