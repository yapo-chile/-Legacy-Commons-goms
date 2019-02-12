package loggers

import (
	"github.schibsted.io/Yapo/goms/pkg/domain"
	"github.schibsted.io/Yapo/goms/pkg/usecases"
)

type fibonacciPrometheusDefaultLogger struct {
	logger  Logger
	metrics MetricsExporter
}

func (l *fibonacciPrometheusDefaultLogger) LogBadInput(n int) {
	l.metrics.IncrementCounter(BadInputError)
	l.logger.Debug("GetNth doesn't like N < 1. Input: %d", n)
}

func (l *fibonacciPrometheusDefaultLogger) LogRepositoryError(i int, x domain.Fibonacci, err error) {
	l.metrics.IncrementCounter(RepositoryError)
	l.logger.Error("Repository refused to save (%d, %d): %s", i, x, err)
}

// MakeFibonacciPrometheusLogger sets up a FibonacciPrometheusLogger instrumented
// via the provided logger & prometheus metrics exporter
func MakeFibonacciPrometheusLogger(logger Logger, prometheus MetricsExporter) usecases.FibonacciPrometheusLogger {
	return &fibonacciPrometheusDefaultLogger{
		logger:  logger,
		metrics: prometheus,
	}
}
