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

	ctx.Given(`^the following matrix A:$`, theFollowingMatrixA)
	ctx.Given(`^the following matrix B:$`, theFollowingMatrixB)
	ctx.Then(`^A \* B is the following 4x4 matrix:$`, aBIsTheFollowingXMatrix)
}

func aBIsTheFollowingXMatrix(arg1, arg2 int, arg3 *godog.Table) error {
	return godog.ErrPending
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

func theFollowingMatrixA(arg1 *godog.Table) error {
	return godog.ErrPending
}

func theFollowingMatrixB(arg1 *godog.Table) error {
	return godog.ErrPending
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

func TestMatrixIdentityMult(t *testing.T) {
	t.Fail()
}

func TestMatrixTranspose(t *testing.T) {
	t.Fail()
}

func TestMatrixDeterminant2x2(t *testing.T) {
	t.Fail()
}

func TestMatrixSubMatrix(t *testing.T) {
	// submatrix(A, r, c) -> returns a matrix which has row r and col c removed from A
	t.Fail()
}

func TestMatrixMinorAndCofactor3x3(t *testing.T) {
	t.Fail()
}

func TestMatrixIsInvertable(t *testing.T) {
	// if determinant is 0 => not invertable
	t.Fail()
}

func TestMatrixInverseCalculation(t *testing.T) {
	t.Fail()
}
