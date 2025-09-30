package observe

import (
	"math"
	"rattata/canvas"
	"rattata/coordinates"
	"rattata/matrices"
	"rattata/rays"
	"sync"
)

type Camera struct {
	Hsize            uint32
	Vsize            uint32
	FOV              float64
	Transform_Matrix matrices.Matrix
	half_width       float64
	half_height      float64
}

func CreateNewCamera(hsize, vsize uint32, view_size float64) Camera {
	_c := Camera{Hsize: hsize, Vsize: vsize, FOV: view_size, Transform_Matrix: matrices.NewIdentityMatrix(4)}
	_c.GetPixelSize()

	return _c
}

func (c *Camera) SetTransformationMatrix(m matrices.Matrix) Camera {
	c.Transform_Matrix = m
	return *c
}

func (c *Camera) GetPixelSize() float64 {
	half_view := math.Tan(c.FOV / 2)
	aspect_ratio := float64(c.Hsize) / float64(c.Vsize)

	if aspect_ratio >= 1 {
		c.half_width = half_view
		c.half_height = half_view / aspect_ratio
	} else {
		c.half_width = half_view * aspect_ratio
		c.half_height = half_view

	}
	return (c.half_width * 2) / float64(c.Hsize)
}

func (c *Camera) RayForPixel(px, py int) rays.Ray {
	dist_per_pix_size := c.GetPixelSize()
	transform_inv_matrix, _ := c.Transform_Matrix.Inverse()

	xoffset := (0.5 + float64(px)) * dist_per_pix_size
	yoffset := (0.5 + float64(py)) * dist_per_pix_size

	world_x := c.half_width - xoffset
	world_y := c.half_height - yoffset

	_, pixel := transform_inv_matrix.Multiply(matrices.CoordinateToMatrix(coordinates.CreatePoint(world_x, world_y, -1)))
	_, origin := transform_inv_matrix.Multiply(matrices.CoordinateToMatrix(coordinates.CreatePoint(0, 0, 0)))

	pixel_point := matrices.MatrixToCoordinate(pixel)
	origin_point := matrices.MatrixToCoordinate(origin)
	direction := *pixel_point.Sub(&origin_point).Norm()

	return rays.NewRay(origin_point, direction)
}

func Render(cam Camera, world World) canvas.Canvas {
	my_canvas := canvas.CreateCanvas(cam.Hsize, cam.Vsize)

	for py := 0; py < my_canvas.GetHeight(); py++ {
		for px := 0; px < my_canvas.GetWidth(); px++ {
			_ray := cam.RayForPixel(px, py)
			c := world.Color_At(_ray)
			my_canvas.WritePixel(uint32(px), uint32(py), canvas.RayColorToCanvasColor(c))
		}
	}

	return my_canvas
}

func RenderParaller(cam Camera, world World, parallel_count int) canvas.Canvas {
	my_canvas := canvas.CreateCanvas(cam.Hsize, cam.Vsize)

	data_chan := make(chan [2]int, parallel_count)

	wg := sync.WaitGroup{}

	point_render_work := func(data_stream <-chan [2]int, _cam Camera, _world World, _canvas canvas.Canvas) {
		wg.Add(1)
		for data := range data_stream {
			py, px := data[0], data[1]
			_ray := _cam.RayForPixel(px, py)
			c := _world.Color_At(_ray)
			_canvas.WritePixel(uint32(px), uint32(py), canvas.RayColorToCanvasColor(c))
		}

		wg.Done()
	}

	for range parallel_count {
		go point_render_work(data_chan, cam, world, my_canvas)
	}

	for y := 0; y < my_canvas.GetHeight(); y++ {
		for x := 0; x < my_canvas.GetWidth(); x++ {
			data_chan <- [2]int{y, x}
		}
	}

	close(data_chan)
	wg.Wait()

	return my_canvas
}
