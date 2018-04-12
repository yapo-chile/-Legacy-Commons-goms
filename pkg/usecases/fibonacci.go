package usecases

import (
	"fmt"
	"github.schibsted.io/Yapo/goms/pkg/domain"
)

type FibonacciInteractor struct {
	Repository domain.FibonacciRepository
}

func (interactor *FibonacciInteractor) GetNth(n int) (domain.Fibonacci, error) {
	// Ensure correct input
	if n <= 0 {
		return -1, fmt.Errorf("There's no such thing as %dth Fibonacci", n)
	}
	// Check if the repository already knows it
	x, err := interactor.Repository.Get(n)
	if err == nil {
		return x, nil
	}
	// Retrieve the latest pair
	latest := interactor.Repository.LatestPair()
	i, x := latest.Next()
	interactor.Repository.Save(i, x)
	// One step closer. Keep trying
	return interactor.GetNth(n)
}
