package usecases

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.schibsted.io/Yapo/goms/pkg/domain"
	"testing"
)

type MockFibonacciRepository struct {
	mock.Mock
}

func (m MockFibonacciRepository) Get(nth int) (domain.Fibonacci, error) {
	ret := m.Called(nth)
	return ret.Get(0).(domain.Fibonacci), ret.Error(1)
}
func (m MockFibonacciRepository) Save(nth int, x domain.Fibonacci) error {
	ret := m.Called(nth, x)
	return ret.Error(0)
}
func (m MockFibonacciRepository) LatestPair() domain.FibonacciPair {
	ret := m.Called()
	return ret.Get(0).(domain.FibonacciPair)
}

func TestFibonacciInteractorGetNthNegative(t *testing.T) {
	m := MockFibonacciRepository{}
	i := FibonacciInteractor{
		Repository: m,
	}

	_, err := i.GetNth(-1)
	assert.Error(t, err)
	m.AssertExpectations(t)
}

func TestFibonacciInteractorGetNthKnown(t *testing.T) {
	m := MockFibonacciRepository{}
	m.On("Get", 1).Return(domain.Fibonacci(1), nil)

	i := FibonacciInteractor{
		Repository: m,
	}

	x, err := i.GetNth(1)
	assert.Equal(t, x, domain.Fibonacci(1))
	assert.NoError(t, err)
	m.AssertExpectations(t)
}

func TestFibonacciInteractorGetNthUnknown(t *testing.T) {
	m := MockFibonacciRepository{}
	m.On("Get", 4).Return(domain.Fibonacci(-1), errors.New("Some error")).Once()
	m.On("LatestPair").Return(domain.FibonacciPair{IA: 1, A: domain.Fibonacci(1), IB: 2, B: domain.Fibonacci(1)}).Once()
	m.On("Save", 3, domain.Fibonacci(2)).Return(nil)

	m.On("Get", 4).Return(domain.Fibonacci(-1), errors.New("Some error")).Once()
	m.On("LatestPair").Return(domain.FibonacciPair{IA: 2, A: domain.Fibonacci(1), IB: 3, B: domain.Fibonacci(2)}).Once()
	m.On("Save", 4, domain.Fibonacci(3)).Return(nil)

	m.On("Get", 4).Return(domain.Fibonacci(3), nil)

	i := FibonacciInteractor{
		Repository: m,
	}

	x, err := i.GetNth(4)
	assert.Equal(t, x, domain.Fibonacci(3))
	assert.NoError(t, err)
	m.AssertExpectations(t)
}
