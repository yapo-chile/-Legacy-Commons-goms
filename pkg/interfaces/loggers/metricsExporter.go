package loggers

import (
	"io"

	"github.schibsted.io/Yapo/goms/pkg/usecases"
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

// metricsExporter exports custom runtime metrics to prometheus
type metricsExporter struct {
	exporter PrometheusExporter
}

// PrometheusExporter allows operations to export metrics to prometheus
type PrometheusExporter interface {
	IncrementCounter(metric MetricType)
	io.Closer
}

// IncrementBadInputCounter increments bad input counter
func (l *metricsExporter) IncrementBadInputCounter() {
	l.exporter.IncrementCounter(BadInputError)
}

// IncrementRepositoryErrorCounter increments repository error counter
func (l *metricsExporter) IncrementRepositoryErrorCounter() {
	l.exporter.IncrementCounter(RepositoryError)
}

// MakeCustomMetricsExporter sets up a MetricsExporter instrumented
// via the provided metrics exporter
func MakeCustomMetricsExporter(prometheus PrometheusExporter) usecases.MetricsExporter {
	return &metricsExporter{
		exporter: prometheus,
	}
}
