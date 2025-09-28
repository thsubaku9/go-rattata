package playground

import (
	"fmt"
	"math"
	"rattata/canvas"
	"rattata/coordinates"
	"rattata/matrices"
	"rattata/observe"
	"rattata/rays"
)

func PerformWorldBuildingDefault() {
	w := observe.NewDefaultWorld()
	cam := observe.CreateNewCamera(400, 400, math.Pi/2)
	from := coordinates.CreatePoint(0, 0, -5)
	to := coordinates.CreatePoint(0, 0, 0)
	up := coordinates.CreateVector(0, 1, 0)
	cam.SetTransformationMatrix(matrices.View_Transform(from, to, up))

	my_canvas := observe.RenderParaller(cam, w, 3)
	fmt.Println(canvas.CanvasToPPMData(my_canvas))
}

func PerformWorldBuildingCustom() {
	floor := rays.NewSphere(coordinates.CreatePoint(0, 0, 0), 1)
}

/*
floor ← sphere()
floor.transform ← scaling(10, 0.01, 10) floor.material ← material() floor.material.color ← color(1, 0.9, 0.9) floor.material.specular ← 0

*/
