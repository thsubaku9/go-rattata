package matrices

import (
	"testing"

	"github.com/cucumber/godog"
)

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

func InitializeScenario(ctx *godog.ScenarioContext) {}

func TestMatrixEquality(t *testing.T) {
	t.Fail()
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
