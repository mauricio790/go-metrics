//(c) Copyright 2022 Hewlett Packard Enterprise Development LP
//All Rights Reserved.
//
//The contents of this software are proprietary and confidential
//to the Hewlett Packard Enterprise Development LP.  No part of this
//program may be photocopied, reproduced, or translated into another
//programming language without prior written consent of the
//Hewlett Packard Enterprise Development LP.

package error

var (
	MetricNotFoundMsg         = "Metric was not found | name=%s |"
	MetricInvalidTypeMsg      = "Metric does not match with value to update | name=%s, type=%s |"
	MetricInvalidOperationMsg = "Metric does not support operation | name=%s, type=%s, operation=%s |"
	CounterNotFoundMsg        = "Counter was not found | name=%s |"
	ValueAssertionInvalidMsg  = "Metric data could not be asserted | value=%v, type=%s |"
)
