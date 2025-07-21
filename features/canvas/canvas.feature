Feature: Canvas related ops
    Scenario: Creating a canvas 
        Given c ← canvas(10, 20) 
        Then c.width = 10
            And c.height = 20
            And every pixel of c is color(0, 0, 0)

    Scenario: Writing pixels to a canvas 
        Given c ← canvas(10, 20)
            And redPixel ← color(255, 0, 0) 
        When write_pixel(c, 2, 3, redPixel) 
        Then pixel_at(c, 2, 3) = redPixel

    Scenario: Constructing the PPM header 
        Given c ← canvas(5, 3)
        When ppm ← canvas_to_ppm(c)
        Then lines 1-3 of ppm are 
            """
            P3 
            5 3 
            255
            """

    Scenario: PPM files are terminated by a newline character 
        Given c ← canvas(5, 3)
        When ppm ← canvas_to_ppm(c)
        Then ppm ends with a newline character