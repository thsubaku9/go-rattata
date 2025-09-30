package matrices

import (
	"rattata/coordinates"
)

func View_Transform(from, to, up coordinates.Coordinate) Matrix {
	forward_vector := to.Sub(&from).Norm()
	up_norm := up.Norm()
	left_vector := forward_vector.CrossP(up_norm)
	true_up := left_vector.CrossP(forward_vector)

	orient_matrix := NewIdentityMatrix(4)

	orient_matrix.Set(0, 0, left_vector.Get(coordinates.X))
	orient_matrix.Set(0, 1, left_vector.Get(coordinates.Y))
	orient_matrix.Set(0, 2, left_vector.Get(coordinates.Z))

	orient_matrix.Set(1, 0, true_up.Get(coordinates.X))
	orient_matrix.Set(1, 1, true_up.Get(coordinates.Y))
	orient_matrix.Set(1, 2, true_up.Get(coordinates.Z))

	orient_matrix.Set(2, 0, -forward_vector.Get(coordinates.X))
	orient_matrix.Set(2, 1, -forward_vector.Get(coordinates.Y))
	orient_matrix.Set(2, 2, -forward_vector.Get(coordinates.Z))

	res, _ := orient_matrix.Multiply(TranslationMatrix(-from.Get(coordinates.X), -from.Get(coordinates.Y), -from.Get(coordinates.Z)))
	return res
}
