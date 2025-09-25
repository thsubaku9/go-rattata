package observe

import (
	"rattata/coordinates"
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

	assert.Equal(t, 3, xs[0].Tvalue)
	assert.Equal(t, 4.0, xs[1].Tvalue)
	assert.Equal(t, 6.0, xs[2].Tvalue)
	assert.Equal(t, 7, xs[3].Tvalue)
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

	assert.Equal(t, rays.Colour{0.38066, 0.47583, 0.2855}, c)
}

// todo -> fix from here onwards ?
func TestWorldColorRayBehind(t *testing.T) {
	w := NewDefaultWorld()
	outer := (*w.ListObjects()[0]).(rays.Sphere)
	inner := (*w.ListObjects()[1]).(rays.Sphere)
	outer.Material.Ambient = 1
	inner.Material.Ambient = 1

	assert.Equal(t, 1, inner.Material.Ambient)

	r := rays.NewRay(coordinates.CreatePoint(0, 0, 0.75), coordinates.CreateVector(0, 0, -1))
	c := w.Color_At(r)
	assert.Equal(t, rays.Colour{1, 1, 1}, c)
}
