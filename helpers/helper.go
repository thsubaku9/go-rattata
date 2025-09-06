package helpers

import (
	"math"
	"rattata/coordinates"
	"testing"
)

func ApproxEqual(t *testing.T, expected, actual, diff_allowed float64) {
	if math.Abs(expected-actual) <= diff_allowed {
		return
	}
	t.Errorf("Value diff larger than permitted -> %f vs %f", expected, actual)
}

func TestApproxEqualCoordinate(t *testing.T, expected, acutal coordinates.Coordinate, diff_allowed float64) {
	t.Helper()
	for i := range expected {
		ApproxEqual(t, expected[i], acutal[i], diff_allowed)
	}
}
