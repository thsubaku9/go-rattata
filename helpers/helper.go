package helpers

import (
	"math"
	"rattata/coordinates"
	"testing"
)

func ApproxEqual(t *testing.T, expected, actual, diff_allowed float64) {
	t.Helper()
	if math.Abs(expected-actual) <= diff_allowed {
		return
	}
	t.Errorf("Diff aprox: expected(%f) vs actual(%f)", expected, actual)
}

func TestApproxEqualMatrix(t *testing.T, expected, actual [][]float64, diff_allowed float64) {
	t.Helper()
	if len(expected) != len(actual) {
		t.Errorf("Rows mismatch")
	} else if len(expected[0]) != len(actual[0]) {
		t.Errorf("Column mismatch")
	} else {
		for r := range expected {
			for c := range expected[r] {
				ApproxEqual(t, expected[r][c], actual[r][c], diff_allowed)
			}
		}
	}
}

func TestApproxEqualCoordinate(t *testing.T, expected, acutal coordinates.Coordinate, diff_allowed float64) {
	t.Helper()
	for i := range expected {
		ApproxEqual(t, expected[i], acutal[i], diff_allowed)
	}
}
