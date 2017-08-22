package main

import (
	"fmt"
	"github.com/exascience/pargo/parallel"
	"github.com/samuel/go-opencl/cl"
	"runtime"
	"strings"
	"sync"
	"time"
)

func main() {
	const N = 512 + 2
	sep := strings.Repeat("=", 48)
	fmt.Printf("%v\nmatrix: %v x %v\n%v\n", sep, N, N, sep)
	heatTestOpenCl(N, N, cl.DeviceTypeGPU)
	heatTestOpenCl(N, N, cl.DeviceTypeCPU)
	heatTest("pargo", N, N, pargoHeatStep)
	heatTest("pargo2", N, N, pargoHeatStep2)
	heatTest("parallel", N, N, parallelHeatStep)
	heatTest("sequential", N, N, sequentialHeatStep)
	fmt.Println(sep)
}

// ----------------------------------------------------------
// --- heat test --------------------------------------------
// ----------------------------------------------------------

const ε = float32(0.001)

func heatTest(title string, M, N int,
	heatStepFun func(u, w *matrix)) {

	sep := strings.Repeat("-", 48)
	fmt.Printf("\n%v\n%v\n%v\n", sep, title, sep)

	u := makeHeatTestMatrix(M, N)
	w := makeHeatTestMatrix(M, N)

	iterations := 0
	δ := ε + 1.0
	start := time.Now()
	for δ >= ε {
		for t := 0; t < 1000; t++ {
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

func parallelHeatStep(w, u *matrix) {
	cores := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(cores)
	blocksize := (w.nrOfRows - 2) / cores
	for row := 1; row < w.nrOfRows-1; row += blocksize {
		go func(from, to int) {
			heatStepRows(w, u, from, to)
			wg.Done()
		}(row, row+blocksize)
	}
	wg.Wait()
}

func pargoHeatStep(w, u *matrix) {
	parallel.Range(1, w.nrOfRows-1, 0,
		func(low, high int) {
			heatStepRows(w, u, low, high)
		})
}

func pargoHeatStep2(w, u *matrix) {
	parallel.Range(1, w.nrOfRows-1, 0, heatStepRowsOn(w, u))
}

// ----------------------------------------------------------

func heatStepRowsOn(w, u *matrix) func(from, to int) {
	return func(from, to int) {
		heatStepRows(w, u, from, to)
	}
}

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

func makeHeatTestMatrix(M, N int) *matrix {
	var mat matrix
	mat.nrOfRows = M
	mat.nrOfColumns = N
	mat.data = make([]float32, M*N)
	mat.fill(74.95107632093934)
	mat.fillBorders(0.0, 100.0, 100.0, 100.0)
	return &mat
}

func (m *matrix) get(r, c int) float32 {
	return m.data[r*m.nrOfColumns+c]
}

func (m *matrix) set(r, c int, v float32) {
	m.data[r*m.nrOfColumns+c] = v
}

func (m *matrix) fill(v float32) {
	for i, _ := range m.data {
		m.data[i] = v
	}
}

func (m *matrix) fillBorders(t, r, b, l float32) {
	N := m.nrOfColumns
	M := m.nrOfRows
	for i := 0; i < N; i++ {
		m.data[i] = t
		m.data[(M-1)*N+i] = b
	}
	for i := 0; i < M; i++ {
		m.data[i*N] = l
		m.data[i*N+N-1] = r
	}
}

func (m *matrix) maxDiff(m2 *matrix) (result float32) {
	for ix, v := range m.data {
		δ := v - m2.data[ix]
		if δ < 0.0 {
			δ = -δ
		}
		if δ > result {
			result = δ
		}
	}
	return result
}

// ----------------------------------------------------------
// --- opencl test ------------------------------------------
// ----------------------------------------------------------

var kernelSource = `
__kernel void heatStencil(
	__global float* a, __global float* b)
{
  int x = 1 + get_global_id(0);
  int y = 1 + get_global_id(1);
  int X = 2 + get_global_size(0);
  int Y = 2 + get_global_size(1);
  int yX = y * X;
  a[x + yX] = (b[x - 1 + yX] + b[x + 1 + yX] +
               b[x + yX - X] + b[x + yX + X]) / 4.0;
}
`

func heatTestOpenCl(M, N int, deviceType cl.DeviceType) {
	context, queue := newOpenClContext(deviceType)
	if context == nil {
		fmt.Println("--- Device not found ---", deviceType)
		return
	}

	w := makeHeatTestMatrix(M, N)
	u := makeHeatTestMatrix(M, N)

	a, _ := context.CreateEmptyBuffer(cl.MemReadWrite, 4*M*N)
	b, _ := context.CreateEmptyBuffer(cl.MemReadWrite, 4*M*N)
	queue.EnqueueWriteBufferFloat32(
		a, true, 0, w.data[:], nil)
	queue.EnqueueWriteBufferFloat32(
		b, true, 0, u.data[:], nil)

	kernel := buildOpenClKernel(
		context, kernelSource, "heatStencil")
	global := []int{N - 2, M - 2}
	local := []int{32, 1}
	cnt := 0
	δ := ε + 1.0
	start := time.Now()
	for δ > ε {
		cnt += 1
		for t := 0; t < 1000; t++ {
			kernel.SetArgs(b, a)
			queue.EnqueueNDRangeKernel(
				kernel, nil, global, local, nil)
			kernel.SetArgs(a, b)
			queue.EnqueueNDRangeKernel(
				kernel, nil, global, local, nil)
		}
		queue.Finish()
		queue.EnqueueReadBufferFloat32(a, true, 0, w.data, nil)
		queue.EnqueueReadBufferFloat32(b, true, 0, u.data, nil)
		δ = w.maxDiff(u)
		fmt.Printf(
			"iters: %6d, δ: %08.6f, w[8][8]: %10.8f\n",
			2000*cnt, δ, w.get(8, 8))
	}
	fmt.Printf("\ntook %6.4f seconds\n",
		time.Now().Sub(start).Seconds())
}

func newOpenClContext(
	deviceType cl.DeviceType) (*cl.Context, *cl.CommandQueue) {
	platforms, err := cl.GetPlatforms()
	if err != nil {
		fmt.Printf("\nFailed to get platforms: %+v", err)
		return nil, nil
	}
	if len(platforms) < 1 {
		fmt.Printf("\nNo OpenCL Platforms Found ")
		return nil, nil
	}
	platform := platforms[0]

	devices, err := platform.GetDevices(cl.DeviceTypeAll)
	if err != nil {
		fmt.Printf("\nFailed to get devices: %+v", err)
		return nil, nil
	}
	if len(devices) == 0 {
		fmt.Printf("\nGetDevices returned no devices")
		return nil, nil
	}
	deviceIndex := -1
	for i, d := range devices {
		if deviceIndex < 0 && d.Type() == deviceType {
			deviceIndex = i
		}
	}
	if deviceIndex < 0 {
		fmt.Printf("\nDid not find right device")
		return nil, nil
	}
	device := devices[deviceIndex]
	sep := strings.Repeat("-", 48)
	fmt.Printf("\n%v\n%v: %v\n%v\n",
		sep, device.Type(), device.Name(), sep)

	context, err := cl.CreateContext([]*cl.Device{device})
	if err != nil {
		fmt.Printf("\nCreateContext failed: %+v", err)
		return nil, nil
	}
	queue, err := context.CreateCommandQueue(device, 0)
	if err != nil {
		fmt.Printf("\nCreateCommandQueue failed: %+v", err)
		return nil, nil
	}
	return context, queue
}

func buildOpenClKernel(
	context *cl.Context, source string,
	kernelName string) *cl.Kernel {
	program := buildOpenClProgram(context, kernelSource)
	kernel := buildKernel(program, "heatStencil")
	return kernel
}

func buildOpenClProgram(
	context *cl.Context, source string) *cl.Program {
	program, err := context.CreateProgramWithSource(
		[]string{source})
	if err != nil {
		fmt.Printf("\nCreateProgramWithSource failed: %+v", err)
	}
	if err := program.BuildProgram(nil, ""); err != nil {
		fmt.Printf("\nBuildProgram failed: %+v", err)
	}
	return program
}

func buildKernel(
	program *cl.Program, kernelName string) *cl.Kernel {
	kernel, err := program.CreateKernel(kernelName)
	if err != nil {
		fmt.Printf("\nCreateKernel failed: %+v", err)
	}
	return kernel
}

// ----------------------------------------------------------
