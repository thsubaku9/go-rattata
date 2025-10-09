Feature: Canvas related ops
    Scenario: Creating a canvas 
        Given c ← canvas(10, 20) 
        Then c.width = 10
            And c.height = 20
            And every pixel of c is color zero

    Scenario: Writing pixels to a canvas 
        Given c ← canvas(10, 20)
            And redPixel ← color(255, 0, 0) 
        When write_pixel(c, 2, 3, redPixel) 
        Then pixel_at(c, 2, 3) = redPixel

    Scenario: Constructing the PPM file properly 
        Given c ← canvas(5, 3)
        When insert random data of size 15
            And ppm ← canvas_to_ppm(c)
        Then header of ppm are
            """
            P3
            5 3
            255
            """
            And each line should try to not exceed 70 chars 
            And ppm ends with a newline character