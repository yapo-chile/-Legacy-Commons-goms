package loggers

import (
	"io"
)

const (
	// EventSuccess represents a successful event type
	EventSuccess = "success"
	// EventError represents an error event type
	EventError = "error"
	// EventInfo represents an info event type
	EventInfo = "info"
	// EventWarning represents an warning event type
	EventWarning = "warning"
	// EventCrit represents an critical event type
	EventCrit = "critical"
)

// Metrics contains needed structs to export metrics data to prometheus
type Metrics struct {
	collector CounterVector
	exporter  MetricsExporter
}

// MetricsExporter allows operations to export metrics to prometheus
type MetricsExporter interface {
	NewCounterVector(name, help string, labels []string) CounterVector
	io.Closer
}

// CounterVector allows operations for data collector using counter selected by labels
type CounterVector interface {
	WithLabelValues(labels ...string) Counter
}

// Counter allows operation for prometheus counter. Counter represents a single
// numerical value that only ever goes up
type Counter interface {
	// Inc increments the counter by 1. Use Add to increment it by arbitrary
	// non-negative values
	Inc()
	// Add adds the given value to the counter. It panics if the value is < 0
	Add(float64)
}
