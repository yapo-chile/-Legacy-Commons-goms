package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.schibsted.io/Yapo/goms/pkg/domain"
)

func TestWayStairRepository(t *testing.T) {
	r := NewStorageWayStair()
	ws, err := r.Get(1)
	assert.Equal(t, domain.WayStair{Ways: 1, Combs: "{1}", Stair: 1}, ws)
	assert.NoError(t, err)

	ws, err = r.Get(0)
	assert.Equal(t, domain.WayStair{}, ws)
	assert.Error(t, err)

	errSave := r.Save(domain.WayStair{Ways: 1, Combs: "{1}", Stair: 1})
	assert.NoError(t, errSave)

	errSave = r.Save(domain.WayStair{Ways: 1, Combs: "", Stair: 1})
	assert.Error(t, errSave)

	errSave = r.Save(domain.WayStair{Ways: 1, Combs: "{1}", Stair: 1})
	assert.NoError(t, errSave)

}
