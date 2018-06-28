package handlers

import (
	"errors"
	"net/http"
	"testing"

	"github.com/Yapo/goutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.schibsted.io/Yapo/goms/pkg/domain"
)

type MockWayStairInteractor struct {
	mock.Mock
}

func (m *MockWayStairInteractor) GetNth(n int) (domain.WayStair, error) {
	args := m.Called(n)
	return args.Get(0).(domain.WayStair), args.Error(1)
}

func TestWayStairHandlerInput(t *testing.T) {
	m := MockWayStairInteractor{}
	h := WayStairHandler{
		Interactor: &m,
	}
	input := h.Input()
	var expected *wayStairRequestInput
	assert.IsType(t, expected, input)
	m.AssertExpectations(t)
}

func TestWayStairHandlerExecuteOK(t *testing.T) {
	m := MockWayStairInteractor{}
	m.On("GetNth", 1).Return(domain.WayStair{Ways: 1, Combs: "{1}"}, nil).Once()
	h := WayStairHandler{Interactor: &m}

	input := wayStairRequestInput{N: 1}
	expectedResponse := &goutils.Response{
		Code: http.StatusOK,
		Body: wayStairRequestOutput{domain.WayStair{Ways: 1, Combs: "{1}"}},
	}

	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)
	assert.Equal(t, expectedResponse, r)

	m.AssertExpectations(t)
}

func TestWayStairHandlerExecuteError(t *testing.T) {
	m := MockWayStairInteractor{}
	m.On("GetNth", -1).Return(domain.WayStair{}, errors.New("kaboom")).Once()
	h := WayStairHandler{Interactor: &m}

	input := wayStairRequestInput{N: -1}
	expectedResponse := &goutils.Response{
		Code: http.StatusBadRequest,
		Body: wayStairRequestError{"kaboom"},
	}

	getter := MakeMockInputGetter(&input, nil)
	r := h.Execute(getter)
	assert.Equal(t, expectedResponse, r)

	m.AssertExpectations(t)
}

func TestWayStairHandlerInputError(t *testing.T) {
	m := MockWayStairInteractor{}
	h := WayStairHandler{Interactor: &m}

	expectedResponse := &goutils.Response{
		Code: http.StatusBadRequest,
		Body: wayStairRequestError{"kaboom"},
	}

	getter := MakeMockInputGetter(nil, expectedResponse)
	r := h.Execute(getter)
	assert.Equal(t, expectedResponse, r)

	m.AssertExpectations(t)
}
