package matrix

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
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
	const maxRow = 3
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

func (X Matrix) Solve(A, Y Matrix) {
	a := mat.NewDense(A.NrOfRows(), A.NrOfColumns(), nil)
	x := mat.NewDense(X.NrOfRows(), X.NrOfColumns(), nil)
	y := mat.NewDense(Y.NrOfRows(), Y.NrOfColumns(), nil)

	N, M := A.NrOfRows(), A.NrOfColumns()
	for rix := 0; rix < N; rix++ {
		for cix := 0; cix < M; cix++ {
			a.Set(rix, cix, A[rix][cix])
		}
	}
	N, M = X.NrOfRows(), X.NrOfColumns()
	for rix := 0; rix < N; rix++ {
		for cix := 0; cix < M; cix++ {
			x.Set(rix, cix, X[rix][cix])
		}
	}
	N, M = Y.NrOfRows(), Y.NrOfColumns()
	for rix := 0; rix < N; rix++ {
		for cix := 0; cix < M; cix++ {
			y.Set(rix, cix, Y[rix][cix])
		}
	}

	err := x.Solve(a, y)
	if err != nil {
		fmt.Println("ERR -", err)
	}

	N, M = A.NrOfRows(), A.NrOfColumns()
	for rix := 0; rix < N; rix++ {
		for cix := 0; cix < M; cix++ {
			A[rix][cix] = a.At(rix, cix)
		}
	}
	N, M = X.NrOfRows(), X.NrOfColumns()
	for rix := 0; rix < N; rix++ {
		for cix := 0; cix < M; cix++ {
			X[rix][cix] = x.At(rix, cix)
		}
	}
	N, M = Y.NrOfRows(), Y.NrOfColumns()
	for rix := 0; rix < N; rix++ {
		for cix := 0; cix < M; cix++ {
			Y[rix][cix] = y.At(rix, cix)
		}
	}
}

func (X Matrix) Solve2(A, B Matrix) {
	N, K := A.NrOfRows(), A.NrOfColumns()
	M := X.NrOfColumns()
	delta := 1000000000.0
	for delta > 1.7 {
		delta = 0.0
		for rix := 0; rix < N; rix++ {
			for cix := 0; cix < M; cix++ {
				if rix != 0 || cix != 0 {
					v := 0.0
					for k := 0; k < K; k++ {
						v += A[rix][k] * X[k][cix]
					}
					v -= B[rix][cix]
					delta += v * v
					v /= float64(K)
					for k := 0; k < K; k++ {
						v0 := X[k][cix]
						X[k][cix] = 0.95*v0 - 0.05*v
					}
				}
			}
		}
		delta = math.Sqrt(delta)
		fmt.Println("delta =", delta)
	}
}

func (mat Matrix) PseudoInverse() Matrix {
	N, M := mat.NrOfRows(), mat.NrOfColumns()
	inv := mat.Transpose()
	delta := 1000000000.0
	for delta > 4.0 {
		delta = 0.0
		for rix := 0; rix < N; rix++ {
			for cix := 0; cix < N; cix++ {
				v := 0.0
				for k := 0; k < M; k++ {
					v += mat[rix][k] * inv[k][cix]
				}
				if rix == cix {
					v -= 1.0
				}
				delta += v * v
				v /= float64(M)
				for k := 0; k < M; k++ {
					v0 := inv[k][cix]
					inv[k][cix] = 0.95*v0 + 0.05*v
				}
			}
		}
		delta = math.Sqrt(delta)
		fmt.Println("delta =", delta)
	}
	return inv
}

//------------------------------------------------------------
