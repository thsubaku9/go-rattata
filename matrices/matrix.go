package matrices

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
