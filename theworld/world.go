package theworld

import (
	"rattata/coordinates"
	"rattata/matrices"
	"rattata/rays"
	"sort"
)

type World struct {
	lightSource *rays.Light
	objects     []*rays.Shape
}

func NewEmptyWorld() World {
	return World{lightSource: nil, objects: make([]*rays.Shape, 0)}
}

func NewDefaultWorld() World {
	lightSrc := rays.NewLightSource(-10, 10, -10, rays.NewWhiteLightColour())
	s1 := rays.NewSphere(coordinates.CreatePoint(0, 0, 0), 1)
	s1.Material = rays.Material{Colour: rays.Colour{0.8, 1.0, 0.6}, Ambient: 0.1, Diffuse: 0.7, Specular: 0.2, Shininess: 200.0}
	s2 := rays.NewSphere(coordinates.CreatePoint(0, 0, 0), 0.5)
	s2.SetTransformation(matrices.ScalingMatrix(0.5, 0.5, 0.5))
	return World{
		lightSource: &lightSrc,
		objects:     make([]*rays.Shape, 0),
	}
}

func (w *World) LightSource() *rays.Light {
	return w.lightSource
}

func (w *World) AddObject(obj rays.Shape) {
	w.objects = append(w.objects, &obj)
}

func (w World) IntersectWithRay(r rays.Ray) []rays.Intersection {
	res := make([]rays.Intersection, 0)

	for _, obj := range w.objects {
		res = append(res, rays.Intersect(*obj, r)...)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Tvalue <= res[j].Tvalue
	})

	return res
}
