package loggers

import (
	"github.schibsted.io/Yapo/goms/pkg/domain"
	"github.schibsted.io/Yapo/goms/pkg/usecases"
)

type wayStairInteractorDefaultLogger struct {
	logger Logger
}

func (l *wayStairInteractorDefaultLogger) LogBadInput(n int) {
	l.logger.Debug("GetNth doesn't receive inputs N < 1 or N > 14. Input: %d", n)
}

func (l *wayStairInteractorDefaultLogger) LogRepositoryError(i int, x domain.WayStair, err error) {
	l.logger.Error("Repository refused to accept this comb(%d) (%+v): %s", i, x, err)
}

// MakeWayStairInteractorLogger sets up a WayStairInteractorLogger instrumented
// via the provided logger
func MakeWayStairInteractorLogger(logger Logger) usecases.WayStairInteractorLogger {
	return &wayStairInteractorDefaultLogger{
		logger: logger,
	}
}
