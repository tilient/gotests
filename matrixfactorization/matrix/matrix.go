package matrix

import (
	"fmt"
	"math"

	. "github.com/tilient/gotests/matrixfactorization/vector"
)

//------------------------------------------------------------

type Matrix []Vector

func NewMatrix(rows, columns int) Matrix {

	mat := make(Matrix, rows)
	for ix := range mat {
		mat[ix] = NewVector(columns)
	}
	return mat
}

func RandomMatrix(rows, columns int) Matrix {
	mat := make(Matrix, rows)
	for ix := range mat {
		mat[ix] = RandomVector(columns)
	}
	return mat
}

//------------------------------------------------------------

func (mat Matrix) NrOfRows() int {
	return len(mat)
}

func (mat Matrix) NrOfColumns() int {
	return len(mat[0])
}

func (mat Matrix) Print() {
	const maxRow = 2
	for ix, row := range mat {
		row.Print()
		fmt.Println()
		if ix >= maxRow {
			fmt.Println(" ..")
			return
		}
	}
}

//------------------------------------------------------------

func (mat Matrix) Min(mat2 Matrix) Matrix {
	m := NewMatrix(mat.NrOfRows(), mat.NrOfColumns())
	for ix := range m {
		m[ix] = mat[ix].Min(mat2[ix])
	}
	return m
}

func (mat Matrix) Mult(mat2 Matrix) Matrix {
	m := NewMatrix(mat.NrOfRows(), mat2.NrOfColumns())
	maxK := mat.NrOfColumns()
	for rix, row := range m {
		for cix := range row {
			m[rix][cix] = 0.0
			for k := 0; k < maxK; k++ {
				m[rix][cix] += mat[rix][k] * mat2[k][cix]
			}
		}
	}
	return m
}

func (mat Matrix) Transpose() Matrix {
	tMat := NewMatrix(mat.NrOfColumns(), mat.NrOfRows())
	for rix, row := range mat {
		for cix, value := range row {
			tMat[cix][rix] = value
		}
	}
	return tMat
}

func (mat Matrix) Abs() Matrix {
	tMat := NewMatrix(mat.NrOfRows(), mat.NrOfColumns())
	for rix, row := range mat {
		for cix, value := range row {
			tMat[rix][cix] = math.Abs(value)
		}
	}
	return tMat
}

//------------------------------------------------------------
