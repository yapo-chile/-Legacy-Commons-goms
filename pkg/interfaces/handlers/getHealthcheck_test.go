package handlers

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/Yapo/goutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockHealthcheckInteractor struct {
	mock.Mock
}

func (m *MockHealthcheckInteractor) GetHealthcheck(ctx context.Context) (string, error) {
	ret := m.Called(ctx)
	return ret.String(0), ret.Error(1)
}

func TestGetHealthcheckOK(t *testing.T) {
	m := MockHealthcheckInteractor{}
	handler := GetHealthcheckHandler{
		GetHealthcheckInteractor: &m,
	}
	var input getHealthcheckHandlerInput
	ctx := mock.AnythingOfType("*context.emptyCtx")
	getter := MakeMockInputGetter(&input, nil)

	m.On("GetHealthcheck", ctx).Return("OK", nil)

	expected := &goutils.Response{
		Code: http.StatusOK,
		Body: getHealthcheckRequestOutput{
			Status: "OK",
		},
	}

	resp := handler.Execute(context.Background(), getter)

	assert.Equal(t, expected, resp)
	m.AssertExpectations(t)
}

func TestGetHealthcheckError(t *testing.T) {
	m := MockHealthcheckInteractor{}
	handler := GetHealthcheckHandler{
		GetHealthcheckInteractor: &m,
	}
	var input getHealthcheckHandlerInput
	ctx := mock.AnythingOfType("*context.emptyCtx")
	getter := MakeMockInputGetter(&input, nil)

	err := errors.New("error")
	m.On("GetHealthcheck", ctx).Return("", err)

	expected := &goutils.Response{
		Code: http.StatusBadRequest,
		Body: getHealthcheckRequestError{err.Error()},
	}

	resp := handler.Execute(context.Background(), getter)

	assert.Equal(t, expected, resp)
	m.AssertExpectations(t)
}
