package main

import "rattata/playground"

func main() {
	// c := canvas.CreateCanvas(5, 5)

	// w, h := c.GetWidth(), c.GetHeight()

	// print(w, h)
	// print(len(c))
	// print(len(c[0]))

	// print(canvas.CanvasToPPMData(c))
	// for i := 0; i < w; i++ {
	// 	for j := 0; j < h; j++ {
	// 		c.WritePixel(uint32(i), uint32(j), canvas.Colour{uint8(rand.Intn(256)), uint8(rand.Intn(256)), uint8(rand.Intn(256)), 255})
	// 	}
	// }

	// print(canvas.CanvasToPPMData(c))

	// playground.ProcessRaySphereIntersection()

	// playground.ProcessPhongReflection()

	playground.PerformWorldBuildingDefault()

}
