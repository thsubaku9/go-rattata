package main

import (
	"rattata/canvas"
)

func main() {
	c := canvas.CreateCanvas(10, 20)

	w, h := c.GetWidth(), c.GetHeight()

	print(w, h)
	print(len(c))
	print(len(c[0]))

	// for i := 0; i < w; i++ {
	// 	for j := 0; j < h; j++ {

	// 		r, g, b :=
	// 			c[i][j].Colour.GetValue(canvas.Red),
	// 			c[i][j].Colour.GetValue(canvas.Green),
	// 			c[i][j].Colour.GetValue(canvas.Blue)

	// 		if r != 0 {
	// 			print("eh")
	// 		} else if g != 0 {
	// 			print("eh")
	// 		} else if b != 0 {
	// 			print("eh")
	// 		}
	// 	}
	// }

}
