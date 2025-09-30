package observe

import (
	"rattata/coordinates"
	"rattata/matrices"
	"rattata/rays"
	"sort"
)

type PreCompData struct {
	Tvalue         float64
	Object         rays.Shape
	Point          coordinates.Coordinate
	OverPoint      coordinates.Coordinate
	EyeVector      coordinates.Coordinate
	NormalVector   coordinates.Coordinate
	EyeInsideShape bool
}

func PreparePrecompData(intersection rays.Intersection, r rays.Ray) PreCompData {
	_preComp := PreCompData{Tvalue: intersection.Tvalue, Object: intersection.Obj, Point: *r.PointAtTime(intersection.Tvalue),
		EyeVector: *r.Direction.Negate(), NormalVector: intersection.Obj.NormalAtPoint(*r.PointAtTime(intersection.Tvalue))}

	_preComp.EyeInsideShape = _preComp.EyeVector.DotP(&_preComp.NormalVector) < 0

	if _preComp.EyeInsideShape {
		_preComp.NormalVector = *_preComp.NormalVector.Negate()
	}

	raising_vector := _preComp.NormalVector.Mul(rays.EPSILON)
	_preComp.OverPoint = *_preComp.Point.Add(raising_vector)
	return _preComp
}

func (pre PreCompData) Shade_Hit(l rays.Light, w World) rays.Colour {
	return rays.Lighting(pre.Object.GetMaterial(), l, pre.OverPoint, pre.EyeVector, pre.NormalVector, w.IsShadowed(pre.OverPoint))
}

type World struct {
	lightSource *rays.Light
	objects     []rays.Shape
}

func NewEmptyWorld() World {
	return World{lightSource: nil, objects: make([]rays.Shape, 0)}
}

func NewDefaultWorld() World {
	lightSrc := rays.NewLightSource(-10, 10, -10, rays.NewWhiteLightColour())

	s1 := rays.NewSphere(coordinates.CreatePoint(0, 0, 0), 1)
	s1.Material = rays.Material{Colour: rays.Colour{0.8, 1.0, 0.6}, Ambient: 0.1, Diffuse: 0.7, Specular: 0.2, Shininess: 200.0}
	s2 := rays.NewSphere(coordinates.CreatePoint(0, 0, 0), 1)
	s2.SetTransformation(matrices.ScalingMatrix(0.5, 0.5, 0.5))

	_objects := make([]rays.Shape, 0)
	_objects = append(_objects, s1)
	_objects = append(_objects, s2)
	return World{
		lightSource: &lightSrc,
		objects:     _objects,
	}
}

func (w *World) LightSource() *rays.Light {
	return w.lightSource
}

func (w *World) SetLightSource(lightSrc *rays.Light) {
	w.lightSource = lightSrc
}

func (w *World) AddObject(obj rays.Shape) {
	w.objects = append(w.objects, obj)
}

func (w *World) ListObjects() []rays.Shape {
	return w.objects
}

func (w *World) RemoveObjectAt(index int) {
	if index < 0 || index >= len(w.objects) {
		return
	}
	w.objects = append(w.objects[0:index], w.objects[index+1:len(w.objects)]...)
}

func (w *World) ReplaceObjectAt(index int, obj rays.Shape) {
	if index < 0 || index >= len(w.objects) {
		return
	}
	w.objects[index] = obj
}

func (w *World) PerformObjectModifications(index int, ops func(obj rays.Shape) rays.Shape) {
	w.ReplaceObjectAt(index, ops(w.ListObjects()[index]))
}

func (w World) IntersectWithRay(r rays.Ray) []rays.Intersection {
	res := make([]rays.Intersection, 0)

	for _, obj := range w.objects {
		res = append(res, rays.Intersect(obj, r)...)
	}

	sort.Slice(res, func(i, j int) bool {
		return res[i].Tvalue <= res[j].Tvalue
	})

	return res
}

func (w World) Color_At(r rays.Ray) rays.Colour {
	xs := w.IntersectWithRay(r)
	res, isOk := rays.Hit(xs)

	if !isOk {
		return rays.Colour{0, 0, 0}
	}

	precomp := PreparePrecompData(*res, r)

	return precomp.Shade_Hit(*w.LightSource(), w)
}

func (w World) IsShadowed(point coordinates.Coordinate) bool {
	_vec := w.lightSource.Origin.Sub(&point)
	dist := _vec.Magnitude()
	direction := _vec.Norm()

	ray := rays.NewRay(point, *direction)

	xs := w.IntersectWithRay(ray)
	h, doesHit := rays.Hit(xs)

	if doesHit && h.Tvalue < dist {
		return true
	}

	return false
}
