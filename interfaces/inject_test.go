package interfaces

import (
	"errors"
	"fmt"
	"gopkg.in/stretchr/testify.v1/assert"
	"gopkg.in/stretchr/testify.v1/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockCalculator struct {
	mock.Mock
}

func (m *MockCalculator) Execute(op string, a, b int) (int, error) {
	args := m.Called(op, a, b)
	return args.Int(0), args.Error(1)
}

func TestInjectHandlerRunNoBody(t *testing.T) {
	m := MockCalculator{}
	h := InjectHandler{Calculator: &m}

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/inject", nil)
	h.Run(w, req)
	r := w.Result()

	assert.Equal(t, r.StatusCode, http.StatusBadRequest)
	assert.Equal(t, w.Body.String(), `{"Error":"EOF"}`)
	m.AssertExpectations(t)
}

func TestInjectHandlerRunBadBody(t *testing.T) {
	m := MockCalculator{}
	m.On("Execute", "", 0, 0).Return(-1, errors.New("Some error"))
	h := InjectHandler{Calculator: &m}

	w := httptest.NewRecorder()
	body := strings.NewReader(`{}`)
	req := httptest.NewRequest("GET", "/inject", body)
	h.Run(w, req)
	r := w.Result()

	assert.Equal(t, r.StatusCode, http.StatusBadRequest)
	assert.Equal(t, w.Body.String(), `{"Error":"Some error"}`)
	m.AssertExpectations(t)
}

func TestInjectHandlerRunOK(t *testing.T) {
	m := MockCalculator{}
	m.On("Execute", "add", 5, 7).Return(42, nil)
	h := InjectHandler{Calculator: &m}

	w := httptest.NewRecorder()
	body := strings.NewReader(`{"Op": "add", "a": 5, "b": 7}`)
	req := httptest.NewRequest("GET", "/inject", body)
	h.Run(w, req)
	r := w.Result()

	assert.Equal(t, r.StatusCode, http.StatusOK)
	assert.Equal(t, w.Body.String(), fmt.Sprintf(`{"Result":%d}`, 42))
	m.AssertExpectations(t)
}
