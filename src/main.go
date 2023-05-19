package main

// Include the telemetry package
import (
	"fmt"
	"time"
	"math/rand"
	//"runtime"
	"metrics/telemetry"
)

func taskA(context telemetry.Context) {
	newContext := telemetry.FunctionName(context)
	if telemetry.IsEnabled() {
		defer telemetry.IncreaseFunctionTracer(newContext, time.Now())
	}

	// Simulate some work
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)	
    taskB(newContext)
}

func taskB(context telemetry.Context) {
	newContext := telemetry.FunctionName(context)
	if telemetry.IsEnabled() {
		defer telemetry.IncreaseFunctionTracer(newContext, time.Now())
	}

	//Simulate some work
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	go taskC(newContext)
}

func taskC(context telemetry.Context) {
	newContext := telemetry.FunctionName(context)
	if telemetry.IsEnabled() {
		defer telemetry.IncreaseFunctionTracer(newContext, time.Now())
	}

	// Simulate some work
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	taskD(newContext)
}

func taskD(context telemetry.Context) {
	newContext := telemetry.FunctionName(context)
	if telemetry.IsEnabled() {
		defer telemetry.IncreaseFunctionTracer(newContext, time.Now())
	}

	// Simulate some work
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
}

func main() {
	// Enable the global metrics
	// Note: this should typically be done via a UnixCTL as
	//       they should start out disabled
	telemetry.Enable()
	context := telemetry.Context{
		ParentFunctionName: "",
		FunctionName:       telemetry.GetRoot(),
		CallID:             "",
	}

	// Simulate some work
	taskA(context)

	// Esperar que termine la go routine
	time.Sleep(time.Duration(1000) * time.Millisecond)
	fmt.Println(telemetry.GetMetricsJSON())
}