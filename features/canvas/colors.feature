Feature: Perform color ops

    Scenario: Colors are (red, green, blue) tuples 
        Given c ← color(255, 45, 184)
        Then c.red = 255
            And c.green = 45
            And c.blue = 184