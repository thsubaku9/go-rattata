package coordinates

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test2IntersectionsWithSphereFromOutside(t *testing.T) {
	r := NewRay(CreatePoint(0, 0, -5), CreateVector(0, 0, 1))
	sph := NewSphere(CreatePoint(0, 0, 0), 1)
	xs := Intersect(sph, r)

	assert.Equal(t, 2, len(xs))
	assert.Equal(t, float32(4.0), xs[0].Tvalue)
	assert.Equal(t, float32(6.0), xs[1].Tvalue)
}

func Test2IntersectionsWithSphereFromInside(t *testing.T) {
	r := NewRay(CreatePoint(0, 0, 0), CreateVector(0, 0, 1))
	sph := NewSphere(CreatePoint(0, 0, 0), 1)
	xs := Intersect(sph, r)

	assert.Equal(t, 2, len(xs))
	assert.Equal(t, float32(-1.0), xs[0].Tvalue)
	assert.Equal(t, float32(1.0), xs[1].Tvalue)
}

func Test1IntersectionWithSphere(t *testing.T) {
	r := NewRay(CreatePoint(0, 1, -5), CreateVector(0, 0, 1))
	sph := NewSphere(CreatePoint(0, 0, 0), 1)
	xs := Intersect(sph, r)

	assert.Equal(t, 2, len(xs))
	assert.Equal(t, float32(5.0), xs[0].Tvalue)
	assert.Equal(t, float32(5.0), xs[1].Tvalue)

}

func Test0IntersectionWithSphere(t *testing.T) {

	r := NewRay(CreatePoint(0, 2, -5), CreateVector(0, 0, 1))
	sph := NewSphere(CreatePoint(0, 0, 0), 1)
	xs := Intersect(sph, r)

	assert.Equal(t, 0, len(xs))
}

func TestHit(t *testing.T) {
	s := NewSphere(CreatePoint(0, 0, 0), 1)
	i1 := NewIntersection(5, s)
	i2 := NewIntersection(7, s)
	i3 := NewIntersection(-3, s)
	i4 := NewIntersection(2, s)
	xs := Intersections(i1, i2, i3, i4)
	i, res := Hit(xs)
	assert.True(t, res)
	assert.Equal(t, i4, *i)

}
