package loggers

import (
	"testing"
)

// There are no return values to assert on, as logger only cause side effects
// to communicate with the outside world. These tests only ensure that the
// loggers don't panic

func TestFibonacciInteractorDefaultLogger(t *testing.T) {
	mock := &loggerMock{t: t}
	l := MakeFibonacciInteractorLogger(mock)
	l.LogBadInput(42)
	l.LogRepositoryError(5, 42, nil)
}
