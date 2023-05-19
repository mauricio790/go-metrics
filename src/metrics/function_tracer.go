// (c) Copyright 2022 Hewlett Packard Enterprise Development LP
//
// Confidential computer software. Valid license from Hewlett Packard
// Enterprise required for possession, use or copying.
//
// Consistent with FAR 12.211 and 12.212, Commercial Computer Software,
// Computer Software Documentation, and Technical Data for Commercial Items
// are licensed to the U.S. Government under vendor's standard commercial
// license.

package metrics

import (
	"math"
	"sync"
	"time"
	"strings"
)

// Metrics DTO
type FunctionTracerMetricsDTO struct {
	Parent        string
	Function      string
	Calls         int
	TotalTimeMs   int
	AverageTimeMs float64
	LowerCeiling  int
	HigherCeiling int

	Children []*FunctionTracerMetricsDTO
}

// FunctionTracer maintains metrics for function calls
// ['fn_1']
//    +-- calls
//    +-- total time (ms)
//    ...
// ['fn_2']
//    +-- calls
//    +-- total time (ms)
//    ...
		
type FunctionTracer struct {
	sync.Mutex
	root string
	metrics map[string]FunctionTracerMetricsDTO
}


// Utility functions

// GetSuffix gets the suffix of a function name 
func GetSuffix(functionName string) (suffix string){
	return functionName[strings.Index(functionName,")")+1:]
}

// GetSuffix gets the suffix of a function name 
func GetName(functionName string) (suffix string){
	return functionName[:strings.Index(functionName,")")+1]
}
// NewMetricSet returns a new MetricSet instance
//
// This will contain one set of 5 metrics:
// 	1. Calls
// 	2. Total time (ms)
// 	3. Avg time (ms) -- calculated on the fly
// 	4. Lower ceiling (ms) -- quickest execution
// 	5. Higher ceiling (ms) -- slowest execution
// for every traced function
func NewFunctionTracer() *FunctionTracer {
	functionTracerMetrics := make(map[string]FunctionTracerMetricsDTO)
	return &FunctionTracer{
		Mutex:   sync.Mutex{},
		metrics: functionTracerMetrics,
	}
}

// AddFunctionTraceMetric adds a new function trace metric
//
// If the traced function wasn't present before, it is added
// automatically
//
// functionName Name of the traced function
// start Function call starting time
func (ft *FunctionTracer) IncreaseFunctionTracer(parentFunction string, functionName string, start time.Time) {
	ft.Lock()
	defer ft.Unlock()

	parentFunctionName := GetName(parentFunction)
	// This function will be called on a defer so the start time is calculated
	// during its deferral and time.Now() will be the ending time when it actually
	// gets executed
	functionTimeMs := time.Now().Sub(start).Milliseconds()

	if len(ft.metrics) == 0{
		ft.root = ft.root + GetSuffix(parentFunctionName) //Set the suffix of the root
	}

	if _, ok := ft.metrics[parentFunctionName]; !ok{
		newFunctionMetrics := FunctionTracerMetricsDTO{
			Parent:        parentFunctionName,
			Function:      parentFunctionName,
			Calls:         int(0),
			TotalTimeMs:   int(0),
			AverageTimeMs: float64(0.0),
			LowerCeiling:  int(math.MaxInt32),
			HigherCeiling: int(0),
			Children: []*FunctionTracerMetricsDTO{},
		}
		ft.metrics[parentFunctionName] = newFunctionMetrics
	}

	// Check to see if the traced function already exists in the metrics map
	// and insert it if needed

	newFunctionMetrics := &FunctionTracerMetricsDTO{
			Parent:        parentFunctionName,
			Function:      functionName,
			Calls:         int(0),
			TotalTimeMs:   int(0),
			AverageTimeMs: float64(0.0),
			LowerCeiling:  int(math.MaxInt32),
			HigherCeiling: int(0),
			Children: []*FunctionTracerMetricsDTO{},
		}
	// Time to update the metrics
	//functionMetrics := ft.metrics[parentFunctionName].Children[functionName] // Silly golang won't allow changing the struct directly...

	newFunctionMetrics.Calls += 1
	newFunctionMetrics.TotalTimeMs += int(functionTimeMs)
	// The average is calculated on the fly when returning the metrics
	if functionTimeMs < int64(newFunctionMetrics.LowerCeiling) {
		newFunctionMetrics.LowerCeiling = int(functionTimeMs)
	}
	if functionTimeMs > int64(newFunctionMetrics.HigherCeiling) {
		newFunctionMetrics.HigherCeiling = int(functionTimeMs)
	}

	// Need to replace the map struct here... clunky but required in golang
	children := ft.metrics[parentFunctionName] 
	children.Children = append(children.Children,newFunctionMetrics)
	ft.metrics[parentFunctionName] = children
}

// getAverage Calculate the average
func getAverage(total int, count int) float64 {
	if count <= 0 {
		return 0
	}

	return float64(total) / float64(count)
}

func (ft *FunctionTracer) GetTree (root string,functionCall string) []*FunctionTracerMetricsDTO{
	
	functionChildren := []*FunctionTracerMetricsDTO{}
	if child,ok := ft.metrics[root]; ok{ //Check to see if root function made calls
		for _,metrics := range child.Children{ //Iterate through all the calls
			childInMap := GetName(metrics.Function) + GetSuffix(metrics.Parent)
			if(functionCall == ft.root || GetSuffix(metrics.Function) == GetSuffix(functionCall) ){ //Check to see if the call was made from the same root call	
				functionChildren = append(functionChildren,metrics)
				functionChildren[len(functionChildren)-1].Children = ft.GetTree(childInMap,metrics.Function)
			}
		}
	}
	return functionChildren

}

// GetFunctionTracerMetrics Get a DTO with the current metrics
//
// This can be invoked at any point in time during metrics collection
func (ft *FunctionTracer) GetFunctionTracerMetrics() FunctionTracerMetricsDTO {
	ft.Lock()
	defer ft.Unlock()

	var functionTraceMetricsDTOSlice []FunctionTracerMetricsDTO
	

	tree := FunctionTracerMetricsDTO{
		Parent:        ft.root,
		Function:      ft.root,
		Calls:         int(0),
		TotalTimeMs:   int(0),
		AverageTimeMs: float64(0.0),
		LowerCeiling:  int(math.MaxInt32),
		HigherCeiling: int(0),
		Children: []*FunctionTracerMetricsDTO{},
	}
	tree.Children=ft.GetTree(ft.root,ft.root)
	
	//Time to build the tree
	for _, metrics := range ft.metrics {
		metrics.AverageTimeMs = getAverage(metrics.TotalTimeMs, metrics.Calls)
		functionTraceMetricsDTOSlice = append(functionTraceMetricsDTOSlice, metrics)
	}

	return tree
}

// Clear Clear all metrics collected at this point in time
//
// Any subsequent traced calls invoked after this point will be
// collected
func (ft *FunctionTracer) Clear() {
	ft.Lock()
	defer ft.Unlock()

	// TODO: the 'clear' method is setting the keys to 'nil'
	//       which is wrong. They should be deleted from the
	//       map instead.

	// Remove all keys from the map
	for function := range ft.metrics {
		delete(ft.metrics, function)
	}
}

func (ft *FunctionTracer) SetRoot(root string) {
	ft.Lock()
	defer ft.Unlock()

	ft.root = root 
}

func (ft *FunctionTracer) GetRoot() (root string){
	ft.Lock()
	defer ft.Unlock()

	return ft.root
}