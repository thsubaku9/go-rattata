package coordinates

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

type tupleCtxKey struct{}

func givenATuple(ctx context.Context, x, y, z, w float32) (context.Context, error) {
	return context.WithValue(ctx, tupleCtxKey{}, Coordinate([4]float32{x, y, z, w})), nil
}

func giveAPoint(ctx context.Context, x, y, z float32) (context.Context, error) {
	return context.WithValue(ctx, tupleCtxKey{}, CreatePoint(x, y, z)), nil
}

func giveAVector(ctx context.Context, x, y, z float32) (context.Context, error) {
	return context.WithValue(ctx, tupleCtxKey{}, CreateVector(x, y, z)), nil
}

func checkCoordinateX(ctx context.Context, expected float32) error {
	return checkCoordinate(ctx, X, expected)
}

func checkCoordinateY(ctx context.Context, expected float32) error {
	return checkCoordinate(ctx, Y, expected)
}

func checkCoordinateZ(ctx context.Context, expected float32) error {
	return checkCoordinate(ctx, Z, expected)
}

func checkCoordinateW(ctx context.Context, expected float32) error {
	return checkCoordinate(ctx, W, expected)
}

func checkCoordinate(ctx context.Context, coordinate CoordinateAxis, expected float32) error {
	coord, ok := ctx.Value(tupleCtxKey{}).(Coordinate)
	if !ok {
		return errors.New("No coord found in ctx")
	}

	switch coordinate {
	case X:
		res := coord.Get(X)
		if expected != res {
			return fmt.Errorf("Value mismatch: %f vs %f", expected, res)
		}
	case Y:
		if expected != coord.Get(Y) {
			return errors.New("Value mismatch")
		}
	case Z:
		if expected != coord.Get(Z) {
			return errors.New("Value mismatch")
		}
	case W:
		if expected != coord.Get(W) {
			return errors.New("Value mismatch")
		}
	default:
		return errors.New("Coordinate axis does not exist")
	}

	return nil
}

func checkCoordinateIsPoint(ctx context.Context) error {
	coord, ok := ctx.Value(tupleCtxKey{}).(Coordinate)
	if !ok {
		return errors.New("No coord found in ctx")
	}

	if coord.Get(W) != 1 {
		return errors.New("Coordinate is not a point")
	}
	return nil
}

func checkCoordinateIsVector(ctx context.Context) error {
	coord, ok := ctx.Value(tupleCtxKey{}).(Coordinate)
	if !ok {
		return errors.New("No coord found in ctx")
	}

	if coord.Get(W) != 0 {
		return errors.New("Coordinate is not a vector")
	}
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Given(`^a ← tuple\((-?\d+\.?\d*), (-?\d+\.?\d*), (-?\d+\.?\d*), (-?\d+\.?\d*)\)$`, givenATuple)
	ctx.Then(`^a\.x = (-?\d+\.?\d*)$`, checkCoordinateX)
	ctx.Then(`^a\.y = (-?\d+\.?\d*)$`, checkCoordinateY)
	ctx.Then(`^a\.z = (-?\d+\.?\d*)$`, checkCoordinateZ)
	ctx.Then(`^a\.w = (-?\d+\.?\d*)$`, checkCoordinateW)
	ctx.Then(`^a is a point$`, checkCoordinateIsPoint)
	ctx.Then(`^a is a vector$`, checkCoordinateIsVector)

	ctx.Given(`^p ← point\((-?\d+\.?\d*), (-?\d+\.?\d*), (-?\d+\.?\d*)\)$`, giveAPoint)
	ctx.Then(`^p is a point$`, checkCoordinateIsPoint)

	ctx.Given(`^v ← vector\((-?\d+\.?\d*), (-?\d+\.?\d*), (-?\d+\.?\d*)\)$`, giveAVector)
	ctx.Then(`^v is a vector$`, checkCoordinateIsVector)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../features/coordinates"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func TestOpsAdd(t *testing.T) {
	v1 := CreateVector(1, 2, 3)
	v2 := CreateVector(3, 2, 1)

	v3 := v1.Add(&v2)

	assert.True(t, v3.IsAVector())
	assert.Equal(t, CreateVector(4, 4, 4), *v3)
}

func TestOpsSub(t *testing.T) {
	v1 := CreatePoint(2, 2, 2)
	v2 := CreateVector(1, 1, 1)

	v3 := v1.Sub(&v2)

	assert.True(t, v3.IsAPoint())
	assert.Equal(t, CreatePoint(1, 1, 1), *v3)
}

func TestScalerOps(t *testing.T) {
	v1 := CreateVector(2, 2, 2)

	v2 := v1.Mul(2)
	v3 := v1.Div(2)

	assert.True(t, v2.IsAVector())
	assert.Equal(t, CreateVector(4, 4, 4), *v2)
	assert.True(t, v3.IsAVector())
	assert.Equal(t, CreateVector(1, 1, 1), *v3)
}
