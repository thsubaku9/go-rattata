package observe

import (
	"math"
	"rattata/coordinates"
	"rattata/helpers"
	"rattata/matrices"
	"testing"
)

func TestUnitPerPixelH(t *testing.T) {
	_c := CreateNewCamera(200, 125, math.Pi/2)

	helpers.ApproxEqual(t, 0.01, _c.GetPixelSize(), 0.0001)
}

func TestUnitPerPixelV(t *testing.T) {
	_c := CreateNewCamera(125, 200, math.Pi/2)

	helpers.ApproxEqual(t, 0.01, _c.GetPixelSize(), 0.0001)
}

func TestRayThroughCenter(t *testing.T) {
	_c := CreateNewCamera(201, 101, math.Pi/2)
	r := _c.RayForPixel(100, 50)

	helpers.TestApproxEqualCoordinate(t, coordinates.CreatePoint(0, 0, 0), r.Origin, 0.0001)
	helpers.TestApproxEqualCoordinate(t, coordinates.CreateVector(0, 0, -1), r.Direction, 0.0001)

}

func TestRayThroughCorner(t *testing.T) {

	_c := CreateNewCamera(201, 101, math.Pi/2)
	r := _c.RayForPixel(0, 0)

	helpers.TestApproxEqualCoordinate(t, coordinates.CreatePoint(0, 0, 0), r.Origin, 0.0001)
	helpers.TestApproxEqualCoordinate(t, coordinates.CreateVector(0.66519, 0.33259, -0.66851), r.Direction, 0.0001)
}

func TestRayWhenCamTransform(t *testing.T) {

	_c := CreateNewCamera(201, 101, math.Pi/2)

	t_mat := matrices.TranslationMatrix(0, -2, 5)
	r_mat := matrices.GivensRotationMatrix3D(coordinates.Y, math.Pi/4)
	_, transform_mat := r_mat.Multiply(t_mat)

	_c.SetTransformationMatrix(transform_mat)

	r := _c.RayForPixel(100, 50)

	helpers.TestApproxEqualCoordinate(t, coordinates.CreatePoint(0, 2, -5), r.Origin, 0.0001)
	// todok -> debug
	helpers.TestApproxEqualCoordinate(t, coordinates.CreateVector(math.Sqrt(2)/2, 0, -math.Sqrt(2)/2), r.Direction, 0.0001)
}
