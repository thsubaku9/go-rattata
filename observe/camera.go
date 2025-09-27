package observe

import (
	"math"
	"rattata/canvas"
	"rattata/coordinates"
	"rattata/matrices"
	"rattata/rays"
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
	xoffset := (0.5 + float64(px)) * dist_per_pix_size
	yoffset := (0.5 + float64(py)) * dist_per_pix_size

	world_x := c.half_width - xoffset
	world_y := c.half_height - yoffset

	transform_inv_matrix, _ := c.Transform_Matrix.Inverse()
	_, pixel := transform_inv_matrix.Multiply(matrices.CoordinateToMatrix(coordinates.CreatePoint(world_x, world_y, -1)))
	_, origin := transform_inv_matrix.Multiply(matrices.CoordinateToMatrix(coordinates.CreatePoint(0, 0, 0)))

	pixel_point := matrices.MatrixToCoordinate(pixel)
	origin_point := matrices.MatrixToCoordinate(origin)
	direction := *pixel_point.Sub(&origin_point).Norm()

	return rays.NewRay(pixel_point, direction)
}

func Render(cam Camera, _world World) canvas.Canvas {
	_canvas := canvas.CreateCanvas(cam.Hsize, cam.Vsize)

	for y := 0; y < _canvas.GetHeight(); y++ {
		for x := 0; x < _canvas.GetWidth(); x++ {
			_ray := cam.RayForPixel(x, y)
			c := _world.Color_At(_ray)
			_canvas.WritePixel(uint32(x), uint32(y), canvas.RayColorToCanvasColor(c))
		}
	}

	return _canvas
}
