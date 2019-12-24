package handlers

import (
	"github.com/stretchr/testify/mock"
)

// MockInputRequest is a mock class
type MockInputRequest struct {
	mock.Mock
}

// Set is a mocked method
func (m *MockInputRequest) Set(input interface{}) OutputRequest {
	args := m.Called(input)
	return args.Get(0).(OutputRequest)
}

// MockOutputRequest is a mock class
type MockOutputRequest struct {
	mock.Mock
}

// FromJSONBody is a mocked method
func (m *MockOutputRequest) FromJSONBody() OutputRequest {
	m.Called()
	return m
}

// FromRawBody is a mocked method
func (m *MockOutputRequest) FromRawBody() OutputRequest {
	m.Called()
	return m
}

// FromPath is a mocked method
func (m *MockOutputRequest) FromPath() OutputRequest {
	m.Called()
	return m
}

// FromQuery is a mocked method
func (m *MockOutputRequest) FromQuery() OutputRequest {
	m.Called()
	return m
}

// FromHeaders is a mocked method
func (m *MockOutputRequest) FromHeaders() OutputRequest {
	m.Called()
	return m
}

// FromCookies is a mocked method
func (m *MockOutputRequest) FromCookies() OutputRequest {
	m.Called()
	return m
}

// FromForm is a mocked method
func (m *MockOutputRequest) FromForm() OutputRequest {
	m.Called()
	return m
}
