package loggers

import (
	"io"
)

// EventType represents the kind of the event exported to prometheus.
// It will works as filter in grafana dashboard
type EventType string

const (
	// EventError represents an error event type
	EventError EventType = "error"
)

// Metrics contains needed structs to export metrics data to prometheus
type Metrics struct {
	// collector is a prometheus data struct to collect metrics and display it on
	// HOST:PROMETHEUS_PORT/metrics endpoint. Add any other collector if you need it.
	collector EventsCollector
	// exporter is a prometheus instance to allow collectors creation. Also allows
	// close the metrics prometheus exporter
	exporter MetricsExporter
}

// MetricsExporter allows operations to export metrics to prometheus.
// Please before naming your metric take a look here:
// https://confluence.schibsted.io/pages/viewpage.action?spaceKey=SPTINF&title=Common+Metrics+Standard
type MetricsExporter interface {
	io.Closer
	NewEventsCollector(name, help string) EventsCollector
}

// EventsCollector allows operations for data collector
type EventsCollector interface {
	CollectEvent(eventName string, eventType EventType)
}
