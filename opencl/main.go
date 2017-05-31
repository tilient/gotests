package main

import (
  "fmt"
	"math/rand"
  "github.com/samuel/go-opencl/cl"
)

var kernelSource = `
__kernel void square(
   __global float* input,
   __global float* output,
   const unsigned int count)
{
   int i = get_global_id(0);
   if(i < count)
       output[i] = input[i] * input[i];
}
`

func main() {
	var data [1024]float32
	for i := 0; i < len(data); i++ {
		data[i] = rand.Float32()
	}

	platforms, err := cl.GetPlatforms()
	if err != nil {
		fmt.Printf("\nFailed to get platforms: %+v", err)
	}
	for i, p := range platforms {
		fmt.Printf("\nPlatform %d:", i)
		fmt.Printf("\n  Name: %s", p.Name())
		fmt.Printf("\n  Vendor: %s", p.Vendor())
		fmt.Printf("\n  Profile: %s", p.Profile())
		fmt.Printf("\n  Version: %s", p.Version())
		fmt.Printf("\n  Extensions: %s", p.Extensions())
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
		if deviceIndex < 0 && d.Type() == cl.DeviceTypeGPU {
			deviceIndex = i
		}
		fmt.Printf("\nDevice %d (%s): %s", i, d.Type(), d.Name())
		fmt.Printf("\n  Address Bits: %d", d.AddressBits())
		fmt.Printf("\n  Available: %+v", d.Available())
		// fmt.Printf("\n  Built-In Kernels: %s", d.BuiltInKernels())
		fmt.Printf("\n  Compiler Available: %+v", d.CompilerAvailable())
		fmt.Printf("\n  Double FP Config: %s", d.DoubleFPConfig())
		fmt.Printf("\n  Driver Version: %s", d.DriverVersion())
		fmt.Printf("\n  Error Correction Supported: %+v", d.ErrorCorrectionSupport())
		fmt.Printf("\n  Execution Capabilities: %s", d.ExecutionCapabilities())
		fmt.Printf("\n  Extensions: %s", d.Extensions())
		fmt.Printf("\n  Global Memory Cache Type: %s", d.GlobalMemCacheType())
		fmt.Printf("\n  Global Memory Cacheline Size: %d KB", d.GlobalMemCachelineSize()/1024)
		fmt.Printf("\n  Global Memory Size: %d MB", d.GlobalMemSize()/(1024*1024))
		fmt.Printf("\n  Half FP Config: %s", d.HalfFPConfig())
		fmt.Printf("\n  Host Unified Memory: %+v", d.HostUnifiedMemory())
		fmt.Printf("\n  Image Support: %+v", d.ImageSupport())
		fmt.Printf("\n  Image2D Max Dimensions: %d x %d", d.Image2DMaxWidth(), d.Image2DMaxHeight())
		fmt.Printf("\n  Image3D Max Dimenionns: %d x %d x %d", d.Image3DMaxWidth(), d.Image3DMaxHeight(), d.Image3DMaxDepth())
		// fmt.Printf("\n  Image Max Buffer Size: %d", d.ImageMaxBufferSize())
		// fmt.Printf("\n  Image Max Array Size: %d", d.ImageMaxArraySize())
		// fmt.Printf("\n  Linker Available: %+v", d.LinkerAvailable())
		fmt.Printf("\n  Little Endian: %+v", d.EndianLittle())
		fmt.Printf("\n  Local Mem Size Size: %d KB", d.LocalMemSize()/1024)
		fmt.Printf("\n  Local Mem Type: %s", d.LocalMemType())
		fmt.Printf("\n  Max Clock Frequency: %d", d.MaxClockFrequency())
		fmt.Printf("\n  Max Compute Units: %d", d.MaxComputeUnits())
		fmt.Printf("\n  Max Constant Args: %d", d.MaxConstantArgs())
		fmt.Printf("\n  Max Constant Buffer Size: %d KB", d.MaxConstantBufferSize()/1024)
		fmt.Printf("\n  Max Mem Alloc Size: %d KB", d.MaxMemAllocSize()/1024)
		fmt.Printf("\n  Max Parameter Size: %d", d.MaxParameterSize())
		fmt.Printf("\n  Max Read-Image Args: %d", d.MaxReadImageArgs())
		fmt.Printf("\n  Max Samplers: %d", d.MaxSamplers())
		fmt.Printf("\n  Max Work Group Size: %d", d.MaxWorkGroupSize())
		fmt.Printf("\n  Max Work Item Dimensions: %d", d.MaxWorkItemDimensions())
		fmt.Printf("\n  Max Work Item Sizes: %d", d.MaxWorkItemSizes())
		fmt.Printf("\n  Max Write-Image Args: %d", d.MaxWriteImageArgs())
		fmt.Printf("\n  Memory Base Address Alignment: %d", d.MemBaseAddrAlign())
		fmt.Printf("\n  Native Vector Width Char: %d", d.NativeVectorWidthChar())
		fmt.Printf("\n  Native Vector Width Short: %d", d.NativeVectorWidthShort())
		fmt.Printf("\n  Native Vector Width Int: %d", d.NativeVectorWidthInt())
		fmt.Printf("\n  Native Vector Width Long: %d", d.NativeVectorWidthLong())
		fmt.Printf("\n  Native Vector Width Float: %d", d.NativeVectorWidthFloat())
		fmt.Printf("\n  Native Vector Width Double: %d", d.NativeVectorWidthDouble())
		fmt.Printf("\n  Native Vector Width Half: %d", d.NativeVectorWidthHalf())
		fmt.Printf("\n  OpenCL C Version: %s", d.OpenCLCVersion())
		// fmt.Printf("\n  Parent Device: %+v", d.ParentDevice())
		fmt.Printf("\n  Profile: %s", d.Profile())
		fmt.Printf("\n  Profiling Timer Resolution: %d", d.ProfilingTimerResolution())
		fmt.Printf("\n  Vendor: %s", d.Vendor())
		fmt.Printf("\n  Version: %s", d.Version())
	}
	if deviceIndex < 0 {
		deviceIndex = 0
	}
	device := devices[deviceIndex]
	fmt.Printf("\nUsing device %d", deviceIndex)
	context, err := cl.CreateContext([]*cl.Device{device})
	if err != nil {
		fmt.Printf("\nCreateContext failed: %+v", err)
	}
	// imageFormats, err := context.GetSupportedImageFormats(0, MemObjectTypeImage2D)
	// if err != nil {
	// 	fmt.Printf("\nGetSupportedImageFormats failed: %+v", err)
	// }
	// fmt.Printf("\nSupported image formats: %+v", imageFormats)
	queue, err := context.CreateCommandQueue(device, 0)
	if err != nil {
		fmt.Printf("\nCreateCommandQueue failed: %+v", err)
	}
	program, err := context.CreateProgramWithSource([]string{kernelSource})
	if err != nil {
		fmt.Printf("\nCreateProgramWithSource failed: %+v", err)
	}
	if err := program.BuildProgram(nil, ""); err != nil {
		fmt.Printf("\nBuildProgram failed: %+v", err)
	}
	kernel, err := program.CreateKernel("square")
	if err != nil {
		fmt.Printf("\nCreateKernel failed: %+v", err)
	}
	for i := 0; i < 3; i++ {
		name, err := kernel.ArgName(i)
		if err == cl.ErrUnsupported {
			break
		} else if err != nil {
			fmt.Printf("\nGetKernelArgInfo for name failed: %+v", err)
			break
		} else {
			fmt.Printf("\nKernel arg %d: %s", i, name)
		}
	}
	input, err := context.CreateEmptyBuffer(cl.MemReadOnly, 4*len(data))
	if err != nil {
		fmt.Printf("\nCreateBuffer failed for input: %+v", err)
	}
	output, err := context.CreateEmptyBuffer(cl.MemReadOnly, 4*len(data))
	if err != nil {
		fmt.Printf("\nCreateBuffer failed for output: %+v", err)
	}
	if _, err := queue.EnqueueWriteBufferFloat32(input, true, 0, data[:], nil); err != nil {
		fmt.Printf("\nEnqueueWriteBufferFloat32 failed: %+v", err)
	}
	if err := kernel.SetArgs(input, output, uint32(len(data))); err != nil {
		fmt.Printf("\nSetKernelArgs failed: %+v", err)
	}

	local, err := kernel.WorkGroupSize(device)
	if err != nil {
		fmt.Printf("\nWorkGroupSize failed: %+v", err)
	}
	fmt.Printf("\nWork group size: %d", local)
	size, _ := kernel.PreferredWorkGroupSizeMultiple(nil)
	fmt.Printf("\nPreferred Work Group Size Multiple: %d", size)

	global := len(data)
	d := len(data) % local
	if d != 0 {
		global += local - d
	}
	if _, err := queue.EnqueueNDRangeKernel(kernel, nil, []int{global}, []int{local}, nil); err != nil {
		fmt.Printf("\nEnqueueNDRangeKernel failed: %+v", err)
	}

	if err := queue.Finish(); err != nil {
		fmt.Printf("\nFinish failed: %+v", err)
	}

  fmt.Println("---0---")
	results := make([]float32, len(data))
  fmt.Println("---1---", results[:4])
	if _, err := queue.EnqueueReadBufferFloat32(output, true, 0, results, nil); err != nil {
		fmt.Printf("\nEnqueueReadBufferFloat32 failed: %+v", err)
	}
  fmt.Println("---2---", results[:4])

	correct := 0
	for i, v := range data {
		if results[i] == v*v {
			correct++
		}
	}
  fmt.Println("---3---", correct)

	if correct != len(data) {
		fmt.Printf("\n%d/%d correct values", correct, len(data))
	}
}
