package main

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"time"

	. "github.com/tilient/gotests/matrixfactorization/matrix"
)

//------------------------------------------------------------

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	fmt.Println("===================================")
	fmt.Println("=== Simple Matrix Factorization ===")
	fmt.Println("===================================")
	n, m := 8, 32
	R := RandomMatrix(m, n)
	for i := 0; i < n; i++ {
		R[0][i] = 0.5 *
			(2*R[1][i] + 3*R[2][i] + 4*math.Sin(R[3][i]))
	}
	K := 8
	P, Q := matrixFactorization(R, K)
	fmt.Println("err =", factorizationError(R, P, Q))
	fmt.Println("-- R ------------------------------")
	R.Print()
	fmt.Println("-- Rp -----------------------------")
	Rp := P.Mult(Q)
	Rp.Print()
	fmt.Println("-- diff ---------------------------")
	diff := Rp.Min(R).Abs()
	diff.Print()
	fmt.Println("-- P ------------------------------")
	P.Print()
	fmt.Println("-- Q ------------------------------")
	Q.Print()
	fmt.Println("===================================")
	fmt.Println("== Prediction                    ==")
	fmt.Println("===================================")
	r := RandomMatrix(m, 1)
	r[0][0] = 0.0
	v := 0.5 * (2*r[1][0] + 3*r[2][0] + 4*math.Sin(r[3][0]))
	fmt.Println("v =", v)
	q := RandomMatrix(K, 1)
	q.Solve(P, r)
	fmt.Println("-- q ------------------------------")
	q.Print()
	fmt.Println("-- rp -- (= P * q) ----------------")
	rp := P.Mult(q)
	rp.Print()
	fmt.Println("-- r ------------------------------")
	r.Print()
	fmt.Println("-----------------------------------")
	for i := 0; i < 10; i++ {
		r = RandomMatrix(m, 1)
		r[0][0] = 0.0
		v = 0.1 * (2*r[1][0] + 3*r[2][0] + 4*math.Sin(r[3][0]))
		q.Solve(P, r)
		rp := P.Mult(q)
		p := rp[0][0]
		//fmt.Println(i, "- v =", v, "p =", p)
		fmt.Println("   [", math.Abs(p-v), "]", v)
	}
	fmt.Println("===================================")
}

//------------------------------------------------------------

func matrixFactorization(R Matrix, K int) (Matrix, Matrix) {
	const alpha = 0.0002 // the learning rate
	const beta = 0.02    // the regularization parameter

	N := R.NrOfRows()
	M := R.NrOfColumns()
	maxSteps := 1000 * K * M * N
	fmt.Println("R =", N, "x", M)
	fmt.Println("K =", K)
	fmt.Println("max steps =", maxSteps)

	nrOfWorkers := 2 * maxInt(1, runtime.NumCPU()-2)
	fmt.Println("nr of workers", nrOfWorkers)
	fmt.Println("-----------------------------------")
	workers := newIntRange(nrOfWorkers)
	maxNrOfIxs := minInt(N, M)
	coordinates := make(chan coordinate, maxNrOfIxs)
	coordsDone := make(chan coordinate, maxNrOfIxs)

	P := RandomMatrix(N, K)
	Q := RandomMatrix(K, M)

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
	return math.Sqrt(e)
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
