// (c) Copyright 2022 Hewlett Packard Enterprise Development LP
//
// Confidential computer software. Valid license from Hewlett Packard
// Enterprise required for possession, use or copying.
//
// Consistent with FAR 12.211 and 12.212, Commercial Computer Software,
// Computer Software Documentation, and Technical Data for Commercial Items
// are licensed to the U.S. Government under vendor's standard commercial
// license.

package metric_types

import (
	"time"
	"sync"

	errWrap "metrics/error"
)

const (
	integerStr       = "Int"
	intTypeStr       = "Counter"
	floatStr         = "Float64"
	floatTypeStr     = "Fraction"
	timeStr          = "Time"
	stringStr        = "String"
	invalidMetricStr = "InvalidMetric"
	incMetricFnName  = "IncreaseMetric"
	decMetricFnName  = "DecreaseMetric"
)

// MetricSet contains a map of ints to use as metrics
type MetricSet struct {
	metrics map[string]interface{}
	sync.RWMutex
}

// NewMetricSet returns a new MetricSet instance
func NewMetricSet() MetricSet {
	metric := make(map[string]interface{})
	return MetricSet{
		metrics: metric,
	}
}

// AddMetric adds a new metric to the MetricSet map
//
// metricName Name of the metric to be added
// value Value to initialize the added metric
func (c MetricSet) AddMetric(metricName string, value interface{}) {
	c.Lock()
	defer c.Unlock()
	c.metrics[metricName] = value
}

// DeleteMetric removes a metric from the MetricSet map
// If specified metric does not exist, does nothing
//
// metricName Name of the metric to be deleted
func (c MetricSet) DeleteMetric(metricName string) {
	c.Lock()
	defer c.Unlock()
	delete(c.metrics, metricName)
}

// IncreaseMetric increases the value of a metric by the increment specified
//
// metricName Name of the metric to increase value
// increment Quantity to add to metric value
// returns error if specified metric does not exist
func (c MetricSet) IncreaseMetric(metricName string, increment interface{}) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.metrics[metricName]; !ok {
		return errWrap.MetricNotFound{metricName}
	}

	// Check the both values are the same type
	switch c.metrics[metricName].(type) {
	case int:
		incValue, ok := increment.(int)
		if !ok {
			return errWrap.MetricInvalidType{metricName, integerStr}
		}
		metricValue := c.metrics[metricName].(int)
		c.metrics[metricName] = metricValue + incValue

	case float64:
		incValue, ok := increment.(float64)
		if !ok {
			return errWrap.MetricInvalidType{metricName, floatStr}
		}
		metricValue := c.metrics[metricName].(float64)
		c.metrics[metricName] = metricValue + incValue

	case string:
		incValue, ok := increment.(string)
		if !ok {
			return errWrap.MetricInvalidType{metricName, stringStr}
		}
		metricValue := c.metrics[metricName].(string)
		c.metrics[metricName] = metricValue + incValue

	case time.Time:
		return errWrap.MetricInvalidOperation{metricName, timeStr, incMetricFnName}

	default:
		return errWrap.MetricNotFound{metricName}
	}

	return nil
}

// DecreaseMetric decreases the value of a metric by the decrement specified
//
// metricName Name of the metric to increase value
// decrement Quantity to subtract to metric value
// returns error if specified metric does not exist
func (c MetricSet) DecreaseMetric(metricName string, decrement interface{}) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.metrics[metricName]; !ok {
		return errWrap.MetricNotFound{metricName}
	}

	// Check the both values are the same type
	switch c.metrics[metricName].(type) {
	case int:
		decValue, ok := decrement.(int)
		if !ok {
			return errWrap.MetricInvalidType{metricName, integerStr}
		}
		metricValue := c.metrics[metricName].(int)
		c.metrics[metricName] = metricValue - decValue

	case float64:
		decValue, ok := decrement.(float64)
		if !ok {
			return errWrap.MetricInvalidType{metricName, floatStr}
		}
		metricValue := c.metrics[metricName].(float64)
		c.metrics[metricName] = metricValue - decValue

	case string:
		return errWrap.MetricInvalidOperation{metricName, stringStr, decMetricFnName}

	case time.Time:
		return errWrap.MetricInvalidOperation{metricName, timeStr, decMetricFnName}

	default:
		return errWrap.MetricNotFound{metricName}
	}

	return nil
}

// ResetMetric sets the value of a metric to nil
//
// metricName Name of the metric to increase value
// returns error if specified metric does not exist
func (c MetricSet) ResetMetric(metricName string) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.metrics[metricName]; !ok {
		return errWrap.MetricNotFound{metricName}
	}
	c.metrics[metricName] = nil
	return nil
}

// ResetAllMetrics sets the value of all metrics to nil
func (c MetricSet) ResetAllMetrics() {
	c.Lock()
	defer c.Unlock()
	for metricName := range c.metrics {
		c.metrics[metricName] = nil
	}
}

// GetMetricValue returns the value of a metric
//
// metricName Name of the metric to get value
// returns error if specified metric does not exist
func (c MetricSet) GetMetricValue(metricName string) (interface{}, error) {
	c.RLock()
	defer c.RUnlock()
	if _, ok := c.metrics[metricName]; !ok {
		return nil, errWrap.MetricNotFound{metricName}
	}
	return c.metrics[metricName], nil
}

// SetMetricValue sets the value of a metric
//
// metricName Name of the metric to get value
// value Value to set the metric
// returns error if specified metric does not exist
func (c MetricSet) SetMetricValue(metricName string, value interface{}) error {
	c.Lock()
	defer c.Unlock()
	if _, ok := c.metrics[metricName]; !ok {
		return errWrap.MetricNotFound{metricName}
	}
	c.metrics[metricName] = value
	return nil
}

// GetAllMetrics returns all metrics in a map with the
// metric name as key
func (c MetricSet) GetAllMetrics() map[string]interface{} {
	c.RLock()
	defer c.RUnlock()
	return c.metrics
}

// GetMetricsNames returns a slice with the name of all metrics
func (c MetricSet) GetMetricsNames() []string {
	c.RLock()
	defer c.RUnlock()
	metricsNames := []string{}
	for name := range c.metrics {
		metricsNames = append(metricsNames, name)
	}
	return metricsNames
}

// GetMetricType returns the type of a metric
//
// metricName Name of the metric to get the type
// returns error if specified metric does not exist
func (c MetricSet) GetMetricType(metricName string) string {
	c.RLock()
	defer c.RUnlock()
	if _, ok := c.metrics[metricName]; !ok {
		return invalidMetricStr
	}

	// Check the both values are the same type
	switch c.metrics[metricName].(type) {
	case int:
		return intTypeStr

	case float64:
		return floatTypeStr

	case string:
		return stringStr

	case time.Time:
		return timeStr

	default:
		return invalidMetricStr
	}
}
