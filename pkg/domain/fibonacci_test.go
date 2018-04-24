package domain

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFibonacciPairNext(t *testing.T) {
	p := FibonacciPair{
		IA: 5, A: 5,
		IB: 6, B: 8,
	}
	i, x := p.Next()
	assert.Equal(t, i, 7)
	assert.Equal(t, x, Fibonacci(13))
}
