package loggers

import (
	"github.mpi-internal.com/Yapo/goms/pkg/usecases"
)

type userProfilePrometheusDefaultLogger struct {
	logger Logger
}

func (l *userProfilePrometheusDefaultLogger) LogBadInput(n string) {
	l.logger.Error("Wrong input type: %d", n)
}

func (l *userProfilePrometheusDefaultLogger) LogRepositoryError(i string, x usecases.UserBasicData, err error) {
	l.logger.Error("Repository refused to save (%d, %d): %s", i, x, err)
}

// MakeUserProfileLogger sets up a userProfileLogger instrumented via the provided logger
func MakeUserProfileLogger(logger Logger) usecases.UserProfilerometheusLogger {
	return &userProfilePrometheusDefaultLogger{
		logger: logger,
	}
}
