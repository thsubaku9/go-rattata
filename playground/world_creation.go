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
	cam := observe.CreateNewCamera(500, 400, math.Pi/2)
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
		floor := rays.NewPlane(coordinates.CreatePoint(0, 0, 0))
		floor.SetTransformation(matrices.ScalingMatrix(10, 0.01, 10))
		floor.Material.Pattern = rays.NewPlainPattern(rays.Colour{1, 0.9, 0.9})
		floor.Material.Specular = 0
		floor.Material.Reflective = 0.4
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
		left_wall.Material.Pattern = rays.NewPlainPattern(rays.Colour{0.6, 0.6, 0.3})
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
		right_wall.Material.Pattern = rays.NewPlainPattern(rays.Colour{1, 0.9, 0.9})
		right_wall.Material.Specular = 0
		pat := rays.NewChecker3D(rays.Colour{1, 0.8, 0.1}, rays.Colour{0.6, 0.3, 0.0})
		pat.SetPatternTransformation(matrices.ScalingMatrix(0.2, 0.2, 0.2))
		right_wall.Material.Pattern = pat

		my_world.AddObject(right_wall)
	}

	{
		middle := rays.NewSphere(coordinates.CreatePoint(0, 0, 0), 1)

		middle.SetTransformation(matrices.TranslationMatrix(-0.5, 1, 0.5))
		// middle.Material.Pattern = rays.NewPlainPattern(rays.Colour{0.1, 1, 0.5})

		_pat := rays.NewXStripe(rays.Colour{0.1, 0.5, 0.5}, rays.Colour{0.8, 0.2, 0.4})
		_pat.SetPatternTransformation(matrices.PerformOrderedChainingOps(matrices.NewIdentityMatrix(4),
			matrices.ScalingMatrix(0.2, 0.2, 0.2),
			matrices.GivensRotationMatrix3DLeftHanded(coordinates.Z, math.Pi/2),
		))

		_pat2 := rays.NewPerturbedPattern(_pat, 0.5)
		middle.Material.Pattern = _pat2
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

		grad_p := rays.NewXGradient(rays.Colour{0.9, 0.0, 0.9}, rays.Colour{0.5, 0.8, 0.2})
		grad_p.SetPatternTransformation(matrices.PerformOrderedChainingOps(matrices.NewIdentityMatrix(4),
			matrices.ScalingMatrix(3, 3, 3),
			matrices.GivensRotationMatrix3DLeftHanded(coordinates.Z, math.Pi/4),
		))

		right.Material.Pattern = grad_p
		// right.Material.Pattern = rays.NewPlainPattern(rays.Colour{0.5, 1, 0.1})
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

		// checker_3d := rays.NewChecker3D(rays.Colour{1, 0.8, 0.1}, rays.Colour{1, 1, 1})
		// checker_3d.SetPatternTransformation(matrices.PerformOrderedChainingOps(matrices.NewIdentityMatrix(4),
		// 	matrices.ScalingMatrix(0.5, 0.5, 0.5),
		// ))
		// left.Material.Pattern = checker_3d

		uvchecker_pattern := rays.NewUnitSphereUVChecker(rays.Colour{1, 1, 1}, rays.Colour{1, 0.8, 0.1}, 16, 16)
		uvchecker_pattern.SetPatternTransformation(matrices.PerformOrderedChainingOps(matrices.NewIdentityMatrix(4),
			matrices.ScalingMatrix(1, 1, 1),
		))
		left.Material.Pattern = uvchecker_pattern

		// left.Material.Pattern = rays.NewPlainPattern(rays.Colour{1, 0.8, 0.1})
		left.Material.Specular = 0.7
		left.Material.Diffuse = 0.7

		my_world.AddObject(left)
	}

	cam := observe.CreateNewCamera(400, 300, math.Pi/3)

	view_t := matrices.View_Transform(coordinates.CreatePoint(0, 1.5, -5), coordinates.CreatePoint(0, 1, 0), coordinates.CreateVector(0, 1, 0))

	cam.SetTransformationMatrix(matrices.PerformOrderedChainingOps(matrices.NewIdentityMatrix(4), view_t))
	my_canvas := observe.RenderParaller(cam, my_world, 8)
	fmt.Println(canvas.CanvasToPPMData(my_canvas))
}
