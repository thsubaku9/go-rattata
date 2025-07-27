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
