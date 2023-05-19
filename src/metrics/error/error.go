//(c) Copyright 2022 Hewlett Packard Enterprise Development LP
//All Rights Reserved.
//
//The contents of this software are proprietary and confidential
//to the Hewlett Packard Enterprise Development LP.  No part of this
//program may be photocopied, reproduced, or translated into another
//programming language without prior written consent of the
//Hewlett Packard Enterprise Development LP.

package error

import (
	"fmt"
)

// MetricNotFound represents an error when a metric name is
// not found among the available metrics
type MetricNotFound struct {
	MetricName string
}

// CounterNotFound represents an error when a counter name is
// not found among the available counter
type CounterNotFound struct {
	CounterName string
}

// MetricInvalidType represents an error when a metric type is
// missmatch among the metric to be updated
type MetricInvalidType struct {
	MetricName string
	MetricType string
}

// MetricInvalidOperation represents an error when a metric type does
// not support metric function
type MetricInvalidOperation struct {
	MetricName string
	MetricType string
	MetricOperation string
}

// ValueAssertionInvalid represents an error when an interface{} value
// can not be converted to a valid value for a metric
type ValueAssertionInvalid struct {
	Value        interface{}
	ExpectedType string
}

// MetricNotFound implements the error interface
func (e MetricNotFound) Error() string {
	err := "Error: " + fmt.Sprintf(MetricNotFoundMsg, e.MetricName)
	return err
}

// CounterNotFound implements the error interface
func (e CounterNotFound) Error() string {
	err := "Error: " + fmt.Sprintf(CounterNotFoundMsg, e.CounterName)
	return err
}

// MetricInvalidType implements the error interface
func (e MetricInvalidType) Error() string {
	err := "Error: " + fmt.Sprintf(MetricInvalidTypeMsg, e.MetricName, e.MetricType)
	return err
}

// MetricInvalidType implements the error interface
func (e MetricInvalidOperation) Error() string {
	err := "Error: " + fmt.Sprintf(MetricInvalidOperationMsg, e.MetricName, e.MetricType, e.MetricOperation)
	return err
}

// ValueAssertionInvalid implements the error interface
func (e ValueAssertionInvalid) Error() string {
	err := "Error: " + fmt.Sprintf(ValueAssertionInvalidMsg, e.Value, e.ExpectedType)
	return err
}
