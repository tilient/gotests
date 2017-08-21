package main

import (
	"fmt"
	"github.com/samuel/go-opencl/cl"
)

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

const (
	N       = 4*512 + 2
	mean    = 74.95107632093934
	epsilon = 0.001
)

type matrix []float32

func main() {

	context, queue := newOpenClContext()
	program := buildOpenClProgram(context, kernelSource)
	kernel := buildKernel(program, "heatStencil")

	data := make(matrix, N*N)
	data2 := make(matrix, N*N)
	data.fill(mean)
	data.fillBorder(100.0)
	data2.fill(mean)
	data2.fillBorder(100.0)
	fmt.Println("---1---", data[10*N+10])

	a, _ := context.CreateEmptyBuffer(cl.MemReadWrite, 4*N*N)
	queue.EnqueueWriteBufferFloat32(a, true, 0, data[:], nil)
	b, _ := context.CreateEmptyBuffer(cl.MemReadWrite, 4*N*N)
	queue.EnqueueWriteBufferFloat32(b, true, 0, data[:], nil)

	global := []int{N - 2, N - 2}
	local := []int{32, 1}
	cnt := 0
	diff := float32(1.0 + epsilon)
	for diff > epsilon {
		cnt += 1
		for t := 0; t < 500; t++ {
			kernel.SetArgs(b, a)
			queue.EnqueueNDRangeKernel(kernel, nil, global, local, nil)
			kernel.SetArgs(a, b)
			queue.EnqueueNDRangeKernel(kernel, nil, global, local, nil)
		}
		queue.Finish()
		queue.EnqueueReadBufferFloat32(a, true, 0, data, nil)
		queue.EnqueueReadBufferFloat32(b, true, 0, data2, nil)
		diff = data.maxDiff(data2)
		fmt.Println(cnt, "---", diff, " --- ", data[10*N+10])
	}

	queue.EnqueueReadBufferFloat32(a, true, 0, data, nil)
	fmt.Println("---2---", data[10*N+10])
}

func newOpenClContext() (*cl.Context, *cl.CommandQueue) {
	platforms, err := cl.GetPlatforms()
	if err != nil {
		fmt.Printf("\nFailed to get platforms: %+v", err)
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
	fmt.Printf("\n--- Candidates ---")
	for i, d := range devices {
		fmt.Printf("\nDevice %d (%s): %s", i, d.Type(), d.Name())
		if deviceIndex < 0 && d.Type() == cl.DeviceTypeGPU {
			deviceIndex = i
		}
	}
	if deviceIndex < 0 {
		deviceIndex = 0
	}
	deviceIndex = 0
	device := devices[deviceIndex]
	fmt.Printf("\n--- Selected ---")
	fmt.Printf("\nDevice %d (%s): %s",
		deviceIndex, device.Type(), device.Name())
	fmt.Printf("\n--- --- ---\n")

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

func buildOpenClProgram(context *cl.Context, source string) *cl.Program {
	program, err := context.CreateProgramWithSource([]string{source})
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
