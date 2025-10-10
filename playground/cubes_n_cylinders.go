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

func CreateGoofyWorld() {
	w := observe.NewEmptyWorld()

	light_src := rays.NewLightSource(5, 10, -10, rays.NewLightColour(1, 1, 1))
	w.SetLightSource(&light_src)

	{
		floor := rays.NewPlane(coordinates.CreatePoint(0, 0, 0))
		floor.Material.Pattern = rays.NewPlainPattern(rays.Colour{0.3, 0.4, 0.9})
		floor.Material.Specular = 0
		floor.Material.Reflective = 0.6
		w.AddObject(floor)
	}

	{
		left_wall := rays.NewPlane(coordinates.CreatePoint(0, 0, 0))

		left_wall_transform_mat := matrices.PerformOrderedChainingOps(matrices.NewIdentityMatrix(4),
			matrices.ScalingMatrix(10, 0.01, 10),
			matrices.GivensRotationMatrix3DLeftHanded(coordinates.X, math.Pi/2),
			matrices.GivensRotationMatrix3DLeftHanded(coordinates.Y, -math.Pi/4),
			matrices.TranslationMatrix(0, 0, 5),
		)

		left_wall.SetTransformation(left_wall_transform_mat)
		left_wall.Material.Pattern = rays.NewPlainPattern(rays.Colour{0.6, 0.6, 0.3})
		left_wall.Material.Specular = 0
		w.AddObject(left_wall)
	}

	{
		right_wall := rays.NewPlane(coordinates.CreatePoint(0, 0, 0))

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

		w.AddObject(right_wall)
	}

	{
		left_cube := rays.NewCube()

		left_cube_transform := matrices.PerformOrderedChainingOps(matrices.NewIdentityMatrix(4),
			matrices.ScalingMatrix(0.33, 0.33, 0.33),
			matrices.TranslationMatrix(-1.5, 0.33, -0.75),
		)

		left_cube.SetTransformation(left_cube_transform)

		checker_3d := rays.NewChecker3D(rays.Colour{1, 0.8, 0.1}, rays.Colour{0.4, 0.2, 0.5})
		checker_3d.SetPatternTransformation(matrices.PerformOrderedChainingOps(matrices.NewIdentityMatrix(4),
			matrices.ScalingMatrix(0.5, 0.5, 0.5),
			// matrices.ShearMatrix(1, 0, 0, 0, 1, 1),
		))
		left_cube.Material.Pattern = checker_3d

		// left.Material.Pattern = rays.NewPlainPattern(rays.Colour{1, 0.8, 0.1})
		left_cube.Material.Specular = 0.7
		left_cube.Material.Diffuse = 1
		left_cube.Material.Shininess = 500
		left_cube.Material.Reflective = 0.5
		w.AddObject(left_cube)
	}

	{
		right_cyl := rays.NewXZCylinder()
		right_cyl.Maximum = 2
		right_cyl.Minimum = 0

		right_cyl_transform := matrices.PerformOrderedChainingOps(matrices.NewIdentityMatrix(4),
			matrices.ScalingMatrix(0.8, 1, 0.8),
			matrices.TranslationMatrix(1.5, 0, 1.5),
		)

		right_cyl.SetTransformation(right_cyl_transform)
		right_cyl.Material.Pattern = rays.NewPlainPattern(rays.Colour{0.7, 0.15, 0.4})

		right_cyl.Material.Ambient = 0.4
		right_cyl.Material.Shininess = 100
		w.AddObject(right_cyl)
	}

	cam := observe.CreateNewCamera(400, 300, math.Pi/3)

	view_t := matrices.View_Transform(coordinates.CreatePoint(0, 1, -10), coordinates.CreatePoint(0, -1, 0), coordinates.CreateVector(0, 1, 0))

	cam.SetTransformationMatrix(view_t)
	// my_canvas := observe.Render(cam, w)
	my_canvas := observe.RenderParaller(cam, w, 8)
	fmt.Println(canvas.CanvasToPPMData(my_canvas))

}
