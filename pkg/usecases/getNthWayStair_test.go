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

<<<<<<< HEAD
func (m MockWayStairRepository) Get(nth int) (domain.WayStair, error) {
=======
func (m MockWayStairRepository) Calculate(nth int) (domain.WayStair, error) {
>>>>>>> c6b76e5c8a1f3d8b66fc66fcd73e1cfcac7fd2b9
	ret := m.Called(nth)
	return ret.Get(0).(domain.WayStair), ret.Error(1)
}

<<<<<<< HEAD
func (m MockWayStairRepository) Save(x domain.WayStair) error {
	ret := m.Called(x)
	return ret.Error(0)
}

=======
>>>>>>> c6b76e5c8a1f3d8b66fc66fcd73e1cfcac7fd2b9
type MockWayStairLogger struct {
	mock.Mock
}

func (m MockWayStairLogger) LogBadInput(x int) {
	m.Called(x)
}
func (m MockWayStairLogger) LogRepositoryError(i int, x domain.WayStair, err error) {
	m.Called(i, x, err)
}

<<<<<<< HEAD
func (m MockWayStairLogger) LogCalculateError(i int, err error) {
	m.Called(i, err)
}
=======
>>>>>>> c6b76e5c8a1f3d8b66fc66fcd73e1cfcac7fd2b9
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
<<<<<<< HEAD
	m.On("Get", 1).Return(e, nil)
=======
	m.On("Calculate", 1).Return(e, nil)
>>>>>>> c6b76e5c8a1f3d8b66fc66fcd73e1cfcac7fd2b9

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
<<<<<<< HEAD
func TestWayStairInteractorGetNthUnknown(t *testing.T) {
	e := domain.WayStair{}
	l := &MockWayStairLogger{}
	m := MockWayStairRepository{}
	l.On("LogRepositoryError", 1, e, errors.New("anything"))
	m.On("Get", 1).Return(e, fmt.Errorf("not found"))
	m.On("Save", domain.WayStair{Stair: 1, Ways: 1, Combs: "{1}"}, nil)
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
=======
>>>>>>> c6b76e5c8a1f3d8b66fc66fcd73e1cfcac7fd2b9

func TestWayStairInteractorGetNthError(t *testing.T) {
	e := domain.WayStair{}
	l := &MockWayStairLogger{}
	m := MockWayStairRepository{}
	l.On("LogRepositoryError", 1, e, errors.New("anything"))
<<<<<<< HEAD
	m.On("Get", 1).Return(e, fmt.Errorf("anything"))
=======
	m.On("Calculate", 1).Return(e, fmt.Errorf("anything"))
>>>>>>> c6b76e5c8a1f3d8b66fc66fcd73e1cfcac7fd2b9

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
