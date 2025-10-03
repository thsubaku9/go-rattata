package rays

import (
	"math"
	"rattata/coordinates"
	"rattata/matrices"
)

type Pattern interface {
	PatternAt(point coordinates.Coordinate) Colour
	PatternTransformation() matrices.Matrix
}

// ------------------------------------ No Pattern ------------------------------------
type PlainPattern struct {
	colorMain Colour
}

func NewPlainPattern(col Colour) PlainPattern {
	return PlainPattern{col}
}

func (p PlainPattern) PatternAt(point coordinates.Coordinate) Colour {
	return p.colorMain
}

func (p PlainPattern) PatternTransformation() matrices.Matrix {
	return matrices.NewIdentityMatrix(4)
}

// ------------------------------------ Stripe Pattern ------------------------------------
type XStripe struct {
	colourA         Colour
	colourB         Colour
	transformMatrix matrices.Matrix
}

func NewXStripe(colA, colB Colour) XStripe {
	return XStripe{colA, colB, matrices.NewIdentityMatrix(4)}
}

func (stripe XStripe) PatternAt(point coordinates.Coordinate) Colour {

	patternTransformationInverse, _ := stripe.transformMatrix.Inverse()
	pattern_point := matrices.MatrixToCoordinate(matrices.PerformOrderedChainingOps(matrices.CoordinateToMatrix(point), patternTransformationInverse))

	if int(math.Floor(pattern_point.Get(coordinates.X)))%2 == 0 {
		return stripe.colourA
	} else {
		return stripe.colourB
	}
}

func (stripe XStripe) PatternTransformation() matrices.Matrix {
	return stripe.transformMatrix
}

func (stripe *XStripe) SetPatternTransformation(_mat matrices.Matrix) {
	stripe.transformMatrix = _mat
}

// ------------------------------------ Gradient Pattern ------------------------------------

type XGradient struct {
	colourA         Colour
	colourB         Colour
	transformMatrix matrices.Matrix
}

func NewXGradient(colA, colB Colour) XGradient {
	return XGradient{colA, colB, matrices.NewIdentityMatrix(4)}
}

func (grad XGradient) PatternAt(point coordinates.Coordinate) Colour {

	patternTransformationInverse, _ := grad.transformMatrix.Inverse()
	pattern_point := matrices.MatrixToCoordinate(matrices.PerformOrderedChainingOps(matrices.CoordinateToMatrix(point), patternTransformationInverse))

	distance := SubColour(grad.colourB, grad.colourA)
	fraction := pattern_point.Get(coordinates.X) - math.Floor(pattern_point.Get(coordinates.X))
	res := AddColour(grad.colourA, MulColour(distance, fraction))
	return res
}

func (grad XGradient) PatternTransformation() matrices.Matrix {
	return grad.transformMatrix
}

func (grad *XGradient) SetPatternTransformation(_mat matrices.Matrix) {
	grad.transformMatrix = _mat
}

// ------------------------------------ Ring Pattern ------------------------------------

type XZRing struct {
	colourA         Colour
	colourB         Colour
	transformMatrix matrices.Matrix
}

func NewXZRing(colA, colB Colour) XZRing {
	return XZRing{colA, colB, matrices.NewIdentityMatrix(4)}
}

func (r XZRing) PatternAt(point coordinates.Coordinate) Colour {
	patternTransformationInverse, _ := r.transformMatrix.Inverse()
	pattern_point := matrices.MatrixToCoordinate(matrices.PerformOrderedChainingOps(matrices.CoordinateToMatrix(point), patternTransformationInverse))

	if int(math.Floor(math.Sqrt(pattern_point.Get(coordinates.X)*pattern_point.Get(coordinates.X)+pattern_point.Get(coordinates.Z)*pattern_point.Get(coordinates.Z))))%2 == 0 {
		return r.colourA
	}

	return r.colourB

}

func (r XZRing) PatternTransformation() matrices.Matrix {
	return r.transformMatrix
}

func (r *XZRing) SetPatternTransformation(_mat matrices.Matrix) {
	r.transformMatrix = _mat
}

// ------------------------------------ Checker Pattern ------------------------------------

type Checker3D struct {
	colourA         Colour
	colourB         Colour
	transformMatrix matrices.Matrix
}

func NewChecker3D(colA, colB Colour) Checker3D {
	return Checker3D{colA, colB, matrices.NewIdentityMatrix(4)}
}

func (chk Checker3D) PatternAt(point coordinates.Coordinate) Colour {
	patternTransformationInverse, _ := chk.transformMatrix.Inverse()
	pattern_point := matrices.MatrixToCoordinate(matrices.PerformOrderedChainingOps(matrices.CoordinateToMatrix(point), patternTransformationInverse))

	if (int(math.Floor(pattern_point.Get(coordinates.X))+math.Floor(pattern_point.Get(coordinates.Y))+math.Floor(pattern_point.Get(coordinates.Z))) % 2) == 0 {
		return chk.colourA
	}

	return chk.colourB

}

func (chk Checker3D) PatternTransformation() matrices.Matrix {
	return chk.transformMatrix
}

func (chk *Checker3D) SetPatternTransformation(_mat matrices.Matrix) {
	chk.transformMatrix = _mat
}

// ------------------------------------ UV Checker Pattern ------------------------------------

type UnitSphereUVChecker struct {
	colourA         Colour
	colourB         Colour
	transformMatrix matrices.Matrix
	width           float64
	height          float64
}

func NewUnitSphereUVChecker(colA, colB Colour, width, height float64) UnitSphereUVChecker {
	return UnitSphereUVChecker{colA, colB, matrices.NewIdentityMatrix(4), width, height}
}

func (chk UnitSphereUVChecker) PatternAt(point coordinates.Coordinate) Colour {
	patternTransformationInverse, _ := chk.transformMatrix.Inverse()
	pattern_point := matrices.MatrixToCoordinate(matrices.PerformOrderedChainingOps(matrices.CoordinateToMatrix(point), patternTransformationInverse))

	u := 0.5 + math.Atan2(pattern_point.Get(coordinates.Z), pattern_point.Get(coordinates.X))/(2*math.Pi)
	v := 0.5 + math.Asin(pattern_point.Get(coordinates.Y))/math.Pi

	if (int(math.Floor(u*chk.width)+math.Floor(v*chk.height)) % 2) == 0 {
		return chk.colourA
	}

	return chk.colourB
}

func (chk UnitSphereUVChecker) PatternTransformation() matrices.Matrix {
	return chk.transformMatrix
}

func (chk *UnitSphereUVChecker) SetPatternTransformation(_mat matrices.Matrix) {
	chk.transformMatrix = _mat
}

// ------------------------------------ Radial Gradient Pattern ------------------------------------

type XZRadialGradient struct {
	colourA         Colour
	colourB         Colour
	transformMatrix matrices.Matrix
}

func NewXZRadialGradient(colA, colB Colour) XZRadialGradient {
	return XZRadialGradient{colA, colB, matrices.NewIdentityMatrix(4)}
}

func (rg XZRadialGradient) PatternAt(point coordinates.Coordinate) Colour {
	patternTransformationInverse, _ := rg.transformMatrix.Inverse()
	pattern_point := matrices.MatrixToCoordinate(matrices.PerformOrderedChainingOps(matrices.CoordinateToMatrix(point), patternTransformationInverse))

	distance := SubColour(rg.colourB, rg.colourA)
	fraction := math.Sqrt(pattern_point.Get(coordinates.X)*pattern_point.Get(coordinates.X) + pattern_point.Get(coordinates.Z)*pattern_point.Get(coordinates.Z))
	res := AddColour(rg.colourA, MulColour(distance, fraction))
	return res
}

func (rg XZRadialGradient) PatternTransformation() matrices.Matrix {
	return rg.transformMatrix
}

func (rg *XZRadialGradient) SetPatternTransformation(_mat matrices.Matrix) {
	rg.transformMatrix = _mat
}

// ------------------------------------ Perturbed Pattern ------------------------------------

type Perturbed struct {
	basePattern     Pattern
	perturbAmount   float64
	transformMatrix matrices.Matrix
}

func NewPerturbedPattern(base Pattern, perturbAmt float64) Perturbed {
	return Perturbed{base, perturbAmt, matrices.NewIdentityMatrix(4)}
}

func (p Perturbed) PatternAt(point coordinates.Coordinate) Colour {
	patternTransformationInverse, _ := p.transformMatrix.Inverse()
	pattern_point := matrices.MatrixToCoordinate(matrices.PerformOrderedChainingOps(matrices.CoordinateToMatrix(point), patternTransformationInverse))

	x2 := pattern_point.Get(coordinates.X) + (PerlinNoise3D(pattern_point.Get(coordinates.X), pattern_point.Get(coordinates.Y), pattern_point.Get(coordinates.Z)) * p.perturbAmount)
	y2 := pattern_point.Get(coordinates.Y) + (PerlinNoise3D(pattern_point.Get(coordinates.Y), pattern_point.Get(coordinates.Z), pattern_point.Get(coordinates.X)) * p.perturbAmount)
	z2 := pattern_point.Get(coordinates.Z) + (PerlinNoise3D(pattern_point.Get(coordinates.Z), pattern_point.Get(coordinates.X), pattern_point.Get(coordinates.Y)) * p.perturbAmount)

	return p.basePattern.PatternAt(coordinates.CreatePoint(x2, y2, z2))
}

func (p Perturbed) PatternTransformation() matrices.Matrix {
	return p.transformMatrix
}

func (p *Perturbed) SetPatternTransformation(_mat matrices.Matrix) {
	p.transformMatrix = _mat
}

func PerlinNoise3D(x, y, z float64) float64 {
	// Scale the input coordinates to control the "frequency" of the noise
	scale := 2.0
	x *= scale
	y *= scale
	z *= scale

	// Determine the grid cell coordinates
	x0 := int(math.Floor(x))
	x1 := x0 + 1
	y0 := int(math.Floor(y))
	y1 := y0 + 1
	z0 := int(math.Floor(z))
	z1 := z0 + 1

	// Relative coordinates within the grid cell
	xf := x - float64(x0)
	yf := y - float64(y0)
	zf := z - float64(z0)

	// Fade curves for each coordinate
	u := fade(xf)
	v := fade(yf)
	w := fade(zf)

	// Hash coordinates of the cube corners
	aaa := hash(x0, y0, z0)
	aba := hash(x0, y1, z0)
	aab := hash(x0, y0, z1)
	abb := hash(x0, y1, z1)
	baa := hash(x1, y0, z0)
	bba := hash(x1, y1, z0)
	bab := hash(x1, y0, z1)
	bbb := hash(x1, y1, z1)

	// And add blended results from 8 corners of the cube
	x1Interp := lerp(grad(aaa, xf, yf, zf), grad(baa, xf-1, yf, zf), u)
	x2Interp := lerp(grad(aba, xf, yf-1, zf), grad(bba, xf-1, yf-1, zf), u)
	y1Interp := lerp(x1Interp, x2Interp, v)

	x3Interp := lerp(grad(aab, xf, yf, zf-1), grad(bab, xf-1, yf, zf-1), u)
	x4Interp := lerp(grad(abb, xf, yf-1, zf-1), grad(bbb, xf-1, yf-1, zf-1), u)
	y2Interp := lerp(x3Interp, x4Interp, v)

	return (lerp(y1Interp, y2Interp, w) + 1) / 2 // Normalize to [0,1]
}

func fade(t float64) float64 {
	return t * t * t * (t*(t*6-15) + 10)
}

func lerp(a, b, t float64) float64 {
	return a + t*(b-a)
}

func grad(hash int, x, y, z float64) float64 {
	h := hash & 15
	u := ifThenElse(h < 8, x, y)
	v := ifThenElse(h < 4, y, ifThenElse(h == 12 || h == 14, x, z))
	return ifThenElse((h&1) == 0, u, -u) + ifThenElse((h&2) == 0, v, -v)
}

func hash(x, y, z int) int {
	// A simple hash function for demonstration purposes
	return (x * 73856093) ^ (y * 19349663) ^ (z * 83492791)
}

func ifThenElse(condition bool, a, b float64) float64 {
	if condition {
		return a
	}
	return b
}
