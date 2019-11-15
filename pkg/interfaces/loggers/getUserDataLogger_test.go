package loggers

import (
	"testing"
)

func TestGetUserDataLogger(t *testing.T) {
	m := &loggerMock{t: t}
	l := MakeGetUserDataLogger(m)
	l.LogBadInput("")
	m.AssertExpectations(t)
}
