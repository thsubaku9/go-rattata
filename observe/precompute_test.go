package observe

import (
	"math"
	"rattata/coordinates"
	"rattata/helpers"
	"rattata/matrices"
	"rattata/rays"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrecompData(t *testing.T) {
	r := rays.NewRay(coordinates.CreatePoint(0, 0, -5), coordinates.CreateVector(0, 0, 1))
	s := rays.NewSphere(coordinates.CreatePoint(0, 0, 0), 1)
	i := rays.Intersection{Tvalue: 4, Obj: s}

	pre := PreparePrecompData(i, r, []rays.Intersection{i})

	assert.Equal(t, i.Tvalue, pre.Tvalue)
	assert.Equal(t, i.Obj, pre.Object)
	assert.Equal(t, coordinates.CreatePoint(0, 0, -1), pre.Point)
	assert.Equal(t, coordinates.CreateVector(0, 0, -1), pre.EyeVector)
	assert.Equal(t, coordinates.CreateVector(0, 0, -1), pre.NormalVector)
	assert.False(t, pre.EyeInsideShape)

	r2 := rays.NewRay(coordinates.CreatePoint(0, 0, 0), coordinates.CreateVector(0, 0, 1))
	i = rays.Intersection{Tvalue: 1, Obj: s}
	pre2 := PreparePrecompData(i, r2, []rays.Intersection{i})

	assert.Equal(t, i.Tvalue, pre2.Tvalue)
	assert.Equal(t, i.Obj, pre2.Object)
	assert.Equal(t, coordinates.CreatePoint(0, 0, 1), pre2.Point)
	assert.Equal(t, coordinates.CreateVector(0, 0, -1), pre2.EyeVector)
	assert.Equal(t, coordinates.CreateVector(0, 0, -1), pre2.NormalVector)
	assert.True(t, pre2.EyeInsideShape)
}

func TestPrecompRefractiveIndices(t *testing.T) {
	A := rays.NewGlassSphere()
	A.SetTransformation(matrices.ScalingMatrix(2, 2, 2))
	A.Material.RefractiveIndex = 1.5

	B := rays.NewGlassSphere()
	B.SetTransformation(matrices.TranslationMatrix(0, 0, -0.25))
	B.Material.RefractiveIndex = 2.0

	C := rays.NewGlassSphere()
	C.SetTransformation(matrices.TranslationMatrix(0, 0, 0.25))
	C.Material.RefractiveIndex = 2.5

	w := NewEmptyWorld()
	w.AddObject(A)
	w.AddObject(B)
	w.AddObject(C)

	r := rays.NewRay(coordinates.CreatePoint(0, 0, -4), coordinates.CreateVector(0, 0, 1))
	xs := w.IntersectWithRay(r)

	expected_n1 := []float64{1.0, 1.5, 2.0, 2.5, 2.5, 1.5}
	expected_n2 := []float64{1.5, 2.0, 2.5, 2.5, 1.5, 1.0}
	for i, intersection := range xs {
		pre := PreparePrecompData(intersection, r, xs)
		helpers.ApproxEqual(t, expected_n1[i], pre.RI_Inbound, 0.0001)
		helpers.ApproxEqual(t, expected_n2[i], pre.RI_Outbound, 0.0001)
	}
}

func TestPrecompTotalInternalReflection(t *testing.T) {
	s := rays.NewGlassSphere()
	s.SetTransformation(matrices.ScalingMatrix(1, 1, 1))
	s.Material.RefractiveIndex = 1.5

	w := NewEmptyWorld()
	w.AddObject(s)

	r := rays.NewRay(coordinates.CreatePoint(0, 0, math.Sqrt(2)/2), coordinates.CreateVector(0, 1, 0))
	xs := []rays.Intersection{rays.NewIntersection(-math.Sqrt(2)/2, s), rays.NewIntersection(math.Sqrt(2)/2, s)}

	pre := PreparePrecompData(xs[1], r, xs)
	c := pre.Refracted_Colour(w, 2)
	assert.Equal(t, rays.Colour{0, 0, 0}, c)
}
