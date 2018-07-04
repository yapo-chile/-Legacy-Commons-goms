package repository

import (
	"fmt"
	"sync"

	"github.schibsted.io/Yapo/goms/pkg/domain"
)

// WayStairRepo instantiate all functions needed so can may be able to
// calculate possible ways and its combinations
type WayStairRepo struct{}

// mapFibonacciRepository is an implementation of domain.FibonacciRepository
// that stores Fibonacci data on a map. It keeps the lastest known pair on a
// separate array to speed up retrieval. The type is intentionally private.
// The correct way to instantiate this type is with NewMapFibonacciRepository.
// This ensures that the required initialization is performed every time.
type storageWayStair struct {
	ways  map[int]domain.Ways
	combs map[int]domain.Combs
	mutex sync.RWMutex
}

// NewMapFibonacciRepository instantiates a fresh mapFibonacciRepository,
// performs the initialization and returns it as a domain.FibonacciRepository.
// The return type prevents others to directly access data members.
func NewStorageWayStair() domain.WayStairRepository {
	var r storageWayStair
	r.ways = map[int]domain.Ways{
		1: 1,
	}
	r.combs = map[int]domain.Combs{
		1: "{1}",
	}
	return &r
}

// Get returns the nth (1 based) Fibonacci should this instance know it.
// Otherwise, will return -1 and error message.
func (r *storageWayStair) Get(nth int) (domain.WayStair, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	response := domain.WayStair{}
	w, foundWays := r.ways[nth]
	c, foundCombs := r.combs[nth]
	if !foundWays || !foundCombs {
		return response, fmt.Errorf("Don't know the %dth WayStair, do you?", nth)
	}
	response.Stair = domain.Stair(nth)
	response.Ways = w
	response.Combs = domain.Combs(c)
	return response, nil
}

// Save sets the nth Fibonacci to x should the last known pair of values end
// at nth-1. Otherwise returns an error message.
func (r *storageWayStair) Save(x domain.WayStair) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if x.Ways < 0 || x.Combs == "" {
		return fmt.Errorf("How do you know this is not possible?")
	}
	r.ways[int(x.Stair)] = x.Ways
	r.combs[int(x.Stair)] = x.Combs
	return nil
}
