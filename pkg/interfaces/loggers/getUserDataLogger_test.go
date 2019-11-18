package loggers

import (
	"testing"
)

func TestGetUserDataLogger(t *testing.T) {
	m := &loggerMock{t: t}
	l := MakeGetUserLogger(m)
	l.LogBadInput("")
	m.AssertExpectations(t)
}
