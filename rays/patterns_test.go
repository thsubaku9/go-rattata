package rays

import (
	"rattata/coordinates"
	"rattata/matrices"
	"testing"
)

// ------------------------------------ No Pattern ------------------------------------
func TestPlainPatternVal(t *testing.T) {
	pp := NewPlainPattern(Colour{1, 0, 0})
	col := pp.PatternAt(coordinates.CreatePoint(0, 0, 0))
	if col != (Colour{1, 0, 0}) {
		t.Errorf("Expected plain pattern to always return the same colour")
	}
}

// ------------------------------------ Stripe Pattern ------------------------------------

func TestXStripePatternVal(t *testing.T) {
	stripe := NewXStripe(Colour{1, 0, 0}, Colour{0, 1, 0})
	col1 := stripe.PatternAt(coordinates.CreatePoint(0, 0, 0))
	if col1 != (Colour{1, 0, 0}) {
		t.Errorf("Expected stripe pattern to return colour A at x=0")
	}
	col2 := stripe.PatternAt(coordinates.CreatePoint(0.9, 0, 0))
	if col2 != (Colour{1, 0, 0}) {
		t.Errorf("Expected stripe pattern to return colour A at x=0.9")
	}
	col3 := stripe.PatternAt(coordinates.CreatePoint(1.0, 0, 0))
	if col3 != (Colour{0, 1, 0}) {
		t.Errorf("Expected stripe pattern to return colour B at x=1.0")
	}
	col4 := stripe.PatternAt(coordinates.CreatePoint(-0.1, 0, 0))
	if col4 != (Colour{0, 1, 0}) {
		t.Errorf("Expected stripe pattern to return colour B at x=-0.1")
	}
	col5 := stripe.PatternAt(coordinates.CreatePoint(-1.0, 0, 0))
	if col5 != (Colour{0, 1, 0}) {
		t.Errorf("Expected stripe pattern to return colour B at x=-1.0")
	}
	col6 := stripe.PatternAt(coordinates.CreatePoint(-1.1, 0, 0))
	if col6 != (Colour{1, 0, 0}) {
		t.Errorf("Expected stripe pattern to return colour A at x=-1.1")
	}
}

func TestXStripePatternWithTransform(t *testing.T) {
	stripe := NewXStripe(Colour{1, 0, 0}, Colour{0, 1, 0})
	stripe.SetPatternTransformation(matrices.ScalingMatrix(2, 2, 2))
	col1 := stripe.PatternAt(coordinates.CreatePoint(1.5, 0, 0))
	if col1 != (Colour{1, 0, 0}) {
		t.Errorf("Expected stripe pattern to return colour A at x=1.5 with pattern scaled by 2")
	}
	col2 := stripe.PatternAt(coordinates.CreatePoint(2.5, 0, 0))
	if col2 != (Colour{0, 1, 0}) {
		t.Errorf("Expected stripe pattern to return colour B at x=2.5 with pattern scaled by 2")
	}
}

// ------------------------------------ Gradient Pattern ------------------------------------

func TestXGradientLinearInterpolation(t *testing.T) {
	grad := NewXGradient(Colour{1, 1, 1}, Colour{0, 0, 0})
	col1 := grad.PatternAt(coordinates.CreatePoint(0, 0, 0))
	if col1 != grad.colourA {
		t.Errorf("Expected gradient pattern to return colourA at x=0, got %v", col1)
	}
	col2 := grad.PatternAt(coordinates.CreatePoint(0.25, 0, 0))
	if col2 != (Colour{0.75, 0.75, 0.75}) {
		t.Errorf("Expected gradient pattern to return colour (0.75, 0.75, 0.75) at x=0.25, got %v", col2)
	}
	col3 := grad.PatternAt(coordinates.CreatePoint(0.5, 0, 0))
	if col3 != (Colour{0.5, 0.5, 0.5}) {
		t.Errorf("Expected	 gradient pattern to return colour (0.5, 0.5, 0.5) at x=0.5, got %v", col3)
	}
	col4 := grad.PatternAt(coordinates.CreatePoint(0.75, 0, 0))
	if col4 != (Colour{0.25, 0.25, 0.25}) {
		t.Errorf("Expected gradient pattern to return colour (0.25, 0.25, 0.25) at x=0.75, got %v", col4)
	}
}

// ------------------------------------ Ring Pattern ------------------------------------

func TestXZRingPatternVal(t *testing.T) {
	ring := NewXZRing(Colour{1, 0, 0}, Colour{0, 1, 0})
	col1 := ring.PatternAt(coordinates.CreatePoint(0, 0, 0))
	if col1 != (Colour{1, 0, 0}) {
		t.Errorf("Expected ring pattern to return colour A at origin")
	}
	col2 := ring.PatternAt(coordinates.CreatePoint(1, 0, 0))
	if col2 != (Colour{0, 1, 0}) {
		t.Errorf("Expected ring pattern to return colour B at x=1")
	}
	col3 := ring.PatternAt(coordinates.CreatePoint(0, 0, 1))
	if col3 != (Colour{0, 1, 0}) {
		t.Errorf("Expected ring pattern to return colour B at z=1")
	}
	col4 := ring.PatternAt(coordinates.CreatePoint(0.708, 0, 0.708))
	if col4 != (Colour{0, 1, 0}) {
		t.Errorf("Expected ring pattern to return colour B at x=0.708, z=0.708")
	}
}

// ------------------------------------ Checker Pattern ------------------------------------
func TestXZRing(t *testing.T) {
	ring := NewChecker3D(Colour{1, 0, 0}, Colour{0, 1, 0})
	col1 := ring.PatternAt(coordinates.CreatePoint(0, 0, 0))
	if col1 != (Colour{1, 0, 0}) {
		t.Errorf("Expected checker pattern to return colour A at origin")
	}
	col2 := ring.PatternAt(coordinates.CreatePoint(0.99, 0, 0.99))
	if col2 != (Colour{1, 0, 0}) {
		t.Errorf("Expected checker pattern to return colour A at x=0.99, z=0.99")
	}
	col3 := ring.PatternAt(coordinates.CreatePoint(1.01, 0, 0.99))
	if col3 != (Colour{0, 1, 0}) {
		t.Errorf("Expected checker pattern to return colour B at x=1.01, z=0.99")
	}
	col4 := ring.PatternAt(coordinates.CreatePoint(0.99, 0, 1.01))
	if col4 != (Colour{0, 1, 0}) {
		t.Errorf("Expected checker pattern to return colour B at x=0.99, z=1.01")
	}
	col5 := ring.PatternAt(coordinates.CreatePoint(1.01, 0, 1.01))
	if col5 != (Colour{1, 0, 0}) {
		t.Errorf("Expected checker pattern to return colour A at x=1.01, z=1.01")
	}
}

// ------------------------------------ UV Checker Pattern ------------------------------------

func TestUnitSphereUVChecker(t *testing.T) {
	chk := NewUnitSphereUVChecker(Colour{1, 1, 1}, Colour{0, 0, 0}, 1, 1)

	col1 := chk.PatternAt(coordinates.CreatePoint(0, 0, 0))
	if col1 != (chk.colourA) {
		t.Errorf("Expected UV checker pattern to return colourA at (0,0,0)")
	}

	col2 := chk.PatternAt(coordinates.CreatePoint(0, 1, 0))
	if col2 != chk.colourB {
		t.Errorf("Expected UV checker pattern to return colourB at (0,1,0)")
	}

	col3 := chk.PatternAt(coordinates.CreatePoint(0, -1, 0))
	if col3 != chk.colourA {
		t.Errorf("Expected UV checker pattern to return colourB at (0,-1,0)")
	}

	col4 := chk.PatternAt(coordinates.CreatePoint(-0.707, -0.707, 0))
	if col4 != chk.colourB {
		t.Errorf("Expected UV checker pattern to return colourB at (-0.707,-0.707,0)")
	}
}
