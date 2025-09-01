package matrices

import (
	"math"
	"rattata/coordinates"
)

func CoordinateToMatrix(c coordinates.Coordinate) Matrix {
	_matrix := NewMatrix(4, 1)
	_matrix.Set(0, 0, c.Get(coordinates.X))
	_matrix.Set(1, 0, c.Get(coordinates.Y))
	_matrix.Set(2, 0, c.Get(coordinates.Z))
	_matrix.Set(3, 0, c.Get(coordinates.W))
	return _matrix
}

func MatrixToCoordinate(mtx Matrix) coordinates.Coordinate {
	return coordinates.CreateCoordinate(mtx.Get(0, 0), mtx.Get(1, 0), mtx.Get(2, 0), mtx.Get(3, 0))
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

func GivensRotationMatrix3D(rotatingAxis coordinates.CoordinateAxis, rad float64) Matrix {
	_matrix := NewIdentityMatrix(4)

	switch rotatingAxis {
	case coordinates.X:
		_matrix.Set(1, 1, float32(math.Cos(rad)))
		_matrix.Set(2, 2, float32(math.Cos(rad)))

		_matrix.Set(1, 2, -float32(math.Sin(rad)))
		_matrix.Set(2, 1, float32(math.Sin(rad)))
	case coordinates.Y:
		_matrix.Set(0, 0, float32(math.Cos(rad)))
		_matrix.Set(2, 2, float32(math.Cos(rad)))

		_matrix.Set(0, 2, -float32(math.Sin(rad)))
		_matrix.Set(2, 0, float32(math.Sin(rad)))
	case coordinates.Z:
		_matrix.Set(0, 0, float32(math.Cos(rad)))
		_matrix.Set(1, 1, float32(math.Cos(rad)))

		_matrix.Set(0, 1, -float32(math.Sin(rad)))
		_matrix.Set(1, 0, float32(math.Sin(rad)))

	}
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

func PerformMatrixRotation(src Matrix, rotatingAxis coordinates.CoordinateAxis, degreeRad float64) Matrix {
	_matrix := GivensRotationMatrix3D(rotatingAxis, degreeRad)

	_, res := _matrix.Multiply(src)

	return res
}

func ShearMatrix(xy, xz, yx, yz, zx, zy float32) Matrix {
	_matrix := NewIdentityMatrix(4)
	_matrix.Set(0, 1, xy)
	_matrix.Set(0, 2, xz)
	_matrix.Set(1, 0, yx)
	_matrix.Set(1, 2, yz)
	_matrix.Set(2, 0, zx)
	_matrix.Set(2, 1, zy)

	return _matrix
}

func PeformMatrixShearing(src Matrix, xy, xz, yx, yz, zx, zy float32) Matrix {
	_matrix := ShearMatrix(xy, xz, yx, yz, zx, zy)
	_, res := _matrix.Multiply(src)
	return res
}

func PerformOrderedChainingOps(src Matrix, opMatrix ...Matrix) Matrix {
	res := src

	for _, op := range opMatrix {
		_, res = op.Multiply(res)
	}
	return res
}
