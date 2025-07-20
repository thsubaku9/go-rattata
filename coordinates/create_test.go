package coordinates

import (
	"context"
	"testing"

	"github.com/cucumber/godog"
)

type tupleCtxKey struct{}

func givenATuple(ctx context.Context, x, y, z, w float32) (context.Context, error) {
	return context.WithValue(ctx, tupleCtxKey{}, []float32{x, y, z, w}), nil
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features", "coordinates"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Given(`^a ← tuple\((\d+), (\d+), (\d+), (\d+)\)$`, givenATuple)
	ctx.Then(`^a.x = (\d+)$`, nil)
	ctx.Then(`^a.y = (\d+)$`, nil)
	ctx.Then(`^a.z = (\d+)$`, nil)
	ctx.Then(`^a.w = (\d+)$`, nil)
	ctx.Then(`^a is a (\w+)$`, nil)
	ctx.Then(`^a is not a (\w+)$`, nil)

	ctx.Given(`^p ← point(4, -4, 3)$`, nil)
	ctx.Then(`^p = tuple(4, -4, 3, 1)$`, nil)

	ctx.Given(`^v ← vector(4, -4, 3)$`, nil)
	ctx.Then(`^v = tuple(4, -4, 3, 0)$`, nil)
}
