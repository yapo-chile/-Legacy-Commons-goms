package usecases

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.schibsted.io/Yapo/goms/pkg/domain"
)

type MockWayStairRepository struct {
	mock.Mock
}

func (m MockWayStairRepository) Get(nth int) (domain.WayStair, error) {
	ret := m.Called(nth)
	return ret.Get(0).(domain.WayStair), ret.Error(1)
}

func (m MockWayStairRepository) Save(x domain.WayStair) error {
	ret := m.Called(x)
	return ret.Error(0)
}

type MockWayStairLogger struct {
	mock.Mock
}

func (m MockWayStairLogger) LogBadInput(x int) {
	m.Called(x)
}
func (m MockWayStairLogger) LogRepositoryError(i int, x domain.WayStair, err error) {
	m.Called(i, x, err)
}

func TestWayStairInteractorGetNthKnown(t *testing.T) {
	e := domain.WayStair{Ways: 1, Combs: "{1}", Stair: 1}
	l := &MockWayStairLogger{}
	m := MockWayStairRepository{}
	m.On("Get", 1).Return(e, nil)

	i := WayStairInteractor{
		Logger:     l,
		Repository: &m,
	}

	x, err := i.GetNth(1)
	assert.Equal(t, e, x)
	assert.NoError(t, err)
	m.AssertExpectations(t)
	l.AssertExpectations(t)
}
func TestWayStairInteractorGetNthUnknown(t *testing.T) {
	e := domain.WayStair{Stair: 1, Ways: 1, Combs: "{1}"}
	l := &MockWayStairLogger{}
	m := MockWayStairRepository{}
	l.On("LogRepositoryError", 1, domain.WayStair{}, errors.New("not found"))
	m.On("Get", 1).Return(domain.WayStair{}, fmt.Errorf("not found"))
	m.On("Save", e).Return(nil)
	i := WayStairInteractor{
		Logger:      l,
		Repository:  &m,
		StairsLimit: 15,
	}

	x, err := i.GetNth(1)
	assert.Equal(t, e, x)
	assert.Nil(t, err)
	m.AssertExpectations(t)
	l.AssertExpectations(t)
}
