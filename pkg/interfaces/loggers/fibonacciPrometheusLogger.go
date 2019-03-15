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
	l.metrics.collector.WithLabelValues("bad_input", EventError).Inc()
	l.logger.Debug("GetNth doesn't like N < 1. Input: %d", n)
}

func (l *fibonacciPrometheusDefaultLogger) LogRepositoryError(i int, x domain.Fibonacci, err error) {
	l.metrics.collector.WithLabelValues("repository", EventError).Inc()
	l.logger.Error("Repository refused to save (%d, %d): %s", i, x, err)
}

// MakeFibonacciPrometheusLogger sets up a FibonacciPrometheusLogger instrumented
// via the provided logger & prometheus metrics exporter
func MakeFibonacciPrometheusLogger(logger Logger, prometheus MetricsExporter) usecases.FibonacciPrometheusLogger {
	counterVector := prometheus.NewCounterVector( // all metrics will be stored in a vector of counters
		"goms_fibonnaci_events_count_total", // metric name
		"fibonnaci events counter",          // metric help
		[]string{"event", "type"},           // labels
	)

	return &fibonacciPrometheusDefaultLogger{
		logger: logger,
		metrics: Metrics{
			collector: counterVector,
			exporter:  prometheus,
		},
	}
}
