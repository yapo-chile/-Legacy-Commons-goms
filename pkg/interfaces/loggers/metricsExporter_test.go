package loggers

import (
	"testing"

	"github.com/stretchr/testify/mock"
)

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

func TestCustomMetricsExporter(t *testing.T) {
	m := &prometheusMock{}
	m.On("IncrementCounter", BadInputError).Once()
	m.On("IncrementCounter", RepositoryError).Once()
	result := MakeCustomMetricsExporter(m)
	result.IncrementBadInputCounter()
	result.IncrementRepositoryErrorCounter()
	m.AssertExpectations(t)
}
