package loggers

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

// There are no return values to assert on, as logger only cause side effects
// to communicate with the outside world. These tests only ensure that the
// loggers don't panic

type prometheusMock struct {
	mock.Mock
}

func (m *prometheusMock) Close() error {
	args := m.Called()
	return args.Error(0)
}

func (m *prometheusMock) IncrementCounter(metric MetricType) {
	m.Called(metric)
}

func TestFibonacciInteractorDefaultLogger(t *testing.T) {
	m := &loggerMock{t: t}
	mPrometheus := &prometheusMock{}
	mPrometheus.On("IncrementCounter", BadInputError).Once()
	mPrometheus.On("IncrementCounter", RepositoryError).Once()
	l := MakeFibonacciPrometheusLogger(m, mPrometheus)
	l.LogBadInput(42)
	l.LogRepositoryError(5, 42, nil)
	m.AssertExpectations(t)
}
