package matrices

import "errors"

type Matrix [][]float32

func (m Matrix) Row() int {
	return len(m)
}

func (m Matrix) Column() int {
	return len(m[0])
}

func (m Matrix) Get(r, c int) float32 {
	return m[r][c]
}

func (m Matrix) Set(r, c int, val float32) {
	m[r][c] = val
}

func (m Matrix) IsEqual(m2 Matrix) bool {
	if len(m) != len(m2) || len(m[0]) != len(m2[0]) {
		return false
	}

	r, c := m.Row(), m.Column()

	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			if m[i][j] != m2[i][j] {
				return false
			}
		}
	}

	return true
}

func (m Matrix) ScaleMul(k float32) Matrix {
	r, c := m.Row(), m.Column()

	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {

			m.Set(i, j, m.Get(i, j)*k)
		}
	}

	return m
}

func (m Matrix) ScaleAdd(k float32) Matrix {
	r, c := m.Row(), m.Column()

	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {

			m.Set(i, j, m.Get(i, j)+k)
		}
	}

	return m
}

func (m Matrix) Multiply(m2 Matrix) (bool, Matrix) {
	if m.Column() != m2.Row() {
		return false, nil
	}
	_matrix := NewMatrix(m.Row(), m2.Column())

	r, c, l := _matrix.Row(), _matrix.Column(), m.Column()

	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			v := float32(0.0)
			for k := 0; k < l; k++ {
				v += m[i][k] * m2[k][j]
			}
			_matrix.Set(i, j, v)
		}
	}

	return true, _matrix
}

func (m Matrix) T() Matrix {
	_matrix := NewMatrix(m.Column(), m.Row())

	r, c := _matrix.Row(), _matrix.Column()

	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			_matrix.Set(i, j, m.Get(j, i))
		}
	}

	return _matrix
}

func (m Matrix) SubMatrix(r_t, c_t int) Matrix {
	// The new matrix will have one less row and one less column
	_matrix := NewMatrix(m.Row()-1, m.Column()-1)

	// Initialize row and column counters for the new submatrix
	newRow := 0
	newCol := 0

	// Iterate through the original matrix
	for i := 0; i < m.Row(); i++ {
		if i == r_t {
			// Skip the target row
			continue
		}
		newCol = 0 // Reset column for each new row in the submatrix
		for j := 0; j < m.Column(); j++ {
			if j == c_t {
				// Skip the target column
				continue
			}
			// Set the value in the new submatrix
			_matrix.Set(newRow, newCol, m.Get(i, j))
			newCol++
		}
		newRow++
	}
	return _matrix
}

func (m Matrix) Determinant() (float32, error) {
	if m.Column() != m.Row() {
		return 0, errors.New("NA")
	}

	n := m.Row()

	switch n {
	case 0:
		return 1, nil
	case 1:
		return m[0][0], nil
	case 2:
		return m[0][0]*m[1][1] - m[0][1]*m[1][0], nil
	}

	res := float32(0.0)

	for i := range n {
		d, _ := m.SubMatrix(0, i).Determinant()
		row_val := m[0][i] * d

		row_val = fetchMatrixPositionalSign(0, i) * row_val
		res += row_val
	}

	return res, nil
}

func (m Matrix) Minor(i, j int) float32 {
	d, _ := m.SubMatrix(i, j).Determinant()
	return d
}

func (m Matrix) Cofactor(i, j int) float32 {
	d, _ := m.SubMatrix(i, j).Determinant()
	return d * fetchMatrixPositionalSign(i, j)
}

func fetchMatrixPositionalSign(row, col int) float32 {
	if (row+col)%2 == 0 {
		return 1
	}

	return -1
}

func (m Matrix) Adj() (Matrix, error) {
	if m.Row() != m.Column() {
		return nil, errors.New("pseudo inverse not supported currently")
	}

	n := m.Row()
	_matrix := NewMatrix(n, n)
	det, _ := m.Determinant()
	for i_row := range n {
		for j_col := range n {
			_matrix[i_row][j_col] = m.Cofactor(i_row, j_col)
		}
	}

	return _matrix.T().ScaleMul(1 / det), nil
}

func IsMatrixInvertableBasedOnDeterminant(val float32) bool {
	return val != 0
}

func NewMatrix(r, c int) Matrix {
	mt := make([][]float32, r, r)

	for i := 0; i < r; i++ {
		mt[i] = make([]float32, c, c)
	}

	return mt
}

func NewIdentityMatrix(ln int) Matrix {
	_matrix := NewMatrix(ln, ln)

	for i := 0; i < ln; i++ {
		_matrix[i][i] = 1
	}

	return _matrix
}
