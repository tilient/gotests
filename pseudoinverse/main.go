package main

import (
	"fmt"
	"math"
	"math/rand"
	// "time"
)

func main() {
	// rand.Seed(time.Now().UTC().UnixNano())
	test02()
}

func test01fun(a, b, c, d float64) float64 {
	return a + 2*b + c + 3*d
}

func test01() {
	nrOfSamples := 56789
	a := randomMatrix(nrOfSamples, 4)
	y := newMatrix(nrOfSamples, 1)
	for i := 0; i < nrOfSamples; i++ {
		y[i][0] = test01fun(a[i][0], a[i][1], a[i][2], a[i][3])
	}
	x := a.pInverse().mult(y)
	fmt.Println("---------")
	for t := 0; t < 10; t++ {
		as := randomMatrix(1, 4)
		y := test01fun(as[0][0], as[0][1], as[0][2], as[0][3])
		yp := as.mult(x)
		//fmt.Println("y:", y, "yp:", yp[0][0])
		diff := math.Abs(y - yp[0][0])
		fmt.Println(" diff:", diff)
	}
	fmt.Println("---------")
}

func test02() {
	N := 8
	a := randomMatrix(N, N).multFactor(10.0)
	x := randomVector(N)
	b := randomVector(N)
	for t := 0; t < 18; t++ {
		x = sor(a, b, x, 0.01)
		bp := a.multVec(x)
		fmt.Println("")
		fmt.Println(b[:2])
		fmt.Println(bp[:2])
	}
}

//#############################################################
// Pseudo Inverse
//  Ben-Israel and Cohen Iteration: Xn+1 = 2 * Xn - Xn * m * Xn
//#############################################################

func (mat matrix) pInverse() matrix {
	N := mat.nrOfRows()
	M := mat.nrOfColumns()
	if N < M {
		panic("*** ERROR *** pInverse with nrOfRows > nrOfColumns")
	}
	res := mat.transpose().multFactor(0.02 / float64(N))
	for !res.isInverseOf(mat, 0.0001) {
		res = res.multFactor(2.0).min(res.mult(mat).mult(res))
	}
	return res
}

//#############################################################
// SOR
//#############################################################

func sor(a matrix, b vector, x vector, w float64) vector {
	n := len(x)
	x_new := newVector(n)
	for i := 0; i < n; i++ {
		x_new[i] = b[i]
		for j := 0; j < i; j++ {
			x_new[i] -= a[i][j] * x_new[j]
		}
		for j := i + 1; j < n; j++ {
			x_new[i] -= a[i][j] * x[j]
		}
		x_new[i] /= a[i][i]
	}
	for i := 0; i < n; i++ {
		x_new[i] = (1.0-w)*x[i] + w*x_new[i]
	}
	return x_new
}

//#############################################################
// Matrix
//#############################################################

type matrix []vector

func newMatrix(rows, columns int) matrix {
	mat := make(matrix, rows)
	for ix := range mat {
		mat[ix] = newVector(columns)
	}
	return mat
}

func randomMatrix(rows, columns int) matrix {
	mat := make(matrix, rows)
	for ix := range mat {
		mat[ix] = randomVector(columns)
	}
	return mat
}

func (mat matrix) print() {
	fmt.Println("[")
	for _, r := range mat {
		fmt.Println(" ", r)
	}
	fmt.Println("]")
}

func (mat matrix) nrOfRows() int {
	return len(mat)
}

func (mat matrix) nrOfColumns() int {
	return len(mat[0])
}

func (mat matrix) min(mat2 matrix) matrix {
	m := newMatrix(mat.nrOfRows(), mat.nrOfColumns())
	for ix := range m {
		m[ix] = mat[ix].min(mat2[ix])
	}
	return m
}

func (mat matrix) multFactor(f float64) matrix {
	m := newMatrix(mat.nrOfRows(), mat.nrOfColumns())
	for rix, row := range mat {
		for cix, v := range row {
			m[rix][cix] = f * v
		}
	}
	return m
}

func (mat matrix) multVec(v vector) vector {
	maxK := len(v)
	m := newVector(mat.nrOfRows())
	for rix := 0; rix < mat.nrOfRows(); rix++ {
		m[rix] = 0.0
		for k := 0; k < maxK; k++ {
			m[rix] += mat[rix][k] * v[k]
		}
	}
	return m
}

func (mat matrix) mult(mat2 matrix) matrix {
	m := newMatrix(mat.nrOfRows(), mat2.nrOfColumns())
	maxK := mat.nrOfColumns()
	for rix := 0; rix < mat.nrOfRows(); rix++ {
		for cix := 0; cix < mat2.nrOfColumns(); cix++ {
			m[rix][cix] = 0.0
			for k := 0; k < maxK; k++ {
				m[rix][cix] += mat[rix][k] * mat2[k][cix]
			}
		}
	}
	return m
}

func (mat matrix) transpose() matrix {
	tMat := newMatrix(mat.nrOfColumns(), mat.nrOfRows())
	for rix, row := range mat {
		for cix, value := range row {
			tMat[cix][rix] = value
		}
	}
	return tMat
}

func (mat matrix) isInverseOf(mat2 matrix, tol float64) bool {
	maxK := mat.nrOfColumns()
	for rix := 0; rix < mat.nrOfRows(); rix++ {
		for cix := 0; cix < mat2.nrOfColumns(); cix++ {
			v := 0.0
			for k := 0; k < maxK; k++ {
				v += mat[rix][k] * mat2[k][cix]
			}
			if rix == cix {
				v = math.Abs(v - 1.0)
			}
			if v > tol {
				return false
			}
		}
	}
	return true
}

//#############################################################
// Vector
//#############################################################

type vector []float64

func newVector(length int) vector {
	return make(vector, length)
}

func randomVector(length int) vector {
	vec := make(vector, length)
	for ix := range vec {
		vec[ix] = 2.0*rand.Float64() - 1.0
	}
	return vec
}

func (vec vector) min(vec2 vector) vector {
	v := newVector(len(vec))
	for ix := range v {
		v[ix] = vec[ix] - vec2[ix]
	}
	return v
}

//#############################################################
