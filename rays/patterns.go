package rays

import "rattata/coordinates"

type Pattern interface {
	PatternAt(point coordinates.Coordinate) Colour
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

// ------------------------------------ Stripe Pattern ------------------------------------
type XStripe struct {
	colourA Colour
	colourB Colour
}

func NewXStripe(colA, colB Colour) XStripe {
	return XStripe{colA, colB}
}

func (stripe XStripe) PatternAt(point coordinates.Coordinate) Colour {
	if int(point.Get(coordinates.X))%2 == 0 {
		return stripe.colourA
	} else {
		return stripe.colourB
	}
}
