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

func TestDefaultWorld(t *testing.T) {
	w := NewDefaultWorld()

	assert.NotNil(t, w.LightSource())
	assert.NotNil(t, w.ListObjects())
	assert.NotEmpty(t, w.ListObjects())
}

func TestWorldIntersection(t *testing.T) {
	w := NewDefaultWorld()
	r := rays.NewRay(coordinates.CreatePoint(0, 0, -5), coordinates.CreateVector(0, 0, 1))

	xs := w.IntersectWithRay(r)

	assert.Equal(t, 4, len(xs))

	assert.Equal(t, 4.0, xs[0].Tvalue)
	assert.Equal(t, 4.5, xs[1].Tvalue)
	assert.Equal(t, 5.5, xs[2].Tvalue)
	assert.Equal(t, 6.0, xs[3].Tvalue)
}

func TestWorldColorRayMiss(t *testing.T) {
	w := NewDefaultWorld()
	r := rays.NewRay(coordinates.CreatePoint(0, 0, -5), coordinates.CreateVector(0, 1, 0))
	c := w.Color_At(r)

	assert.Equal(t, rays.Colour{0, 0, 0}, c)
}

func TestWorldColorRayHit(t *testing.T) {
	w := NewDefaultWorld()
	r := rays.NewRay(coordinates.CreatePoint(0, 0, -5), coordinates.CreateVector(0, 0, 1))
	c := w.Color_At(r)

	expected_c := rays.Colour{0.38066, 0.47583, 0.2855}

	for i := range expected_c {
		helpers.ApproxEqual(t, expected_c[i], c[i], 0.0001)
	}

}

func TestWorldColorRayBehind(t *testing.T) {
	w := NewDefaultWorld()

	w.PerformObjectModifications(1, func(obj rays.Shape) rays.Shape {
		inner := (w.ListObjects()[1]).(rays.Sphere)
		inner.Material.Ambient = 1
		return inner
	})

	r := rays.NewRay(coordinates.CreatePoint(0, 0, 0.75), coordinates.CreateVector(0, 0, -1))
	c := w.Color_At(r)
	assert.Equal(t, rays.Colour{1, 1, 1}, c)
}

func TestNoShadowWhenNothingCollinear(t *testing.T) {
	w := NewDefaultWorld()
	p := coordinates.CreatePoint(0, 10, 0)
	in_shadow := w.IsShadowed(p)
	assert.False(t, in_shadow)
}

func TestShadowWhenObjectBetweenAndLight(t *testing.T) {
	w := NewDefaultWorld()
	p := coordinates.CreatePoint(10, -10, 10)
	in_shadow := w.IsShadowed(p)
	assert.True(t, in_shadow)
}

func TestNoShadowWhenObjectBehindLight(t *testing.T) {
	w := NewDefaultWorld()
	p := coordinates.CreatePoint(-20, 20, -20)
	in_shadow := w.IsShadowed(p)
	assert.False(t, in_shadow)
}

func TestNoShadowWhenObjectBehindPoint(t *testing.T) {
	w := NewDefaultWorld()
	p := coordinates.CreatePoint(-2, 2, -2)
	in_shadow := w.IsShadowed(p)
	assert.False(t, in_shadow)
}

func TestShadeHitWithReflectiveObject(t *testing.T) {
	w := NewDefaultWorld()
	plane := rays.NewPlane(coordinates.CreatePoint(0, -1, 0))
	plane.Material.Reflective = 0.5
	plane.SetTransformation(matrices.TranslationMatrix(0, -1, 0))
	w.AddObject(plane)

	r := rays.NewRay(coordinates.CreatePoint(0, 0, -3), coordinates.CreateVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	c := w.Color_At(r)
	expected_c := rays.Colour{0.87677, 0.92436, 0.82918}

	for i := range expected_c {
		helpers.ApproxEqual(t, expected_c[i], c[i], 0.0001)
	}
}

/*
Scenario: The reflected color at the maximum recursive depth Given w ← default_world()
And shape ← plane() with:
| material.reflective | 0.5 | | transform | translation(0, -1, 0) |
And shape is added to w
And r ← ray(point(0, 0, -3), vector(0, -√2/2, √2/2)) And i ← intersection(√2, shape)
When comps ← prepare_computations(i, r) And color ← reflected_color(w, comps, 0)
Then color = color(0, 0, 0)

*/
