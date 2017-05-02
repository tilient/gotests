package matrix

import (
	"fmt"
	"math/rand"

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

func RandomMatrix(rows, columns int, factor float64) Matrix {
	mat := make(Matrix, rows)
	for ix := range mat {
		mat[ix] = RandomVector(columns, factor)
	}
	return mat
}

func ArtificialMatrix(rows, columns int) Matrix {
	mat := RandomMatrix(rows, columns, 1.0)
	for ix, row := range mat {
		row[ix%columns] += 20.0
	}
	for ix := range mat {
		rix := rand.Intn(rows)
		mat[ix], mat[rix] = mat[rix], mat[ix]
	}
	for ix := range mat {
		mat[ix][0] = mat[ix][1] + mat[ix][2]
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
	for _, row := range mat {
		row.Print()
		fmt.Println()
	}
}

//------------------------------------------------------------

func (mat Matrix) min(mat2 Matrix) Matrix {
	m := NewMatrix(mat.NrOfRows(), mat.NrOfColumns())
	for ix := range m {
		m[ix] = mat[ix].Min(mat2[ix])
	}
	return m
}

func (mat Matrix) mult(mat2 Matrix) Matrix {
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

func (mat Matrix) transpose() Matrix {
	tMat := NewMatrix(mat.NrOfColumns(), mat.NrOfRows())
	for rix, row := range mat {
		for cix, value := range row {
			tMat[cix][rix] = value
		}
	}
	return tMat
}

//------------------------------------------------------------
