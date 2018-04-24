package interfaces

import (
	"errors"
	"github.com/Yapo/goutils"
	"github.schibsted.io/Yapo/goms/pkg/domain"
	"gopkg.in/stretchr/testify.v1/assert"
	"gopkg.in/stretchr/testify.v1/mock"
	"net/http"
	"testing"
)

type MockFibonacciInteractor struct {
	mock.Mock
}

func (m *MockFibonacciInteractor) GetNth(n int) (domain.Fibonacci, error) {
	args := m.Called(n)
	return args.Get(0).(domain.Fibonacci), args.Error(1)
}

func TestFibonacciHandlerInput(t *testing.T) {
	m := MockFibonacciInteractor{}
	h := FibonacciHandler{Interactor: &m}
	input := h.Input()
	var expected *fibonacciRequestInput
	assert.IsType(t, expected, input)
	m.AssertExpectations(t)
}

func TestFibonacciHandlerExecuteOK(t *testing.T) {
	m := MockFibonacciInteractor{}
	m.On("GetNth", 5).Return(domain.Fibonacci(5), nil).Once()
	h := FibonacciHandler{Interactor: &m}

	input := fibonacciRequestInput{N: 5}
	expectedResponse := &goutils.Response{
		Code: http.StatusOK,
		Body: fibonacciRequestOutput{5},
	}

	r := h.Execute(&input)
	assert.Equal(t, expectedResponse, r)

	m.AssertExpectations(t)
}

func TestFibonacciHandlerExecuteError(t *testing.T) {
	m := MockFibonacciInteractor{}
	m.On("GetNth", -1).Return(domain.Fibonacci(0), errors.New("kaboom")).Once()
	h := FibonacciHandler{Interactor: &m}

	input := fibonacciRequestInput{N: -1}
	expectedResponse := &goutils.Response{
		Code: http.StatusBadRequest,
		Body: fibonacciRequestError{"kaboom"},
	}

	r := h.Execute(&input)
	assert.Equal(t, expectedResponse, r)

	m.AssertExpectations(t)
}
