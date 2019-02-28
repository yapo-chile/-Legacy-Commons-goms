package handlers

import (
	"github.com/stretchr/testify/mock"
)

type MockInputRequest struct {
	mock.Mock
}

func (m *MockInputRequest) Set(input interface{}) InputRequest {
	args := m.Called(input)
	return args.Get(0).(InputRequest)
}

func (m *MockInputRequest) FromJsonBody() InputRequest {
	args := m.Called()
	return args.Get(0).(InputRequest)
}

func (m *MockInputRequest) FromPath() InputRequest {
	args := m.Called()
	return args.Get(0).(InputRequest)
}

func (m *MockInputRequest) FromQuery() InputRequest {
	args := m.Called()
	return args.Get(0).(InputRequest)
}

func (m *MockInputRequest) FromHeaders() InputRequest {
	args := m.Called()
	return args.Get(0).(InputRequest)
}
