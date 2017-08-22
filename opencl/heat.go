package main

import (
	"fmt"
	"github.com/samuel/go-opencl/cl"
	"sync"
	"time"
)

const (
	block   = 16
	N       = block*128 + 2
	mean    = 74.95107632093934
	epsilon = 0.001
	clGPU   = cl.DeviceTypeGPU
	clCPU   = cl.DeviceTypeCPU
)

func main() {
	fmt.Printf("\nmatrix: %v x %v\n\n", N, N)
	fmt.Printf("\ntook %v seconds\n", openClTest(clGPU))
	fmt.Printf("\ntook %v seconds\n", openClTest(clCPU))
	fmt.Printf("\ntook %v seconds\n", test02())
	fmt.Printf("\ntook %v seconds\n", test01())
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

func openClTest(deviceType cl.DeviceType) float64 {
	context, queue := newOpenClContext(deviceType)
	if context == nil {
		fmt.Println("--- Device not found ---", deviceType)
		return 0.0
	}

	mat1 := makeHeatTestMatrix()
	a, _ := context.CreateEmptyBuffer(cl.MemReadWrite, 4*N*N)
	queue.EnqueueWriteBufferFloat32(a, true, 0, mat1[:], nil)

	mat2 := makeHeatTestMatrix()
	b, _ := context.CreateEmptyBuffer(cl.MemReadWrite, 4*N*N)
	queue.EnqueueWriteBufferFloat32(b, true, 0, mat2[:], nil)

	kernel := buildOpenClKernel(
		context, kernelSource, "heatStencil")
	global := []int{N - 2, N - 2}
	local := []int{32, 1}
	cnt := 0
	diff := float32(1.0 + epsilon)
	start := time.Now()
	for diff > epsilon {
		cnt += 1
		for t := 0; t < 500; t++ {
			kernel.SetArgs(b, a)
			queue.EnqueueNDRangeKernel(
				kernel, nil, global, local, nil)
			kernel.SetArgs(a, b)
			queue.EnqueueNDRangeKernel(
				kernel, nil, global, local, nil)
		}
		queue.Finish()
		queue.EnqueueReadBufferFloat32(a, true, 0, mat1, nil)
		queue.EnqueueReadBufferFloat32(b, true, 0, mat2, nil)
		diff = mat1.maxDiff(mat2)
		fmt.Println(
			"iteration: ", 1000*cnt,
			", diff: ", diff,
			" check: ", mat1[10*N+10])
	}
	return time.Now().Sub(start).Seconds()
}

// ----------------------------------------------------------
// --- opencl -----------------------------------------------
// ----------------------------------------------------------

func newOpenClContext(
	deviceType cl.DeviceType) (*cl.Context, *cl.CommandQueue) {
	platforms, err := cl.GetPlatforms()
	if err != nil {
		fmt.Printf("\nFailed to get platforms: %+v", err)
	}
	if len(platforms) < 1 {
		fmt.Printf("\nNo OpenCL Platforms Found ")
		return nil, nil
	}
	platform := platforms[0]

	devices, err := platform.GetDevices(cl.DeviceTypeAll)
	if err != nil {
		fmt.Printf("\nFailed to get devices: %+v", err)
	}
	if len(devices) == 0 {
		fmt.Printf("\nGetDevices returned no devices")
	}
	deviceIndex := -1
	for i, d := range devices {
		if deviceIndex < 0 && d.Type() == deviceType {
			deviceIndex = i
		}
	}
	if deviceIndex < 0 {
		deviceIndex = 0
		return nil, nil
	}
	device := devices[deviceIndex]
	fmt.Printf("\nDevice %d (%s): %s\n",
		deviceIndex, device.Type(), device.Name())
	fmt.Println("--------------")

	context, err := cl.CreateContext([]*cl.Device{device})
	if err != nil {
		fmt.Printf("\nCreateContext failed: %+v", err)
	}
	queue, err := context.CreateCommandQueue(device, 0)
	if err != nil {
		fmt.Printf("\nCreateCommandQueue failed: %+v", err)
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

func buildKernel(program *cl.Program, kernelName string) *cl.Kernel {
	kernel, err := program.CreateKernel(kernelName)
	if err != nil {
		fmt.Printf("\nCreateKernel failed: %+v", err)
	}
	return kernel
}

// --- heat -------------------------------------------------

func test01() float64 {
	fmt.Println("\nNormal Version")
	fmt.Println("--------------")

	u := makeHeatTestMatrix()
	w := makeHeatTestMatrix()

	iterations := 0
	var diff float32 = 1.0 + epsilon
	start := time.Now()
	for epsilon <= diff {
		for t := 0; t < 500; t++ {
			heatStep(w, u)
			heatStep(u, w)
		}
		iterations += 1000
		diff = w.maxDiff(u)
		fmt.Println(
			"iteration: ", iterations,
			", diff: ", diff,
			" check: ", w[10*N+10])
	}
	return time.Now().Sub(start).Seconds()
}

func test02() float64 {
	fmt.Println("\nParallel Version")
	fmt.Println("----------------")

	u := makeHeatTestMatrix()
	w := makeHeatTestMatrix()

	iterations := 0
	var diff float32 = 1.0 + epsilon
	start := time.Now()
	for epsilon <= diff {
		for t := 0; t < 500; t++ {
			pHeatStep(w, u)
			pHeatStep(u, w)
		}
		iterations += 1000
		diff = w.maxDiff(u)
		fmt.Println(
			"iteration: ", iterations,
			", diff: ", diff,
			" check: ", w[10*N+10])
	}
	return time.Now().Sub(start).Seconds()
}

func heatStepLines(w matrix, u matrix, ix int) {
	for i := ix; i < ix+block; i++ {
		for j := 1; j < N-1; j++ {
			w[(i*N)+j] = (u[(i-1)*N+j] + u[(i+1)*N+j] +
				u[(i*N)+j-1] + u[(i*N)+j+1]) / 4.0
		}
	}
}

func heatStep(w matrix, u matrix) {
	for i := 1; i < N-1; i += block {
		heatStepLines(w, u, i)
	}
}

func pHeatStep(w matrix, u matrix) {
	var wg sync.WaitGroup
	wg.Add((N - 2) / block)
	for i := 1; i < N-1; i += block {
		go func(i int) {
			heatStepLines(w, u, i)
			wg.Done()
		}(i)
	}
	wg.Wait()
}

//-----------------------------------------------------------
// --- matrices ---------------------------------------------
//-----------------------------------------------------------

type matrix []float32

func makeHeatTestMatrix() matrix {
	m := make(matrix, N*N)
	m.fill(mean)
	m.fillBorder(100.0)
	return m
}

func (m matrix) fill(v float32) {
	for i, _ := range m {
		m[i] = v
	}
}

func (m matrix) fillBorder(v float32) {
	for i := 0; i < N; i++ {
		m[i*N] = v
		m[i*N+N-1] = v
		m[(N-1)*N+i] = v
		m[i] = 0.0
	}
}

func (m matrix) maxDiff(m2 matrix) float32 {
	result := float32(0.0)
	for ix, v := range m {
		diff := v - m2[ix]
		if diff < 0.0 {
			diff = -diff
		}
		if diff > result {
			result = diff
		}
	}
	return result
}

// ----------------------------------------------------------
