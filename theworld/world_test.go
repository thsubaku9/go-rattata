package theworld

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

	assert.Equal(t, 3.585786437626905, xs[0].Tvalue)
	assert.Equal(t, 4.0, xs[1].Tvalue)
	assert.Equal(t, 6.0, xs[2].Tvalue)
	assert.Equal(t, 6.414213562373095, xs[3].Tvalue)
}
