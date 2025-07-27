package canvas

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/cucumber/godog"
	"github.com/stretchr/testify/assert"
)

const (
	colorHolderkey int = iota
	canvasHolderKey
	ppmDataKey
)

func givenAColour(ctx context.Context, r, g, b int) (context.Context, error) {
	return context.WithValue(ctx, colorHolderkey, Colour([4]uint8{uint8(r), uint8(g), uint8(b), 255})), nil
}

func checkRed(ctx context.Context, valExpected int) error {
	c, ok := ctx.Value(colorHolderkey).(Colour)

	if !ok {
		return errors.New("Value not found")
	}

	if !checkColourValueMatch(c, Red, uint8(valExpected)) {
		return errors.New("Value not matched")
	}
	return nil
}

func checkGreen(ctx context.Context, valExpected int) error {
	c, ok := ctx.Value(colorHolderkey).(Colour)

	if !ok {
		return errors.New("Value not found")
	}

	if !checkColourValueMatch(c, Green, uint8(valExpected)) {
		return errors.New("Value not matched")
	}
	return nil
}

func checkBlue(ctx context.Context, valExpected int) error {
	c, ok := ctx.Value(colorHolderkey).(Colour)

	if !ok {
		return errors.New("Value not found")
	}

	if !checkColourValueMatch(c, Blue, uint8(valExpected)) {
		return errors.New("Value not matched")
	}
	return nil
}

func checkColourValueMatch(c Colour, mode ColourMode, valExpected uint8) bool {
	return c.GetValue(mode) == valExpected
}

func createCanvas(ctx context.Context, arg1, arg2 int) (context.Context, error) {
	return context.WithValue(ctx, canvasHolderKey, CreateCanvas(uint32(arg1), uint32(arg2))), nil
}

func checkHeight(ctx context.Context, arg1 int) error {
	c, ok := ctx.Value(canvasHolderKey).(Canvas)

	if !ok {
		return errors.New("Value not found")
	}

	if c.GetHeight() != arg1 {
		return errors.New("Height mismatch")
	}

	return nil
}

func checkWidth(ctx context.Context, arg1 int) error {
	c, ok := ctx.Value(canvasHolderKey).(Canvas)

	if !ok {
		return errors.New("Value not found")
	}

	if c.GetWidth() != arg1 {
		return errors.New("Height mismatch")
	}

	return nil
}

func everyPixelOfCIsColor(ctx context.Context) error {
	c, ok := ctx.Value(canvasHolderKey).(Canvas)

	if !ok {
		return errors.New("Value not found")
	}

	w, h := c.GetWidth(), c.GetHeight()

	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			r, g, b :=
				c[i][j].Colour.GetValue(Red),
				c[i][j].Colour.GetValue(Green),
				c[i][j].Colour.GetValue(Blue)

			if r != 0 {
				return errors.New("Not 0 value")
			} else if g != 0 {
				return errors.New("Not 0 value")
			} else if b != 0 {
				return errors.New("Not 0 value")
			}
		}
	}
	return nil
}

func pixel_atcRedPixel(arg1, arg2 int) error {
	return godog.ErrPending
}

func redPixelColor(arg1, arg2, arg3 int) error {
	return godog.ErrPending
}

func write_pixelcRedPixel(arg1, arg2 int) error {
	return godog.ErrPending
}

func perform_canvas_to_ppm(ctx context.Context) (context.Context, error) {
	c, ok := ctx.Value(canvasHolderKey).(Canvas)

	if !ok {
		return ctx, errors.New("Value not found")
	}

	return context.WithValue(ctx, ppmDataKey, CanvasToPPMData(c)), nil
}

func eachLineShouldTryToNotExceedChars(arg1 int) error {
	return godog.ErrPending
}

func headerOfPpmAre(ctx context.Context, arg1 *godog.DocString) error {
	ppmData, ok := ctx.Value(ppmDataKey).(string)

	if !ok {
		return errors.New("Value not found")
	}
	splits := strings.Split(ppmData, "\n")

	testContentString := arg1.Content

	if testContentString != strings.Join(splits[0:3], "\n") {
		return errors.New("header did not match")
	}
	return nil
}

func insertRandomDataOfSize(arg1 int) error {
	return godog.ErrPending
}

func ppmEndsWithANewlineCharacter() error {
	return godog.ErrPending
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Given(`^c ← color\((\d+), (\d+), (\d+)\)$`, givenAColour)
	ctx.Then(`^c\.blue = (\d+)$`, checkBlue)
	ctx.Then(`^c\.green = (\d+)$`, checkGreen)
	ctx.Then(`^c\.red = (\d+)$`, checkRed)

	ctx.Given(`^c ← canvas\((\d+), (\d+)\)$`, createCanvas)
	ctx.When(`^ppm ← canvas_to_ppm\(c\)$`, perform_canvas_to_ppm)
	ctx.Then(`^c\.height = (\d+)$`, checkHeight)
	ctx.Then(`^c\.width = (\d+)$`, checkWidth)
	ctx.Then(`^every pixel of c is color zero$`, everyPixelOfCIsColor)

	ctx.Step(`^pixel_at\(c, (\d+), (\d+)\) = redPixel$`, pixel_atcRedPixel)
	ctx.Step(`^redPixel ← color\((\d+), (\d+), (\d+)\)$`, redPixelColor)
	ctx.Step(`^write_pixel\(c, (\d+), (\d+), redPixel\)$`, write_pixelcRedPixel)

	ctx.Step(`^each line should try to not exceed (\d+) chars$`, eachLineShouldTryToNotExceedChars)
	ctx.Step(`^header of ppm are$`, headerOfPpmAre)
	ctx.Step(`^insert random data of size (\d+)$`, insertRandomDataOfSize)
	ctx.Step(`^ppm ends with a newline character$`, ppmEndsWithANewlineCharacter)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../features/canvas"},
			TestingT: t, // Testing instance that will run subtests.
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func TestColourAdd(t *testing.T) {
	c1 := Colour{100, 200, 60, 255}
	c2 := Colour{100, 100, 0, 255}
	c3 := c1.Add(&c2)

	assert.Equal(t, Colour{200, 255, 60, 255}, *c3)
}

func TestColourSub(t *testing.T) {
	c1 := Colour{100, 200, 60, 255}
	c2 := Colour{100, 100, 0, 255}
	c3 := c1.Sub(&c2)

	assert.Equal(t, Colour{0, 100, 60, 255}, *c3)
}

func TestScalarMul(t *testing.T) {
	c1 := Colour{100, 200, 60, 120}
	c3 := c1.Mul(2)

	assert.Equal(t, Colour{200, 255, 120, 120}, *c3)
}

func TestColorMul(t *testing.T) {
	c1 := Colour{10, 20, 60, 255}
	c2 := Colour{10, 10, 0, 255}
	c3 := c1.Blend(&c2)

	assert.Equal(t, Colour{100, 200, 0, 255}, *c3)
}
