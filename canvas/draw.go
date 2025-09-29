package canvas

import (
	"fmt"
	"log"
	"os"
)

const PPM_LINE_LIM = 70

type Pixel struct {
	Colour
}

type Canvas [][]Pixel

// func CreateCanvas(w, h uint32) Canvas {

// 	grids := make([][]Pixel, w, w)

// 	for i := uint32(0); i < w; i++ {
// 		grids[i] = make([]Pixel, h, h)

// 		for j := uint32(0); j < h; j++ {
// 			grids[i][j] = Pixel{NewColour()}
// 		}
// 	}

// 	return Canvas(grids)
// }

func CreateCanvas(w, h uint32) Canvas {

	grids := make([][]Pixel, h, h)

	for i := uint32(0); i < h; i++ {
		grids[i] = make([]Pixel, w, w)

		for j := uint32(0); j < w; j++ {
			grids[i][j] = Pixel{NewColour()}
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

func (c Canvas) WritePixel(colIdx, rowIdx uint32, col Colour) {
	c[rowIdx][colIdx].Colour = col
}

func (c Canvas) ReadPixel(colIdx, rowIdx uint32) Pixel {
	return c[rowIdx][colIdx]
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

	for _rowIdx := 0; _rowIdx < h; _rowIdx++ {
		for _colIdx := 0; _colIdx < w; _colIdx++ {
			r, g, b := fmt.Sprint(myCanvas[_rowIdx][_colIdx].Colour.GetValue(Red)),
				fmt.Sprint(myCanvas[_rowIdx][_colIdx].Colour.GetValue(Green)),
				fmt.Sprint(myCanvas[_rowIdx][_colIdx].Colour.GetValue(Blue))

			{
				if cur_len+len(r)+1 >= PPM_LINE_LIM {
					ppmData = AppendWithLine(ppmData, line_data)
					line_data = ""
					cur_len = 0
				}

				line_data = Append(line_data, r, " ")
				cur_len += len(r) + 1
			}

			{
				if cur_len+len(g)+1 >= PPM_LINE_LIM {
					ppmData = AppendWithLine(ppmData, line_data)
					line_data = ""
					cur_len = 0
				}

				line_data = Append(line_data, g, " ")
				cur_len += len(b) + 1
			}

			{
				if cur_len+len(b)+1 >= PPM_LINE_LIM {
					ppmData = AppendWithLine(ppmData, line_data)
					line_data = ""
					cur_len = 0
				}

				line_data = Append(line_data, b, " ")
				cur_len += len(r) + 1
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
