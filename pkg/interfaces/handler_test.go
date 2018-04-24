package interfaces

import (
	"github.com/Yapo/goutils"
	"gopkg.in/stretchr/testify.v1/assert"
	"gopkg.in/stretchr/testify.v1/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Input() HandlerInput {
	args := m.Called()
	return args.Get(0).(HandlerInput)
}

func (m *MockHandler) Execute(input HandlerInput) *goutils.Response {
	args := m.Called(input)
	return args.Get(0).(*goutils.Response)
}

type DummyInput struct {
	X int
}

func TestJsonJandlerFuncOK(t *testing.T) {
	h := MockHandler{}
	input := &DummyInput{}
	response := &goutils.Response{
		Code: 42,
		Body: goutils.GenericError{"That's some bad hat, Harry"},
	}
	h.On("Input").Return(input).Once()
	h.On("Execute", input).Return(response).Once()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/someurl", strings.NewReader("{}"))
	fn := MakeJsonHandlerFunc(&h)
	fn(w, r)

	assert.Equal(t, 42, w.Code)
	assert.Equal(t, `{"ErrorMessage":"That's some bad hat, Harry"}`, w.Body.String())
	h.AssertExpectations(t)
}

func TestJsonJandlerFuncParseError(t *testing.T) {
	h := MockHandler{}
	input := &DummyInput{}
	h.On("Input").Return(input).Once()

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/someurl", strings.NewReader("{"))
	fn := MakeJsonHandlerFunc(&h)
	fn(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, `{"ErrorMessage":"unexpected EOF"}`, w.Body.String())
	h.AssertExpectations(t)
}
