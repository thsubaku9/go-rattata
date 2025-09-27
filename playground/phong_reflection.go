package playground

import (
	"fmt"
	"rattata/canvas"
	"rattata/coordinates"
	"rattata/matrices"
	"rattata/rays"
)

func ProcessPhongReflection() {

	ray_origin := coordinates.CreatePoint(0, 0, -5)
	wall_z := 10.0
	wall_size := 7.0

	canvas_pixels := 500.0
	pixel_size := wall_size / canvas_pixels
	half := wall_size / 2

	my_canvas := canvas.CreateCanvas(uint32(canvas_pixels), uint32(canvas_pixels))
	color := canvas.NewColour()
	color.SetValue(canvas.Red, 255)
	sph := rays.NewSphere(coordinates.CreatePoint(0, 0, 0), 1)
	sph.Material.Colour = rays.Colour{1, 0.2, 1}
	sph.SetTransformation(matrices.ScalingMatrix(1, 0.7, 1))

	sph2 := rays.NewSphere(coordinates.CreatePoint(1, 1, 2), 0.5)
	sph2.Material.Colour = rays.Colour{0.5, 0.5, 0.0}

	light := rays.NewLightSource(-10, 10, -10, rays.NewWhiteLightColour())

	for y := 0; y < my_canvas.GetHeight(); y++ {
		world_y := half - pixel_size*float64(y)

		for x := 0; x < my_canvas.GetWidth(); x++ {
			world_x := -half + pixel_size*float64(x)

			position := coordinates.CreatePoint(world_x, world_y, wall_z)

			cur_ray := rays.NewRay(ray_origin, *position.Sub(&ray_origin).Norm())
			xs := rays.Intersect(sph, cur_ray)

			if intersect, isPresent := rays.Hit(xs); isPresent {
				point := cur_ray.PointAtTime(intersect.Tvalue)
				normal_vector := sph.NormalAtPoint(*point)
				eye_vector := *cur_ray.Direction.Negate()

				color := rays.Lighting(sph.Material, light, *point, eye_vector, normal_vector)
				my_canvas.WritePixel(uint32(x), uint32(y), canvas.RayColorToCanvasColor(color))
			}

			xs2 := rays.Intersect(sph2, cur_ray)
			if intersect, isPresent := rays.Hit(xs2); isPresent {
				point := cur_ray.PointAtTime(intersect.Tvalue)
				normal_vector := sph2.NormalAtPoint(*point)
				eye_vector := *cur_ray.Direction.Negate()

				color := rays.Lighting(sph2.Material, light, *point, eye_vector, normal_vector)
				my_canvas.WritePixel(uint32(x), uint32(y), canvas.RayColorToCanvasColor(color))
			}
		}
	}
	fmt.Println(canvas.CanvasToPPMData(my_canvas))

}
