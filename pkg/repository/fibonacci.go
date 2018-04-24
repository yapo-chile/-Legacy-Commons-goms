package repository

import (
	"fmt"
	"github.schibsted.io/Yapo/goms/pkg/domain"
)

type MapFibonacciRepository struct {
	storage map[int]domain.Fibonacci
	latest  []int
}

func NewMapFibonacciRepository() (r MapFibonacciRepository) {
	r.storage = map[int]domain.Fibonacci{
		1: 1,
		2: 1,
	}
	r.latest = []int{1, 2}
	return
}

func (r *MapFibonacciRepository) Get(nth int) (domain.Fibonacci, error) {
	f, found := r.storage[nth]
	if !found {
		return -1, fmt.Errorf("Don't know the %dth Fibonacci, do you?", nth)
	}
	return f, nil
}

func (r *MapFibonacciRepository) Save(nth int, x domain.Fibonacci) error {
	if nth != r.latest[1]+1 {
		return fmt.Errorf("How do you know the %dth Fibonacci number?", nth)
	}
	r.storage[nth] = x
	r.latest[0]++
	r.latest[1]++
	return nil
}

func (r *MapFibonacciRepository) LatestPair() domain.FibonacciPair {
	return domain.FibonacciPair{
		IA: r.latest[0],
		IB: r.latest[1],
		A:  r.storage[r.latest[0]],
		B:  r.storage[r.latest[1]],
	}
}
