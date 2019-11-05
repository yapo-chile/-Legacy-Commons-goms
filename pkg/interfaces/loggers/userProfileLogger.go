package loggers

import (
	"github.mpi-internal.com/Yapo/goms/pkg/interfaces/handlers"
)

type userProfilePrometheusDefaultLogger struct {
	logger Logger
}

func (l *userProfilePrometheusDefaultLogger) LogBadRequest(input interface{}) {
	l.logger.Error("Bad request with input: %+v", input)
}

func (l *userProfilePrometheusDefaultLogger) LogErrorGettingInternalData(err error) {
	l.logger.Error("Error getting internal data %+v ", err)
}

// MakeUserProfilePrometheusDefaultLogger sets up a InternalUserDataHandlerLogger instrumented
// via the provided logger
func MakeUserProfilePrometheusDefaultLogger(logger Logger) handlers.UserProfilePrometheusDefaultLogger {
	return &userProfilePrometheusDefaultLogger{
		logger: logger,
	}
}
