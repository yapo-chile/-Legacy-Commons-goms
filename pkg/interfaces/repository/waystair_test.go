package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.schibsted.io/Yapo/goms/pkg/domain"
)

func TestWayStairRepositoryCreation(t *testing.T) {
	r := NewWayStair()

	ws, err := r.Calculate(1)
	assert.Equal(t, domain.WayStair{Ways: 1, Combs: "{1}"}, ws)
	assert.NoError(t, err)

	ws, err = r.Calculate(2)
	assert.Equal(t, domain.WayStair{Ways: 2, Combs: "{11,2}"}, ws)
	assert.NoError(t, err)

	ws, err = r.Calculate(-1)
	assert.Equal(t, domain.WayStair{}, ws)
	assert.Error(t, err)
}
