package repository

import (
	"fmt"
	"sync"

	"github.schibsted.io/Yapo/goms/pkg/domain"
)

// WayStairRepo instantiate all functions needed so can may be able to
// calculate possible ways and its combinations
type WayStairRepo struct{}

// StorageWayStair is an implementation of domain.WayStair
// that stores the result data of a nth on a map indexing the requested valid nth
// as index in a separate array to speed up retrieval.
type StorageWayStair struct {
	ways  map[int]domain.Ways
	combs map[int]domain.Combs
	mutex sync.RWMutex
}

// NewStorageWayStair instantiates a fresh StorageWayStair,
// performs the initialization and returns it as a domain.WayStair.
// The return type prevents others to directly access data members.
func NewStorageWayStair() domain.WayStairRepository {
	var r StorageWayStair
	r.ways = map[int]domain.Ways{
		1: 1,
	}
	r.combs = map[int]domain.Combs{
		1: "{1}",
	}
	return &r
}

// Get returns the nth index mapped WayStair.
// Otherwise, will return an empty domain.WayStair and error message.
func (r *StorageWayStair) Get(nth int) (domain.WayStair, error) {
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

// Save sets the nth WayStair to the mapped WayStair.
// Otherwise returns an error message.
func (r *StorageWayStair) Save(x domain.WayStair) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if x.Ways < 0 || x.Combs == "" {
		return fmt.Errorf("How do you know this is not possible?")
	}
	r.ways[int(x.Stair)] = x.Ways
	r.combs[int(x.Stair)] = x.Combs
	return nil
}
