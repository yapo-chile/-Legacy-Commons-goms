package loggers

import (
	"testing"

	"github.schibsted.io/Yapo/goms/pkg/domain"
)

// There are no return values to assert on, as logger only cause side effects
// to communicate with the outside world. These tests only ensure that the
// loggers don't panic

func TestWayStairInteractorDefaultLogger(t *testing.T) {
	l := MakeWayStairInteractorLogger(loggerMock{t: t})
	l.LogBadInput(42)
	l.LogRepositoryError(5, domain.WayStair{}, nil)
}
