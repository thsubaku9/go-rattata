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
	return AddColour(grad.colourA, MulColour(distance, fraction))
}

func (grad XGradient) PatternTransformation() matrices.Matrix {
	return grad.transformMatrix
}

func (grad *XGradient) SetPatternTransformation(_mat matrices.Matrix) {
	grad.transformMatrix = _mat
}
