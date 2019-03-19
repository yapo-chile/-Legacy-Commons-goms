package loggers

import (
	"github.schibsted.io/Yapo/goms/pkg/domain"
	"github.schibsted.io/Yapo/goms/pkg/usecases"
)

type fibonacciPrometheusDefaultLogger struct {
	logger  Logger
	metrics Metrics
}

func (l *fibonacciPrometheusDefaultLogger) LogBadInput(n int) {
	l.metrics.collector.CollectEvent("bad_input", EventError)
	l.logger.Debug("GetNth doesn't like N < 1. Input: %d", n)
}

func (l *fibonacciPrometheusDefaultLogger) LogRepositoryError(i int, x domain.Fibonacci, err error) {
	l.metrics.collector.CollectEvent("repository", EventError)
	l.logger.Error("Repository refused to save (%d, %d): %s", i, x, err)
}

// MakeFibonacciPrometheusLogger sets up a FibonacciPrometheusLogger instrumented
// via the provided logger & prometheus metrics exporter
func MakeFibonacciPrometheusLogger(logger Logger, prometheus MetricsExporter) usecases.FibonacciPrometheusLogger {
	collector := prometheus.NewEventsCollector(
		"goms_fibonacci_events_total", // metric name
		"fibonacci events counter",    // metric help
	)
	return &fibonacciPrometheusDefaultLogger{
		logger: logger,
		metrics: Metrics{
			collector: collector,
			exporter:  prometheus,
		},
	}
}
