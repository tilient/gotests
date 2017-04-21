package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

//------------------------------------------------------------

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func newIntRange(n int) []int {
	r := make([]int, n)
	for i := 0; i < n; i++ {
		r[i] = i
	}
	return r
}

//------------------------------------------------------------

type Vector []float64

func newVector(length int) Vector {
	return make(Vector, length)
}

func randomVector(length int, factor float64) Vector {
	vec := make(Vector, length)
	for ix := range vec {
		vec[ix] = factor * rand.Float64()
	}
	return vec
}

func (vec Vector) min(vec2 Vector) Vector {
	v := newVector(len(vec))
	for ix := range v {
		v[ix] = vec[ix] - vec2[ix]
	}
	return v
}

func (vec Vector) print() {
	for _, v := range vec {
		fmt.Printf("%5.2f ", v)
	}
}

//------------------------------------------------------------

type Matrix []Vector

func newMatrix(rows, columns int) Matrix {

	mat := make(Matrix, rows)
	for ix := range mat {
		mat[ix] = newVector(columns)
	}
	return mat
}

func randomMatrix(rows, columns int, factor float64) Matrix {
	mat := make(Matrix, rows)
	for ix := range mat {
		mat[ix] = randomVector(columns, factor)
	}
	return mat
}

func artificialMatrix(rows, columns int) Matrix {
	mat := randomMatrix(rows, columns, 1.0)
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

func (mat Matrix) nrOfRows() int {
	return len(mat)
}

func (mat Matrix) nrOfColumns() int {
	return len(mat[0])
}

func (mat Matrix) min(mat2 Matrix) Matrix {
	m := newMatrix(mat.nrOfRows(), mat.nrOfColumns())
	for ix := range m {
		m[ix] = mat[ix].min(mat2[ix])
	}
	return m
}

func (mat Matrix) print() {
	for _, row := range mat {
		row.print()
		fmt.Println()
	}
}

func (mat Matrix) Transpose() Matrix {
	tMat := newMatrix(mat.nrOfColumns(), mat.nrOfRows())
	for rix, row := range mat {
		for cix, value := range row {
			tMat[cix][rix] = value
		}
	}
	return tMat
}

func (mat1 Matrix) Mult(mat2 Matrix) Matrix {
	mat := newMatrix(mat1.nrOfRows(), mat2.nrOfColumns())
	maxK := mat1.nrOfColumns()
	for rix, row := range mat {
		for cix := range row {
			mat[rix][cix] = 0.0
			for k := 0; k < maxK; k++ {
				mat[rix][cix] += mat1[rix][k] * mat2[k][cix]
			}
		}
	}
	return mat
}

//------------------------------------------------------------

func matrixMultErrorAt(R, P, Q Matrix, i, j int) float64 {
	K := Q.nrOfRows()
	e := R[i][j]
	for k := 0; k < K; k++ {
		e -= P[i][k] * Q[k][j]
	}
	return e
}

func factorizationError(R, P, Q Matrix) float64 {
	N := R.nrOfRows()
	M := R.nrOfColumns()
	K := Q.nrOfRows()
	e := 0.0
	for i := 0; i < N; i++ {
		for j := 0; j < M; j++ {
			pRij := 0.0
			for k := 0; k < K; k++ {
				Pik := P[i][k]
				Qkj := Q[k][j]
				//? e += (beta / 2.0) * (Pik*Pik + Qkj*Qkj)
				pRij += Pik * Qkj
			}
			eRij := R[i][j] - pRij
			e += eRij * eRij
		}
	}
	return e
}

//------------------------------------------------------------

type coordinate struct {
	i, j int
}

//------------------------------------------------------------

type coordinateGenerator struct {
	n, m   int
	randsI map[int]struct{}
	randsJ map[int]struct{}
}

func newCoordinateGenerator(N, M int) coordinateGenerator {
	maxNrOfIxs := minInt(N, M)
	result := coordinateGenerator{
		N, M,
		make(map[int]struct{}, maxNrOfIxs),
		make(map[int]struct{}, maxNrOfIxs)}
	return result
}

func (gen coordinateGenerator) new() coordinate {
	vi := rand.Intn(gen.n)
	vj := rand.Intn(gen.m)
	for true {
		if _, ok := gen.randsI[vi]; ok {
			vi = rand.Intn(gen.n)
			continue
		}
		if _, ok := gen.randsJ[vj]; ok {
			vj = rand.Intn(gen.m)
			continue
		}
		break
	}
	gen.randsI[vi] = struct{}{}
	gen.randsJ[vj] = struct{}{}
	return coordinate{vi, vj}
}

func (gen coordinateGenerator) release(c coordinate) {
	delete(gen.randsI, c.i)
	delete(gen.randsJ, c.j)
}

//------------------------------------------------------------

func matrixFactorization(R Matrix, K int) (Matrix, Matrix) {
	const alpha = 0.0002 // the learning rate
	const beta = 0.02    // the regularization parameter

	N := R.nrOfRows()
	M := R.nrOfColumns()
	maxSteps := 1000 * K * M * N
	fmt.Println("-----------------------------------")
	fmt.Println("R =", M, "x", N)
	fmt.Println("K =", K)
	fmt.Println("max steps =", maxSteps)

	nrOfWorkers := 2 * maxInt(1, runtime.NumCPU()-2)
	fmt.Println("nr of workers", nrOfWorkers)
	fmt.Println("-----------------------------------")
	workers := newIntRange(nrOfWorkers)
	maxNrOfIxs := minInt(N, M)
	coordinates := make(chan coordinate, maxNrOfIxs)
	coordsDone := make(chan coordinate, maxNrOfIxs)

	P := randomMatrix(N, K, 10.0)
	Q := randomMatrix(K, M, 10.0)

	for id := range workers {
		go func(id int) {
			for c := range coordinates {
				i, j := c.i, c.j
				e := matrixMultErrorAt(R, P, Q, i, j)
				for k := 0; k < K; k++ {
					P[i][k] += alpha * (2*e*Q[k][j] - beta*P[i][k])
					Q[k][j] += alpha * (2*e*P[i][k] - beta*Q[k][j])
				}
				coordsDone <- coordinate{i, j}
			}
		}(id)
	}

	gen := newCoordinateGenerator(N, M)
	for cnt := 0; cnt < maxNrOfIxs; cnt++ {
		coordinates <- gen.new()
	}
	for step := 0; step < maxSteps; step++ {
		gen.release(<-coordsDone)
		coordinates <- gen.new()
	}
	for cnt := 0; cnt < maxNrOfIxs; cnt++ {
		gen.release(<-coordsDone)
	}
	close(coordinates)
	close(coordsDone)
	return P, Q
}

//------------------------------------------------------------

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	fmt.Println("=== simple matrix factorization ===")
	R := artificialMatrix(4, 4)
	K := 2
	P, Q := matrixFactorization(R, K)
	fmt.Println("err:", factorizationError(R, P, Q))
	fmt.Println("===================================")
}

//------------------------------------------------------------
