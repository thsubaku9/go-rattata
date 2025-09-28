package playground

import (
	"fmt"
	"rattata/canvas"
	"rattata/coordinates"
	"rattata/rays"
)

func ProcessRaySphereIntersection() {
	ray_origin := coordinates.CreatePoint(0, 0, -5)
	wall_z := 10.0
	wall_size := 7.0

	canvas_pixels := 200.0
	pixel_size := wall_size / canvas_pixels
	half := wall_size / 2

	my_canvas := canvas.CreateCanvas(uint32(canvas_pixels), uint32(canvas_pixels))
	color := canvas.NewColour()
	color.SetValue(canvas.Red, 255)
	sph := rays.NewSphere(coordinates.CreatePoint(0, 0, 0), 1)

	for y := 0; y < my_canvas.GetHeight(); y++ {
		world_y := half - pixel_size*float64(y)

		for x := 0; x < my_canvas.GetWidth(); x++ {
			world_x := -half + pixel_size*float64(x)

			position := coordinates.CreatePoint(world_x, world_y, wall_z)

			r := rays.NewRay(ray_origin, *position.Sub(&ray_origin).Norm())
			xs := rays.Intersect(sph, r)

			if _, isPresent := rays.Hit(xs); isPresent {
				my_canvas.WritePixel(uint32(x), uint32(y), color)
			}
		}
	}
	fmt.Println(canvas.CanvasToPPMData(my_canvas))
}
