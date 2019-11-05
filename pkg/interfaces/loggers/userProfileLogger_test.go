package loggers

import (
	"testing"
)

func TestUserProfileLogger(t *testing.T) {
	m := &loggerMock{t: t}
	l := MakeUserProfilePrometheusDefaultLogger(m)
	l.LogBadRequest(nil)
	l.LogErrorGettingInternalData(nil)
	m.AssertExpectations(t)
}
