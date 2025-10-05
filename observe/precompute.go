package observe

import (
	"fmt"
	"rattata/coordinates"
	"rattata/rays"
)

type PreCompData struct {
	Tvalue           float64
	Object           rays.Shape
	Point            coordinates.Coordinate
	OverPoint        coordinates.Coordinate
	UnderPoint       coordinates.Coordinate
	EyeVector        coordinates.Coordinate
	NormalVector     coordinates.Coordinate
	ReflectiveVector coordinates.Coordinate
	EyeInsideShape   bool
	RI_Inbound       float64
	RI_Outbound      float64
}

func PreparePrecompData(intersection rays.Intersection, r rays.Ray, xs []rays.Intersection) PreCompData {

	_preComp := PreCompData{Tvalue: intersection.Tvalue, Object: intersection.Obj, Point: *r.PointAtTime(intersection.Tvalue),
		EyeVector: *r.Direction.Negate(), NormalVector: intersection.Obj.NormalAtPoint(*r.PointAtTime(intersection.Tvalue))}

	_preComp.EyeInsideShape = _preComp.EyeVector.DotP(&_preComp.NormalVector) < 0

	if _preComp.EyeInsideShape {
		_preComp.NormalVector = *_preComp.NormalVector.Negate()
	}

	raising_vector := _preComp.NormalVector.Mul(rays.EPSILON)
	_preComp.OverPoint = *_preComp.Point.Add(raising_vector)
	_preComp.UnderPoint = *_preComp.Point.Sub(raising_vector)
	_preComp.ReflectiveVector = rays.ReflectVector(r.Direction, _preComp.NormalVector)

	_preComp.RI_Inbound, _preComp.RI_Outbound = obtainInboundOutboundRefractiveIndices(intersection, xs)
	return _preComp
}

func obtainInboundOutboundRefractiveIndices(intersection rays.Intersection, xs []rays.Intersection) (float64, float64) {

	fmt.Printf("%p \n", &intersection)
	containers := make([]rays.Shape, 0)
	var ri_inbound, ri_outbound float64 = 1.0, 1.0

	for _, intersection1 := range xs {
		if intersection.Tvalue == intersection1.Tvalue && intersection.Obj.Id() == intersection1.Obj.Id() { // Intersection meet point
			if len(containers) == 0 {
				ri_inbound = 1.0
			} else {
				ri_inbound = containers[len(containers)-1].GetMaterial().RefractiveIndex
			}
		}

		// Check if the object is already in the container
		intersection1_index_in_container := -1

		for idx, obj := range containers {
			if obj.Id() == intersection1.Obj.Id() {
				intersection1_index_in_container = idx
				break
			}
		}

		if intersection1_index_in_container != -1 { // Object is already in the container, so remove it
			containers = append(containers[:intersection1_index_in_container], containers[intersection1_index_in_container+1:]...)
		} else {
			containers = append(containers, intersection1.Obj) // Object is not in the container, so add it
		}

		if intersection.Tvalue == intersection1.Tvalue && intersection.Obj.Id() == intersection1.Obj.Id() { // Intersection leave point
			if len(containers) == 0 {
				ri_outbound = 1.0
			} else {
				ri_outbound = containers[len(containers)-1].GetMaterial().RefractiveIndex
			}
		}

	}

	return ri_inbound, ri_outbound
}

func (pre PreCompData) Shade_Hit(l rays.Light, w World, limit uint) rays.Colour {
	lighting_value := rays.Lighting(pre.Object, l, pre.OverPoint, pre.EyeVector, pre.NormalVector, w.IsShadowed(pre.OverPoint))
	reflected_value := pre.Reflected_Colour(w, limit)
	return rays.AddColour(lighting_value, reflected_value)
}

func (pre PreCompData) Reflected_Colour(w World, limit uint) rays.Colour {
	if pre.Object.GetMaterial().Reflective == 0 {
		return rays.Colour{0, 0, 0}
	}

	reflect_ray := rays.NewRay(pre.OverPoint, pre.ReflectiveVector)
	color := w.Color_At(reflect_ray, limit-1)

	return rays.MulColour(color, pre.Object.GetMaterial().Reflective)
}

func (pre PreCompData) Refracted_Colour(w World, limit uint) rays.Colour {
	if pre.Object.GetMaterial().Transparency == 0 || limit == 0 {
		return rays.Colour{0, 0, 0}
	}

	// todo perform refraction calculations here
	return rays.Colour{1, 1, 1}
}
