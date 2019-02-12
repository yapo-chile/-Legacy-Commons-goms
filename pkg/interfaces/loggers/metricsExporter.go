package loggers

import (
	"io"
)

// MetricType is the datatype to represent a prometheus metric type
type MetricType int

// Metrics exposer constants, add all the metrics that you need
const (
	// BadInputError represents counter of bad input errors
	BadInputError MetricType = iota
	// RepositoryError represents counter of repository errors
	RepositoryError
)

// MetricsExporter allows operations to export metrics to prometheus
type MetricsExporter interface {
	IncrementCounter(metric MetricType)
	io.Closer
}
