package main

import (
	"fmt"
	"github.com/exascience/pargo/parallel"
	"runtime"
	"strings"
	"sync"
	"time"
)

// ----------------------------------------------------------

func main() {
	const N = 1024 + 2
	sep := strings.Repeat("=", 48)
	fmt.Printf("%v\nmatrix: %v x %v (with %v cores)\n%v\n",
		sep, N, N, runtime.NumCPU(), sep)
	heatTest("pargo", N, N, pargoHeatStep)
	heatTest("pargo2", N, N, pargo2HeatStep)
	heatTest("parallel", N, N, parallelHeatStep)
	heatTest("sequential", N, N, sequentialHeatStep)
	fmt.Println(sep)
}

// ----------------------------------------------------------
// --- heat test --------------------------------------------
// ----------------------------------------------------------

func heatTest(title string, M, N int,
	heatStepFun func(u, w *matrix)) {

	sep := strings.Repeat("-", 48)
	fmt.Printf("\n%v\nheat test - %v\n%v\n", sep, title, sep)

	u := makeFilledMatrix(M, N, 75.0)
	u.fillBorders(0.0, 100.0, 100.0, 100.0)
	w := u.copy()

	const ε = float32(0.001)
	δ := ε + 1.0
	iterations := 0
	start := time.Now()
	for δ >= ε {
		for s := 0; s < 1000; s++ {
			heatStepFun(w, u)
			heatStepFun(u, w)
		}
		δ = w.maxDiff(u)
		iterations += 2000
		fmt.Printf(
			"iters: %6d, δ: %08.6f, w[8][8]: %10.8f\n",
			iterations, δ, w.get(8, 8))
	}
	fmt.Printf("\ntook %6.4f seconds\n",
		time.Now().Sub(start).Seconds())
}

// ----------------------------------------------------------

func sequentialHeatStep(w, u *matrix) {
	heatStepRows(w, u, 1, w.nrOfColumns-1)
}

func pargoHeatStep(w, u *matrix) {
	parallel.Range(1, w.nrOfRows-1, 0,
		func(low, high int) {
			heatStepRows(w, u, low, high)
		})
}

func pargo2HeatStep(w, u *matrix) {
	stepFun := func(from, to int) {
		heatStepRows(w, u, from, to)
	}
	parallel.Range(1, w.nrOfRows-1, 0, stepFun)
}

func parallelHeatStep(w, u *matrix) {
	cores := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(cores)
	blocksize := 1 + ((w.nrOfRows - 2) / cores)
	for c := 0; c < cores; c++ {
		go func(c int) {
			from := 1 + c*blocksize
			to := 1 + (c+1)*blocksize
			if to > (w.nrOfRows - 2) {
				to = w.nrOfRows - 1
			}
			heatStepRows(w, u, from, to)
			wg.Done()
		}(c)
	}
	wg.Wait()
}

// ----------------------------------------------------------

func heatStepRows(w, u *matrix, from, to int) {
	for row := from; row < to; row++ {
		heatStepRow(w, u, row)
	}
}

func heatStepRow(w, u *matrix, row int) {
	for col := 1; col < w.nrOfColumns-1; col++ {
		w.set(row, col,
			(u.get(row-1, col)+u.get(row+1, col)+
				u.get(row, col-1)+u.get(row, col+1))/4.0)
	}
}

//-----------------------------------------------------------
// --- matrix -----------------------------------------------
//-----------------------------------------------------------

type matrix struct {
	nrOfRows    int
	nrOfColumns int
	data        []float32
}

func makeMatrix(M, N int) *matrix {
	var mat matrix
	mat.nrOfRows = M
	mat.nrOfColumns = N
	mat.data = make([]float32, M*N, M*N)
	return &mat
}

func makeFilledMatrix(M, N int, v float32) *matrix {
	m := makeMatrix(M, N)
	m.fill(v)
	return m
}

func (m *matrix) copy() *matrix {
	mat := *m
	mat.data = make([]float32, len(m.data), len(m.data))
	copy(mat.data, m.data)
	return &mat
}

func (m *matrix) get(r, c int) float32 {
	return m.data[r*m.nrOfColumns+c]
}

func (m *matrix) set(r, c int, v float32) {
	m.data[r*m.nrOfColumns+c] = v
}

func (m *matrix) fill(v float32) {
	for ix := range m.data {
		m.data[ix] = v
	}
}

func (m *matrix) fillBorders(t, r, b, l float32) {
	for i := 0; i < m.nrOfColumns; i++ {
		m.set(0, i, t)
		m.set(0, m.nrOfRows-1, t)
	}
	for i := 0; i < m.nrOfRows; i++ {
		m.set(i, 0, l)
		m.set(i, m.nrOfColumns-1, l)
	}
}

func (m *matrix) maxDiff(m2 *matrix) (result float32) {
	for ix, v := range m.data {
		result = max(result, abs(v-m2.data[ix]))
	}
	return result
}

// ----------------------------------------------------------

func abs(v float32) float32 {
	if v < 0.0 {
		return -v
	}
	return v
}

func max(a, b float32) float32 {
	if a < b {
		return b
	}
	return a
}

// ----------------------------------------------------------
