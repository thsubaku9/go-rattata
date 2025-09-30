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
	cam := observe.CreateNewCamera(400, 200, math.Pi/2)
	from := coordinates.CreatePoint(0, 0, -5)
	to := coordinates.CreatePoint(0, 0, 0)
	up := coordinates.CreateVector(0, 1, 0)
	cam.SetTransformationMatrix(matrices.View_Transform(from, to, up))

	my_canvas := observe.RenderParaller(cam, w, 3)
	fmt.Println(canvas.CanvasToPPMData(my_canvas))
}

func PerformWorldBuildingCustom() {

	my_world := observe.NewEmptyWorld()
	light_src := rays.NewLightSource(-10, 10, -10, rays.NewLightColour(1, 1, 1))
	my_world.SetLightSource(&light_src)

	{
		floor := rays.NewSphere(coordinates.CreatePoint(0, 0, 0), 1)
		floor.SetTransformation(matrices.PerformOrderedChainingOps(matrices.NewIdentityMatrix(4), matrices.ScalingMatrix(10, 0.01, 10)))

		floor.Material.Colour = rays.Colour{1, 0.9, 0.9}
		floor.Material.Specular = 0
		my_world.AddObject(floor)
	}

	{
		left_wall := rays.NewSphere(coordinates.CreatePoint(0, 0, 0), 1)

		left_wall_transform_mat := matrices.PerformOrderedChainingOps(matrices.NewIdentityMatrix(4),
			matrices.ScalingMatrix(10, 0.01, 10),
			matrices.GivensRotationMatrix3DLeftHanded(coordinates.X, math.Pi/2),
			matrices.GivensRotationMatrix3DLeftHanded(coordinates.Y, -math.Pi/4),
			matrices.TranslationMatrix(0, 0, 5),
		)

		left_wall.SetTransformation(left_wall_transform_mat)
		left_wall.Material.Colour = rays.Colour{1, 0.9, 0.9}
		left_wall.Material.Specular = 0
		my_world.AddObject(left_wall)
	}

	{
		right_wall := rays.NewSphere(coordinates.CreatePoint(0, 0, 0), 1)

		right_wall_transform_mat := matrices.PerformOrderedChainingOps(matrices.NewIdentityMatrix(4),
			matrices.ScalingMatrix(10, 0.01, 10),
			matrices.GivensRotationMatrix3DLeftHanded(coordinates.X, math.Pi/2),
			matrices.GivensRotationMatrix3DLeftHanded(coordinates.Y, math.Pi/4),
			matrices.TranslationMatrix(0, 0, 5),
		)

		right_wall.SetTransformation(right_wall_transform_mat)
		right_wall.Material.Colour = rays.NewLightColour(1, 0.9, 0.9)
		right_wall.Material.Specular = 0
		my_world.AddObject(right_wall)
	}

	{
		middle := rays.NewSphere(coordinates.CreatePoint(0, 0, 0), 1)

		middle.SetTransformation(matrices.TranslationMatrix(-0.5, 1, 0.5))
		middle.Material.Colour = rays.NewLightColour(0.1, 1, 0.5)
		middle.Material.Specular = 0.3
		middle.Material.Diffuse = 0.7

		my_world.AddObject(middle)

	}

	{
		right := rays.NewSphere(coordinates.CreatePoint(0, 0, 0), 1)

		right_sph_transform := matrices.PerformOrderedChainingOps(matrices.NewIdentityMatrix(4),
			matrices.ScalingMatrix(0.5, 0.5, 0.5),
			matrices.TranslationMatrix(1.5, 0.5, -0.5),
		)
		right.SetTransformation(right_sph_transform)

		right.Material.Colour = rays.NewLightColour(0.5, 1, 0.1)
		right.Material.Specular = 0.7
		right.Material.Diffuse = 0.3
		my_world.AddObject(right)
	}

	{
		left := rays.NewSphere(coordinates.CreatePoint(0, 0, 0), 1)

		left_sph_transform := matrices.PerformOrderedChainingOps(matrices.NewIdentityMatrix(4),
			matrices.ScalingMatrix(0.33, 0.33, 0.33),
			matrices.TranslationMatrix(-1.5, 0.33, -0.75),
		)
		left.SetTransformation(left_sph_transform)

		left.Material.Colour = rays.NewLightColour(1, 0.8, 0.1)
		left.Material.Specular = 0.7
		left.Material.Diffuse = 0.7

		my_world.AddObject(left)
	}

	cam := observe.CreateNewCamera(300, 200, math.Pi/3)

	view_t := matrices.View_Transform(coordinates.CreatePoint(0, 1.5, -5), coordinates.CreatePoint(0, 1, 0), coordinates.CreateVector(0, 1, 0))

	cam.SetTransformationMatrix(matrices.PerformOrderedChainingOps(matrices.NewIdentityMatrix(4), view_t))
	// matrices.GivensRotationMatrix3DLeftHanded(coordinates.Z, math.Pi/2),
	// matrices.ScalingMatrix(1, -1, 1)))

	my_canvas := observe.RenderParaller(cam, my_world, 5)
	// my_canvas := observe.Render(cam, my_world)
	fmt.Println(canvas.CanvasToPPMData(my_canvas))
}
