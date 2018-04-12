package repository

import (
	"github.schibsted.io/Yapo/goms/pkg/domain"
	"gopkg.in/stretchr/testify.v1/assert"
	"testing"
)

func TestFibonacciRepositoryCreation(t *testing.T) {
	r := NewMapFibonacciRepository()

	x, err := r.Get(1)
	assert.Equal(t, x, domain.Fibonacci(1))
	assert.NoError(t, err)

	x, err = r.Get(2)
	assert.Equal(t, x, domain.Fibonacci(1))
	assert.NoError(t, err)

	x, err = r.Get(3)
	assert.Equal(t, x, domain.Fibonacci(-1))
	assert.Error(t, err)

	p := r.LatestPair()
	expected := domain.FibonacciPair{
		IA: 1, A: domain.Fibonacci(1),
		IB: 2, B: domain.Fibonacci(1),
	}

	assert.Equal(t, p, expected)
}

func TestFibonacciRepositorySaveWildGuess(t *testing.T) {
	r := NewMapFibonacciRepository()

	err := r.Save(5, 32)
	assert.Error(t, err)
}

func TestFibonacciRepositorySaveReplace(t *testing.T) {
	r := NewMapFibonacciRepository()

	err := r.Save(2, 42)
	assert.Error(t, err)
}

func TestFibonacciRepositorySaveNext(t *testing.T) {
	r := NewMapFibonacciRepository()

	err := r.Save(3, 2)
	assert.NoError(t, err)

	p := r.LatestPair()
	expected := domain.FibonacciPair{
		IA: 2, A: domain.Fibonacci(1),
		IB: 3, B: domain.Fibonacci(2),
	}

	assert.Equal(t, p, expected)
}
