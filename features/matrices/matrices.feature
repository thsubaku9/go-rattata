Feature: Matrix features
    Scenario: Constructing and inspecting a 4x4 matrix 
        Given the following 4x4 matrix M:
        | 1   | 2  |  3  |4   | 
        | 5.5 | 6.5| 7.5 | 8.5| 
        | 9   | 10 | 11  | 12 | 
        | 13.5|14.5| 15.5|16.5|
        Then M[0,0] = 1
            And M[0,3] = 4 
            And M[1,0] = 5.5 
            And M[1,2] = 7.5 
            And M[2,2] = 11 
            And M[3,0] = 13.5 
            And M[3,2] = 15.5

    Scenario: Multiplying two matrices 
        Given the following matrix A:
        |1|2|3|4| 
        |5|6|7|8| 
        |9|8|7|6| 
        |5|4|3|2|
            And the following matrix B: 
            |-2|1|2| 3|
            | 3|2|1|-1|
            | 4|3|6| 5|
            | 1|2|7| 8|
        Then A * B is the following 4x4 matrix:
        |20| 22| 50| 48| 
        |44| 54|114|108| 
        |40| 58|110|102| 
        |16| 26| 46| 42|
