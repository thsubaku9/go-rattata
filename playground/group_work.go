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

func AHexagon() {

	my_world := observe.NewEmptyWorld()
	light_src := rays.NewLightSource(-10, 10, -10, rays.NewLightColour(1, 1, 1))
	my_world.SetLightSource(&light_src)

	full_hex := hexagon_full()

	my_world.AddObject(full_hex)

	cam := observe.CreateNewCamera(400, 300, math.Pi/3)
	view_t := matrices.View_Transform(coordinates.CreatePoint(0, 2, -3), coordinates.CreatePoint(0, -1, 0), coordinates.CreateVector(0, 1, 0))

	cam.SetTransformationMatrix(matrices.PerformOrderedChainingOps(matrices.NewIdentityMatrix(4), view_t))
	my_canvas := observe.RenderParaller(cam, my_world, 8)
	fmt.Println(canvas.CanvasToPPMData(my_canvas))
}

func hexagon_corner() rays.Sphere {
	corner := rays.NewCenteredSphere()

	transform_mat := matrices.PerformOrderedChainingOps(matrices.NewIdentityMatrix(4),
		matrices.ScalingMatrix(0.25, 0.25, 0.25),
		matrices.TranslationMatrix(0, 0, -1),
	)

	corner.SetTransformation(transform_mat)

	return corner
}

func hexagon_edge() rays.XZCylinder {
	edge := rays.NewXZCylinder()
	edge.Minimum = 0
	edge.Maximum = 1
	edge.Closed = true

	transform_mat := matrices.PerformOrderedChainingOps(matrices.NewIdentityMatrix(4),
		matrices.ScalingMatrix(0.25, 1, 0.25),
		matrices.GivensRotationMatrix3D(coordinates.Z, -math.Pi/2),
		matrices.GivensRotationMatrix3DLeftHanded(coordinates.Y, -math.Pi/6),
		matrices.TranslationMatrix(0, 0, -1),
	)

	edge.SetTransformation(transform_mat)
	return edge
}

func hexagon_side() rays.Group {
	grp := rays.NewGroup()

	corner := hexagon_corner()
	edge := hexagon_edge()
	grp.IndoctrinateShapeToGroup(&edge)

	grp.IndoctrinateShapeToGroup(&corner)

	return grp
}

func hexagon_full() rays.Shape {
	hex := rays.NewGroup()

	for i := 0; i < 6; i++ {
		side := hexagon_side()
		rot_transform := matrices.GivensRotationMatrix3DLeftHanded(coordinates.Y, float64(i)*math.Pi/3)
		side.SetTransformation(rot_transform)

		hex.IndoctrinateShapeToGroup(&side)
	}

	return hex
}
