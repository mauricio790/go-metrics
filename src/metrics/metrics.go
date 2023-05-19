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
	"fmt"
	"sync"

	metricTypes "metrics/metrictypes"
)

// Type of metric
type MetricType int

// Available metric types
const (
	InvalidMetric MetricType = iota
	Counter
	Fraction
	String
	Time
)

var metricCapabilitiesMap = map[string]MetricType{
	"InvalidMetric": InvalidMetric,
	"Counter":       Counter,
	"Fraction":      Fraction,
	"String":        String,
	"Time":          Time,
}

// Metrics is a struct to keep record of metrics
type Metrics struct {
	mutex      sync.Mutex
	metricData metricTypes.MetricSet
}

// NewMetrics returns a new Metrics struct with the specified metrics
//
// metrics Map that contains the name and type of the metrics to be added
func NewMetrics(metrics map[string]interface{}) Metrics {
	metricSet := metricTypes.NewMetricSet()
	// Initialize metrics values, according to its metric type
	for metricName, metricType := range metrics {
		metricSet.AddMetric(metricName, metricType)
	}

	// Create and return metric structure
	return Metrics{
		metricData: metricSet,
	}
}

// IncreaseMetricValue increases the value of the specified metric
//
// metricName Name of the counter to increase value
// returns error if specified metric does not exist
func (metric Metrics) IncreaseMetricValue(metricName string, increment interface{}) error {
	metric.mutex.Lock()
	defer metric.mutex.Unlock()

	return metric.metricData.IncreaseMetric(metricName, increment)
}

// DecreaseMetricValue decreases the value of the specified metric
//
// metricName Name of the metric to increase value
// returns error if specified metric does not exist
func (metric Metrics) DecreaseMetricValue(metricName string, decrement interface{}) error {
	metric.mutex.Lock()
	defer metric.mutex.Unlock()

	return metric.metricData.DecreaseMetric(metricName, decrement)
}

// ReadMetric returns the value as a string of the specified metric
//
// metricName Name of the metric to be read
// returns error if specified metric does not exist
func (metric Metrics) ReadMetric(metricName string) (interface{}, error) {
	metric.mutex.Lock()
	defer metric.mutex.Unlock()

	value, err := metric.metricData.GetMetricValue(metricName)
	return value, err
}

// ReadMetricAsFloat64 returns the value as a float64 of the specified metric
//
// metricName Name of the metric to be read
// returns error if specified metric does not exist
func (metric Metrics) ReadMetricAsFloat64(metricName string) (float64, error) {
	iValue, err := metric.ReadMetric(metricName)
	if err != nil {
		return 0.0, fmt.Errorf("Unable to read metric |name=%s, value=%v, error=%s", metricName, iValue, err)
	}

	value, ok := iValue.(float64)
	if !ok {
		return 0.0, fmt.Errorf("Unable to cast metric |name=%s, value=%v, type=float64", metricName, iValue)
	}

	return value, err
}

// SetMetric sets the value of a metric
//
// metricName Name of the metric to be read
// value Value to set
// returns error if specified metric does not exist or if value type is invalid
func (metric Metrics) SetMetric(metricName string, value interface{}) error {
	metric.mutex.Lock()
	defer metric.mutex.Unlock()

	return metric.metricData.SetMetricValue(metricName, value)
}

// ResetMetric resets the value of the specified metric
//
// metricName Name of the metric to be reset
// returns error if specified metric does not exist
func (metric Metrics) ResetMetric(metricName string) error {
	metric.mutex.Lock()
	defer metric.mutex.Unlock()

	return metric.metricData.ResetMetric(metricName)
}

// ResetMetric resets the value of all metrics
func (metric Metrics) ResetAllMetrics() {
	metric.mutex.Lock()
	defer metric.mutex.Unlock()

	metric.metricData.ResetAllMetrics()
}

// GetMetricNames returns a slice with all the available metric names
func (metric Metrics) GetMetricNames() []string {
	metric.mutex.Lock()
	defer metric.mutex.Unlock()

	metricNamesList := metric.metricData.GetMetricsNames()

	return metricNamesList
}

// GetMetricType returns the type of a metric
//
// metricName Name of the metric to get the type
// returns error if specified metric does not exist
func (metric Metrics) GetMetricType(metricName string) MetricType {
	metric.mutex.Lock()
	defer metric.mutex.Unlock()

	metricTypeStr := metric.metricData.GetMetricType(metricName)
	// Get the enum based on the string value
	return metricCapabilitiesMap[metricTypeStr]
}

// DeleteMetric removes the metric name from Metrics map
// If specified metric does not exist, does nothing
//
// metricName Name of the metric to be removed
func (metric Metrics) DeleteMetric(metricName string) {
	metric.mutex.Lock()
	defer metric.mutex.Unlock()

	metric.metricData.DeleteMetric(metricName)
}

// GetAllMetrics returns the underlying mapping of metric name to metric value.
func (metric Metrics) GetAllMetrics() map[string]interface{} {
	return metric.metricData.GetAllMetrics()
}
