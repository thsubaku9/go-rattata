package matrices

import "rattata/coordinates"

func CoordinateToMatrix(c coordinates.Coordinate) Matrix {
	_matrix := NewMatrix(4, 1)
	_matrix.Set(0, 0, c.Get(coordinates.X))
	_matrix.Set(1, 0, c.Get(coordinates.Y))
	_matrix.Set(2, 0, c.Get(coordinates.Z))
	_matrix.Set(3, 0, c.Get(coordinates.W))
	return _matrix
}

func MatrixToCoordinate(mtx Matrix) coordinates.Coordinate {
	_c := coordinates.CreateCoordinate(mtx.Get(0, 0), mtx.Get(1, 0), mtx.Get(2, 0), mtx.Get(3, 0))

	return _c
}

func TranslationMatrix(x, y, z float32) Matrix {
	_matrix := NewIdentityMatrix(4)
	_matrix.Set(0, 3, x)
	_matrix.Set(1, 3, y)
	_matrix.Set(2, 3, z)

	return _matrix
}

func ScalingMatrix(x, y, z float32) Matrix {
	_matrix := NewIdentityMatrix(4)
	_matrix.Set(0, 0, x)
	_matrix.Set(1, 1, y)
	_matrix.Set(2, 2, z)

	return _matrix
}

func PeformMatrixTranslation(src Matrix, x, y, z float32) Matrix {
	_matrix := TranslationMatrix(x, y, z)
	_, res := _matrix.Multiply(src)

	return res
}

func PeformMatrixScaling(src Matrix, x, y, z float32) Matrix {
	_matrix := ScalingMatrix(x, y, z)
	_, res := _matrix.Multiply(src)

	return res
}

func PeformMatrixRotation() {

}

func PeformMatrixShearing() {

}
