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

func (m *prometheusMock) NewCounterVector(name, help string, labels []string) CounterVector {
	args := m.Called(name, help, labels)
	return args.Get(0).(CounterVector)
}

func (m *prometheusMock) Close() error {
	args := m.Called()
	return args.Error(0)
}

type mockCollector struct {
	mock.Mock
}

func (m *mockCollector) WithLabelValues(labels ...string) Counter {
	args := m.Called(labels)
	return args.Get(0).(Counter)
}

type mockCounter struct {
	mock.Mock
}

func (m *mockCounter) Inc() {
	m.Called()
}

func (m *mockCounter) Add(value float64) {
	m.Called(value)
}

func TestFibonacciInteractorDefaultLogger(t *testing.T) {
	m := &loggerMock{t: t}
	mPrometheus := &prometheusMock{}
	mCollector := &mockCollector{}
	mCounter := &mockCounter{}
	mPrometheus.On("NewCounterVector",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string"),
		mock.AnythingOfType("[]string")).Return(mCollector).Once()

	l := MakeFibonacciPrometheusLogger(m, mPrometheus)

	mCollector.On("WithLabelValues", mock.AnythingOfType("[]string")).Return(mCounter)
	mCounter.On("Inc")

	l.LogBadInput(42)
	l.LogRepositoryError(5, 42, nil)
	m.AssertExpectations(t)
	mPrometheus.AssertExpectations(t)
	mCollector.AssertExpectations(t)
	mCounter.AssertExpectations(t)
}
