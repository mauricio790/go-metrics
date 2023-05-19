// (c) Copyright 2022 Hewlett Packard Enterprise Development LP
//
// Confidential computer software. Valid license from Hewlett Packard
// Enterprise required for possession, use or copying.
//
// Consistent with FAR 12.211 and 12.212, Commercial Computer Software,
// Computer Software Documentation, and Technical Data for Commercial Items
// are licensed to the U.S. Government under vendor's standard commercial
// license.

package telemetry

import (
	"encoding/json"
	"runtime"
	"sync"
	"time"
	"strings"
	gometrics "metrics"
)

//////////////////////////////////////////////////////////

type Context struct{
	ParentFunctionName string
	FunctionName       string 
	CallID             string  
}

// Utility functions

// FunctionName Obtains the current function name and its parent function name 
//
// Ex. rdlabs.hpecorp.net/restlib/table.(*RowGetter).GetRows
func FunctionName(context Context) (Context) {

	newContext := Context{
		ParentFunctionName: context.FunctionName,
		FunctionName: getFunctionName(2),
		CallID: context.CallID,
	}

	//When the ID is empty it means we're creating 
	//a new child of the main function so we get a new ID
	if context.CallID == ""{
		stack := make([]byte, 2048)
		length := runtime.Stack(stack, false)
		goRoutineStack := (string)(stack[:length])
		parentFunctionSlice := strings.Fields(goRoutineStack)
		//Add the suffix to the function name to track it in the tree
		newContext.CallID = parentFunctionSlice[len(parentFunctionSlice)-1]
	}
	
	//add id to function names
	newContext.ParentFunctionName = newContext.ParentFunctionName[:strings.Index(newContext.ParentFunctionName,")")+1] + 
									newContext.CallID 
	newContext.FunctionName += newContext.CallID 
	return newContext
}

func getFunctionName(skips int)(functionName string){

	counter, _, _, success := runtime.Caller(skips) //Number of times we want to go up in the stack
	functionName = runtime.FuncForPC(counter).Name() + "()"

	if success{
		return functionName
	}
	return "N/A"
}

//////////////////////////////////////////////////////////

// Telemetry object
type telemetry struct {
	sync.Mutex
	enabled        bool
	functionTracer *gometrics.FunctionTracer
}

// NewTelemetry Create a new Telemetry object
func NewTelemetry() *telemetry {
	return &telemetry{
		Mutex:          sync.Mutex{},
		enabled:        false,
		functionTracer: gometrics.NewFunctionTracer(),
	}
}

// Enable Enable metrics collection by the Telemetry object
func (t *telemetry) Enable() {
	t.Lock()
	defer t.Unlock()

	t.enabled = true
}

// Disable Disable metrics collection by the Telemetry object
func (t *telemetry) Disable() {
	t.Lock()
	defer t.Unlock()

	t.enabled = false
	t.functionTracer.Clear()
}

// Clear Clear metrics collected by the Telemetry object
func (t *telemetry) Clear() {
	t.functionTracer.Clear()
}

// GetMetricsJSON Get a JSON array containing all per-function collected metrics
func (t *telemetry) GetMetricsJSON() string {
	functionTracerMetrics := globalTelemetry.functionTracer.GetFunctionTracerMetrics()
	telemetryMetricsJSON, err := json.Marshal(functionTracerMetrics)
	if err != nil {
		return "{\"error\": \"Could not marshal the telemetry metrics\"}"
	}

	return string(telemetryMetricsJSON)
}

// IsEnabled Get whether metrics collection is enabled or not
func (t *telemetry) IsEnabled() bool {
	// A mutex here won't help _much_ for now but will be costly
	// so leave it unprotected
	return t.enabled
}

func (t* telemetry) SetRoot(root string) {
	t.functionTracer.SetRoot(root)
}

func (t* telemetry) GetRoot() (string) {
	return t.functionTracer.GetRoot()
}

// IncreaseFunctionTracer Increase/update the traced function metrics
func (t *telemetry) IncreaseFunctionTracer(context Context, start time.Time) {
	t.functionTracer.IncreaseFunctionTracer(context.ParentFunctionName, context.FunctionName, start)
}

//////////////////////////////////////////////////////////

// Global Telemetry object
var globalTelemetry *telemetry

// Telemetry public API

// Initialize the global Telemetry object
//
// Note: this is automatically invoked only _once_ by golang
func init() {
	globalTelemetry = NewTelemetry()
}

// Telemetry global instance API

// Leave the infra open for eventual non-global telemetry instances by
// treating even the global one as an object

// Enable Enable global Telemetry metrics collection
func Enable() {
	globalTelemetry.Enable()
	SetRoot()
}

//SetRoot Sets root of the tree
func SetRoot(){
	parentName := getFunctionName(3)
	globalTelemetry.SetRoot(parentName)
	
}

//SetRoot Sets root of the tree
func GetRoot() (string){
	return globalTelemetry.GetRoot()
}


// Disable Disable global Telemetry metrics collection
func Disable() {
	globalTelemetry.Disable()
}

// Clear Clear global Telemetry metrics collected
func Clear() {
	globalTelemetry.Clear()
}

// GetMetricsJSON Get a JSON array containing all global per-function collected metrics
func GetMetricsJSON() string {
	return globalTelemetry.GetMetricsJSON()
}

// IsEnabled Get whether global Telemetry metrics collection is enabled or not
func IsEnabled() bool {
	return globalTelemetry.IsEnabled()
}

// IncreaseFunctionTracer Increase/update global Telemetry traced function metrics
func IncreaseFunctionTracer(context Context, start time.Time) {
	globalTelemetry.IncreaseFunctionTracer(context, start)
}
