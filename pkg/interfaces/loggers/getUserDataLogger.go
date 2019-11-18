package loggers

import (
	"github.mpi-internal.com/Yapo/goms/pkg/usecases"
)

type getUserDataPrometheusDefaultLogger struct {
	logger Logger
}

func (l *getUserDataPrometheusDefaultLogger) LogBadInput(n string) {
	l.logger.Error("Wrong input type: %d", n)
}

// MakeGetUserDataLogger sets up a userProfileLogger instrumented via the provided logger
func MakeGetUserDataLogger(logger Logger) usecases.GetUserDataPrometheusDefaultLogger {
	return &getUserDataPrometheusDefaultLogger{
		logger: logger,
	}
}
