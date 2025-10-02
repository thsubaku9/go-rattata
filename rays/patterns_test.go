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

// todo
