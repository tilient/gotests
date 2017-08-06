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

func testFun(a, b, c, d float64) float64 {
	return math.Cos(a) + 2*b*c + math.Sin(d)
}

func test03() {
	n := 128
	m := 4 * 128

	A := mat.NewDense(n, m, nil)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			A.Set(i, j, 2.0*rand.Float64()-1.0)
		}
	}
	y := make([]float64, n)
	for i := 0; i < n; i++ {
		//y[i] = 2.0*rand.Float64() - 1.0
		y[i] = testFun(
			A.At(i, 0), A.At(i, 1),
			A.At(i, 2), A.At(i, 3))
	}
	Y := mat.NewVecDense(n, y)

	x := make([]float64, m)
	X := mat.NewVecDense(m, x)
	yp := make([]float64, n)
	Yp := mat.NewVecDense(n, yp)

	fmt.Println("Solving A * X = Y")
	err := X.SolveVec(A, Y)
	if err != nil {
		fmt.Println("ERR -", err)
	}

	Yp.MulVec(A, X)
	fmt.Println("x: ", x[:3])
	fmt.Println("y: ", y[:3])
	fmt.Println("yp:", yp[:3])
	delta := 0.1
	for mat.EqualApprox(Y, Yp, delta) {
		delta /= 10.0
	}
	fmt.Println("delta:", 10.0*delta)
}
