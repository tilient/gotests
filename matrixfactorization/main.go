package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

  . "github.com/tilient/gotests/matrixfactorization/matrix"
)

//------------------------------------------------------------

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	fmt.Println("=== simple matrix factorization ===")
	R := ArtificialMatrix(22, 22)
	K := 3
	P, Q := matrixFactorization(R, K)
	fmt.Println("err:", factorizationError(R, P, Q))
	fmt.Println("===================================")
}

//------------------------------------------------------------

func matrixFactorization(R Matrix, K int) (Matrix, Matrix) {
	const alpha = 0.0002 // the learning rate
	const beta = 0.02    // the regularization parameter

	N := R.NrOfRows()
	M := R.NrOfColumns()
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

	P := RandomMatrix(N, K, 10.0)
	Q := RandomMatrix(K, M, 10.0)

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

func matrixMultErrorAt(R, P, Q Matrix, i, j int) float64 {
	K := Q.NrOfRows()
	e := R[i][j]
	for k := 0; k < K; k++ {
		e -= P[i][k] * Q[k][j]
	}
	return e
}

func factorizationError(R, P, Q Matrix) float64 {
	N := R.NrOfRows()
	M := R.NrOfColumns()
	K := Q.NrOfRows()
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
