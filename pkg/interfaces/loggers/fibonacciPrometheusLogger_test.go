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

func (m *prometheusMock) NewEventsCollector(name, help string) EventsCollector {
	args := m.Called(name, help)
	return args.Get(0).(EventsCollector)
}

func (m *prometheusMock) Close() error {
	args := m.Called()
	return args.Error(0)
}

type mockCollector struct {
	mock.Mock
}

func (m *mockCollector) CollectEvent(eventName string, eventType EventType) {
	m.Called(eventName, eventType)
}

func TestFibonacciInteractorDefaultLogger(t *testing.T) {
	m := &loggerMock{t: t}
	mPrometheus := &prometheusMock{}
	mCollector := &mockCollector{}
	mPrometheus.On("NewEventsCollector",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("string")).Return(mCollector).Once()

	l := MakeFibonacciPrometheusLogger(m, mPrometheus)

	mCollector.On("CollectEvent",
		mock.AnythingOfType("string"),
		mock.AnythingOfType("EventType"))

	l.LogBadInput(42)
	l.LogRepositoryError(5, 42, nil)
	m.AssertExpectations(t)
	mPrometheus.AssertExpectations(t)
	mCollector.AssertExpectations(t)
}
