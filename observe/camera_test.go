package observe

import (
	"math"
	"rattata/helpers"
	"testing"
)

/*
Scenario: The pixel size for a horizontal canvas Given c ← camera(200, 125, π/2)
Then c.pixel_size = 0.01
Scenario: The pixel size for a vertical canvas Given c ← camera(125, 200, π/2)
Then c.pixel_size = 0.01

*/

func TestUnitPerPixelH(t *testing.T) {
	_c := CreateNewCamera(200, 125, math.Pi/2)

	helpers.ApproxEqual(t, 0.01, _c.GetPixelSize(), 0.0001)
}

func TestUnitPerPixelV(t *testing.T) {
	_c := CreateNewCamera(125, 200, math.Pi/2)

	helpers.ApproxEqual(t, 0.01, _c.GetPixelSize(), 0.0001)
}
