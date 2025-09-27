package observe

import (
	"rattata/coordinates"
	"rattata/helpers"
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
