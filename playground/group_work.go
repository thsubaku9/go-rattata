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

	// Adjust light source position for better illumination
	light_src := rays.NewLightSource(-3, 5, -3, rays.NewLightColour(1, 1, 1))
	my_world.SetLightSource(&light_src)

	// Create the hexagon
	full_hex := hexagon_full()
	my_world.AddObject(full_hex)

	// Adjust camera position and orientation for a top-inclined view
	cam := observe.CreateNewCamera(400, 300, math.Pi/3)
	view_t := matrices.View_Transform(
		coordinates.CreatePoint(0, 4, -2), // Camera position (higher and farther back)
		coordinates.CreatePoint(0, 0, 0),  // Look at the center of the hexagon
		coordinates.CreateVector(0, 1, 0), // Up direction
	)
	cam.SetTransformationMatrix(view_t)

	// Render the scene
	my_canvas := observe.RenderParaller(cam, my_world, 8)
	fmt.Println(canvas.CanvasToPPMData(my_canvas))
}

func hexagon_corner_a() rays.Sphere {
	corner := rays.NewCenteredSphere()

	transform_mat := matrices.PerformOrderedChainingOps(matrices.NewIdentityMatrix(4),
		matrices.ScalingMatrix(0.25, 0.25, 0.25),
		matrices.TranslationMatrix(0, 0, -1),
	)

	corner.SetTransformation(transform_mat)

	return corner
}
func hexagon_corner_b() rays.Sphere {
	corner := rays.NewCenteredSphere()

	transform_mat := matrices.PerformOrderedChainingOps(matrices.NewIdentityMatrix(4),
		matrices.ScalingMatrix(0.25, 0.25, 0.25),
		matrices.TranslationMatrix(0, 0, 1),
	)

	corner.SetTransformation(transform_mat)

	return corner
}

func hexagon_edge() rays.XZCylinder {
	edge := rays.NewXZCylinder()
	edge.Minimum = 0
	edge.Maximum = 1

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

	corner_a := hexagon_corner_a()
	corner_b := hexagon_corner_b()
	edge := hexagon_edge()
	grp.IndoctrinateShapeToGroup(&edge)

	grp.IndoctrinateShapeToGroup(&corner_a)
	grp.IndoctrinateShapeToGroup(&corner_b)

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
