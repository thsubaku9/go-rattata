package canvas

import (
	"fmt"
	"log"
	"os"
)

const PPM_LINE_LIM = 70

type Grid struct {
	Colour
}

type Canvas [][]Grid

func CreateCanvas(h, w uint32) Canvas {

	grids := make([][]Grid, w, w)

	for i := uint32(0); i < w; i++ {
		grids[i] = make([]Grid, h, h)

		for j := uint32(0); j < h; j++ {
			grids[i][j] = Grid{NewColour()}
		}
	}

	return Canvas(grids)
}

func (c Canvas) GetHeight() int {
	return len(c)
}

func (c Canvas) GetWidth() int {
	return len(c[0])
}

func (c Canvas) WritePixel(x, y uint32, col Colour) {
	c[x][y].Colour = col
}

/*
Notice how the first row of pixels comes first, then the second row, and so forth. Further, each row is terminated by a new line.
In addition, no line in a PPM file should be more than 70 characters long. Most image programs tend to accept PPM images with lines longer than that, but it’s a good idea to add new lines as needed to keep the lines shorter. (Just be careful to put the new line where a space would have gone, so you don’t split a number in half!
*/
func CanvasToPPMData(myCanvas Canvas) string {
	ppmData := ""
	ppmData = AppendWithLine(ppmData, "P3")
	ppmData = AppendWithLine(ppmData, fmt.Sprintf("%d %d", myCanvas.GetWidth(), myCanvas.GetHeight()))
	ppmData = AppendWithLine(ppmData, "255")

	line_data := ""
	cur_len := 0
	w, h := myCanvas.GetWidth(), myCanvas.GetHeight()

	for i := 0; i < w; i++ {
		for j := 0; j < h; j++ {
			r, g, b := fmt.Sprint(myCanvas[i][j].Colour.GetValue(Red)),
				fmt.Sprint(myCanvas[i][j].Colour.GetValue(Green)),
				fmt.Sprint(myCanvas[i][j].Colour.GetValue(Blue))

			if cur_len+len(r)+1 < PPM_LINE_LIM {
				line_data = Append(line_data, r, " ")
				cur_len += len(r) + 1
			} else {
				ppmData = AppendWithLine(ppmData, line_data)
				line_data = ""
				cur_len = 0
			}

			if cur_len+len(g)+1 < PPM_LINE_LIM {
				line_data = Append(line_data, g, " ")
				cur_len += len(r) + 1
			} else {
				ppmData = AppendWithLine(ppmData, line_data)
				line_data = ""
				cur_len = 0
			}

			if cur_len+len(b)+1 < PPM_LINE_LIM {
				line_data = Append(line_data, b, " ")
				cur_len += len(r) + 1
			} else {
				ppmData = AppendWithLine(ppmData, line_data)
				line_data = ""
				cur_len = 0
			}
		}
	}

	if cur_len > 0 {
		ppmData = AppendWithLine(ppmData, line_data)
		line_data = ""
		cur_len = 0
	}

	return ppmData
}

func AppendWithLine(src string, datum string) string {
	return Append(src, datum, "\n")
}

func Append(src, datum, suffix string) string {
	return src + datum + suffix
}

func SaveToPath(fileName, ppmData string) {
	file, err := os.Create(fileName + ".ppm")
	if err != nil {
		log.Default().Printf("File creation failed")
	}
	defer file.Close()

	_, err = file.WriteString(ppmData)

	if err != nil {
		log.Default().Printf("File write failed")
	}
}
