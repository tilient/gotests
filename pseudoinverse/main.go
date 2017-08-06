package main

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
	"math"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	test03()
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
	// try to solve X for A * X = B given A and B
	N := 128
	M := N
	a := randomMatrix(N, M)
	b := randomMatrix(N, 1)
	x := a.pInverse().mult(b)
	bp := a.mult(x)

	fmt.Println("b  ", b[:2])
	fmt.Println("bp ", bp[:2])
	fmt.Println(" diff:", b.dif(bp))
}

func test03() {
	// try to solve X for A * X = Y given A and Y
	n := 512
	m := 512

	y := make([]float64, n)
	for i := 0; i < n; i++ {
		y[i] = 2.0*rand.Float64() - 1.0
	}
	Y := mat.NewVecDense(n, y)

	A := mat.NewDense(n, m, nil)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			A.Set(i, j, 2.0*rand.Float64()-1.0)
		}
	}
	x := make([]float64, m)
	X := mat.NewVecDense(m, x)
	yp := make([]float64, n)
	Yp := mat.NewVecDense(n, yp)

	// Option 1

	fmt.Println("Solving A * X = Y")
	err := X.SolveVec(A, Y)
	if err != nil {
		fmt.Println("ERR -", err)
	}

	Yp.MulVec(A, X)
	fmt.Println("y: ", y[:3])
	fmt.Println("yp:", yp[:3])

	// Option 2

	fmt.Println("Calculating Ainv")
	Ainv := pInv(A)
	fmt.Println("Calculated Ainv")
	X.MulVec(Ainv, Y)
	Yp.MulVec(A, X)
	fmt.Println("y: ", y[:3])
	fmt.Println("yp:", yp[:3])
}

//#############################################################
// Pseudo Inverse
//  Ben-Israel and Cohen Iteration: Xn+1 = 2 * Xn - Xn * m * Xn
//#############################################################

func pInv(m *mat.Dense) *mat.Dense {
	N, M := m.Dims()
	res := mat.DenseCopyOf(m.T())
	res.Scale(0.02/float64(N), res)
	tmp2 := mat.NewDense(M, M, nil)
	tmp1 := mat.NewDense(M, N, nil)
	tmp3 := mat.NewDense(M, N, nil)
	for !isInvOf(res, m, 0.01) {
		//res = (2.0 * res) - (res * m * res)
		tmp1.Scale(2.0, res)
		tmp2.Mul(res, m)
		tmp3.Mul(tmp2, res)
		res.Sub(tmp1, tmp3)
	}
	return res
}

func isInvOf(m1, m2 *mat.Dense, tol float64) bool {
	N, _ := m1.Dims()
	_, M := m2.Dims()
	mid := mat.NewDense(N, M, nil)
	mid.Mul(m1, m2)
	for rix := 0; rix < N; rix++ {
		for cix := 0; cix < M; cix++ {
			v := mid.At(rix, cix)
			if rix == cix {
				v -= 1.0
			}
			v = math.Abs(v)
			if v > tol {
				return false
			}
		}
	}
	return true
}

//#############################################################

func (mat matrix) pInverse() matrix {
	N := mat.nrOfRows()
	M := mat.nrOfColumns()
	if N > M {
		panic("*** ERROR *** pInverse with nrOfRows > nrOfColumns")
	}
	res := mat.transpose().multFactor(0.02 / float64(N))
	for !res.isInverseOf(mat, 0.0001) {
		res = res.multFactor(2.0).min(res.mult(mat).mult(res))
	}
	return res
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

func (mat matrix) dif(mat2 matrix) float64 {
	diff := 0.0
	for rix := 0; rix < mat.nrOfRows(); rix++ {
		for cix := 0; cix < mat2.nrOfColumns(); cix++ {
			d := mat[rix][cix] - mat2[rix][cix]
			diff += d * d
		}
	}
	return math.Sqrt(diff)
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

func (vec vector) dif(vec2 vector) float64 {
	diff := 0.0
	for ix := range vec {
		d := vec[ix] - vec2[ix]
		diff += d * d
	}
	return math.Sqrt(diff)
}

func (vec vector) asMatrix() matrix {
	m := newMatrix(len(vec), 1)
	for ix := range vec {
		m[ix][0] = vec[ix]
	}
	return m
}

//#############################################################
