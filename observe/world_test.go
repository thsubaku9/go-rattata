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
	c := w.Color_At(r, 1)

	assert.Equal(t, rays.Colour{0, 0, 0}, c)
}

func TestWorldColorRayHit(t *testing.T) {
	w := NewDefaultWorld()
	r := rays.NewRay(coordinates.CreatePoint(0, 0, -5), coordinates.CreateVector(0, 0, 1))
	c := w.Color_At(r, 1)

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
	c := w.Color_At(r, 1)
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
	c := w.Color_At(r, 2)
	expected_c := rays.Colour{0.87677, 0.92436, 0.82918}

	for i := range expected_c {
		helpers.ApproxEqual(t, expected_c[i], c[i], 0.0001)
	}
}

func TestShadeHitWithRecursiveDepth(t *testing.T) {
	w := NewDefaultWorld()
	plane := rays.NewPlane(coordinates.CreatePoint(0, -1, 0))
	plane.Material.Reflective = 0.5
	plane.SetTransformation(matrices.TranslationMatrix(0, -1, 0))
	w.AddObject(plane)

	r := rays.NewRay(coordinates.CreatePoint(0, 0, -3), coordinates.CreateVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	c := w.Color_At(r, 5)
	expected_c := rays.Colour{0.87677, 0.92436, 0.82918}

	for i := range expected_c {
		helpers.ApproxEqual(t, expected_c[i], c[i], 0.0001)
	}

	c = w.Color_At(r, 0)
	expected_c = rays.Colour{0.0, 0.0, 0.0}

	for i := range expected_c {
		helpers.ApproxEqual(t, expected_c[i], c[i], 0.0001)
	}

}

func TestShadeHitWithRefractiveObject(t *testing.T) {
	w := NewDefaultWorld()
	plane := rays.NewPlane(coordinates.CreatePoint(0, 0, 0))
	plane.Material.Transparency = 0.5
	plane.Material.RefractiveIndex = 1.5
	plane.SetTransformation(matrices.TranslationMatrix(0, -1, 0))
	w.AddObject(plane)

	sph := rays.NewSphere(coordinates.CreatePoint(0, 0, 0), 1)
	sph.Material.Pattern = rays.NewPlainPattern(rays.Colour{1.0, 0.0, 0.0})
	sph.Material.Ambient = 0.5
	sph.SetTransformation(matrices.TranslationMatrix(0, -3.5, -0.5))

	w.AddObject(sph)

	r := rays.NewRay(coordinates.CreatePoint(0, 0, -3), coordinates.CreateVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))
	c := w.Color_At(r, 2)
	expected_c := rays.Colour{0.936422, 0.68642, 0.68642}

	for i := range expected_c {
		helpers.ApproxEqual(t, expected_c[i], c[i], 0.0001)
	}
}

func TestShadeHitWithReflectiveRefractiveObject(t *testing.T) {
	w := NewDefaultWorld()
	plane := rays.NewPlane(coordinates.CreatePoint(0, 0, 0))
	plane.Material.Reflective = 0.5
	plane.Material.Transparency = 0.5
	plane.Material.RefractiveIndex = 1.5
	plane.SetTransformation(matrices.TranslationMatrix(0, -1, 0))
	w.AddObject(plane)

	sph := rays.NewSphere(coordinates.CreatePoint(0, 0, 0), 1)
	sph.Material.Pattern = rays.NewPlainPattern(rays.Colour{1.0, 0.0, 0.0})
	sph.Material.Ambient = 0.5
	sph.SetTransformation(matrices.TranslationMatrix(0, -3.5, -0.5))

	w.AddObject(sph)

	r := rays.NewRay(coordinates.CreatePoint(0, 0, -3), coordinates.CreateVector(0, -math.Sqrt(2)/2, math.Sqrt(2)/2))

	c := w.Color_At(r, 3)
	expected_c := rays.Colour{1.296087, 0.69643, 0.69243}

	for i := range expected_c {
		helpers.ApproxEqual(t, expected_c[i], c[i], 0.0001)
	}
}
