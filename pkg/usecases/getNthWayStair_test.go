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

func (m MockWayStairRepository) WayStairs(nth int) (domain.WayStair, error) {
	ret := m.Called(nth)
	return ret.Get(0).(domain.WayStair), ret.Error(1)
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

func TestWayStairInteractorGetNthNegative(t *testing.T) {
	l := &MockWayStairLogger{}
	m := MockWayStairRepository{}
	i := WayStairInteractor{
		Logger:     l,
		Repository: m,
	}

	l.On("LogBadInput", -1)

	_, err := i.GetNth(-1)
	assert.Error(t, err)
	m.AssertExpectations(t)
	l.AssertExpectations(t)
}

func TestWayStairInteractorGetNthTooHigh(t *testing.T) {
	l := &MockWayStairLogger{}
	m := MockWayStairRepository{}
	i := WayStairInteractor{
		Logger:     l,
		Repository: m,
	}

	l.On("LogBadInput", 15)

	_, err := i.GetNth(15)
	assert.Error(t, err)
	m.AssertExpectations(t)
	l.AssertExpectations(t)
}

func TestWayStairInteractorGetNthKnown(t *testing.T) {
	e := domain.WayStair{Ways: 1, Combs: "{1}"}
	l := &MockWayStairLogger{}
	m := MockWayStairRepository{}
	m.On("WayStairs", 1).Return(e, nil)

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

func TestWayStairInteractorGetNthError(t *testing.T) {
	e := domain.WayStair{}
	l := &MockWayStairLogger{}
	m := MockWayStairRepository{}
	l.On("LogRepositoryError", 1, e, errors.New("anything"))
	m.On("WayStairs", 1).Return(e, fmt.Errorf("anything"))

	i := WayStairInteractor{
		Logger:     l,
		Repository: &m,
	}

	x, err := i.GetNth(1)
	assert.Equal(t, e, x)
	assert.Error(t, err)
	m.AssertExpectations(t)
	l.AssertExpectations(t)
}
