package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	mux "gopkg.in/gorilla/mux.v1"

	"github.com/Yapo/goutils"
)

func MakeMockInputGetter(input HandlerInput, response *goutils.Response) InputGetter {
	return func() (HandlerInput, *goutils.Response) {
		return input, response
	}
}

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Input(ir InputRequest) HandlerInput {
	args := m.Called(ir)
	return args.Get(0).(HandlerInput)
}

func (m *MockHandler) Execute(getter InputGetter) *goutils.Response {
	args := m.Called(getter)
	_, response := getter()
	if response != nil {
		return response
	}
	return args.Get(0).(*goutils.Response)
}

type MockInputHandler struct {
	mock.Mock
}

func (m *MockInputHandler) Input() (HandlerInput, *goutils.Response) {
	args := m.Called()
	return args.Get(0).(HandlerInput), args.Get(1).(*goutils.Response)
}

func (m *MockInputHandler) NewInputRequest(r *http.Request) InputRequest {
	args := m.Called(r)
	return args.Get(0).(InputRequest)
}

func (m *MockInputHandler) SetInputRequest(ri InputRequest, hi HandlerInput) {
	m.Called(ri, hi)
}

type mockCacheDriver struct {
	mock.Mock
}

func (m *mockCacheDriver) SetCache(input interface{}, data json.RawMessage) error {
	args := m.Called(input, data)
	return args.Error(0)
}

func (m *mockCacheDriver) GetCache(input interface{}) (json.RawMessage, error) {
	args := m.Called(input)
	return args.Get(0).(json.RawMessage), args.Error(1)
}

type MockPanicHandler struct {
	mock.Mock
}

func (m *MockPanicHandler) Input(ir InputRequest) HandlerInput {
	args := m.Called(ir)
	return args.Get(0).(HandlerInput)
}
func (m *MockPanicHandler) Execute(getter InputGetter) *goutils.Response {
	m.Called(getter)
	panic("dead")
}

type MockLogger struct {
	mock.Mock
}

func (m *MockLogger) LogRequestStart(r *http.Request) {
	m.Called(r)
}
func (m *MockLogger) LogRequestEnd(r *http.Request, response *goutils.Response) {
	m.Called(r, response)
}
func (m *MockLogger) LogRequestPanic(r *http.Request, response *goutils.Response, err interface{}) {
	m.Called(r, response, err)
}
func (m *MockLogger) LogResponseFromCache(r *http.Request, resp *goutils.Response) {
	m.Called(r, resp)
}
func (m *MockLogger) LogErrorSettingCache(r *http.Request, err error) {
	m.Called(r, err)
}

type DummyInput struct {
	X int
}

type DummyOutput struct {
	Y string
}

type TestParam struct {
	Param1 string `get:"param1"`
	Param2 string `get:"param2"`
}

type TestParamInt struct {
	Param1 int `get:"param3"`
	Param2 int `get:"param4"`
}

type TestParamStruct struct {
	Param1 TestParam    `get:"param5"`
	Param2 TestParamInt `get:"param6"`
}

func TestJsonHandlerFuncOK(t *testing.T) {
	h := MockHandler{}
	ih := MockInputHandler{}
	mCache := mockCacheDriver{}
	mMockInputRequest := MockInputRequest{}
	l := MockLogger{}
	input := &DummyInput{}
	response := &goutils.Response{
		Code: 42,
		Body: DummyOutput{"That's some bad hat, Harry"},
	}
	getter := mock.AnythingOfType("handlers.InputGetter")
	h.On("Execute", getter).Return(response).Once()
	h.On("Input", mock.AnythingOfType("*handlers.MockInputRequest")).Return(input).Once()

	ih.On("NewInputRequest", mock.AnythingOfType("*http.Request")).Return(&mMockInputRequest)
	ih.On("Input").Return(input, response)
	ih.On(
		"SetInputRequest",
		mock.AnythingOfType("*handlers.MockInputRequest"),
		mock.AnythingOfType("*handlers.DummyInput"),
	)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/someurl", strings.NewReader("{}"))

	l.On("LogRequestStart", r)
	l.On("LogRequestEnd", r, response)

	fn := MakeJSONHandlerFunc(&h, &l, &ih, &mCache, false)
	fn(w, r)

	assert.Equal(t, 42, w.Code)
	assert.Equal(t, `{"Y":"That's some bad hat, Harry"}`+"\n", w.Body.String())
	h.AssertExpectations(t)
	ih.AssertExpectations(t)
	mMockInputRequest.AssertExpectations(t)
	mCache.AssertExpectations(t)
	l.AssertExpectations(t)
}

func TestJsonHandlerFuncOKWithGetCacheOK(t *testing.T) {
	h := MockHandler{}
	ih := MockInputHandler{}
	mCache := mockCacheDriver{}
	mMockInputRequest := MockInputRequest{}
	l := MockLogger{}
	input := &DummyInput{}
	response := &goutils.Response{
		Code: 42,
		Body: DummyOutput{"That's some bad hat, Harry"},
	}
	h.On("Input", mock.AnythingOfType("*handlers.MockInputRequest")).Return(input)
	var cache json.RawMessage
	cache, _ = json.Marshal(response)
	mCache.On("GetCache", input).Return(cache, nil)
	ih.On("NewInputRequest", mock.AnythingOfType("*http.Request")).Return(&mMockInputRequest)
	ih.On("Input").Return(input, response)
	ih.On(
		"SetInputRequest",
		mock.AnythingOfType("*handlers.MockInputRequest"),
		mock.AnythingOfType("*handlers.DummyInput"),
	)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/someurl", strings.NewReader("{}"))

	l.On("LogRequestStart", r)
	l.On("LogResponseFromCache", mock.AnythingOfType("*http.Request"), mock.AnythingOfType("*goutils.Response"))

	fn := MakeJSONHandlerFunc(&h, &l, &ih, &mCache, true)
	fn(w, r)

	assert.Equal(t, 42, w.Code)
	assert.Equal(t, `{"Y":"That's some bad hat, Harry"}`+"\n", w.Body.String())
	h.AssertExpectations(t)
	ih.AssertExpectations(t)
	mMockInputRequest.AssertExpectations(t)
	mCache.AssertExpectations(t)
	l.AssertExpectations(t)
}

func TestJsonHandlerFuncOKWithSetCacheError(t *testing.T) {
	h := MockHandler{}
	ih := MockInputHandler{}
	mCache := mockCacheDriver{}
	mMockInputRequest := MockInputRequest{}
	l := MockLogger{}
	input := &DummyInput{}
	response := &goutils.Response{
		Code: 42,
		Body: DummyOutput{"That's some bad hat, Harry"},
	}
	getter := mock.AnythingOfType("handlers.InputGetter")
	h.On("Input", mock.AnythingOfType("*handlers.MockInputRequest")).Return(input)
	h.On("Execute", getter).Return(response).Once()
	mCache.On("GetCache", input).Return(json.RawMessage{}, fmt.Errorf("CACHE NOT FOUND"))
	e := fmt.Errorf("error setting cache")
	mCache.On("SetCache", mock.AnythingOfType("*handlers.DummyInput"),
		mock.AnythingOfType("json.RawMessage")).Return(e)
	ih.On("NewInputRequest", mock.AnythingOfType("*http.Request")).Return(&mMockInputRequest)
	ih.On("Input").Return(input, response)
	ih.On(
		"SetInputRequest",
		mock.AnythingOfType("*handlers.MockInputRequest"),
		mock.AnythingOfType("*handlers.DummyInput"),
	)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/someurl", strings.NewReader("{}"))

	l.On("LogRequestStart", r)
	l.On("LogErrorSettingCache", r, e)
	l.On("LogRequestEnd", r, response)

	fn := MakeJSONHandlerFunc(&h, &l, &ih, &mCache, true)
	fn(w, r)

	assert.Equal(t, 42, w.Code)
	assert.Equal(t, `{"Y":"That's some bad hat, Harry"}`+"\n", w.Body.String())
	h.AssertExpectations(t)
	ih.AssertExpectations(t)
	mMockInputRequest.AssertExpectations(t)
	mCache.AssertExpectations(t)
	l.AssertExpectations(t)
}

func TestJsonHandlerFuncOK2(t *testing.T) {
	h := MockHandler{}
	ih := MockInputHandler{}
	mMockInputRequest := MockInputRequest{}
	mCache := mockCacheDriver{}
	l := MockLogger{}
	input := &DummyInput{}
	response := &goutils.Response{
		Code: 42,
		Body: DummyOutput{"That's some bad hat, Harry"},
	}
	getter := mock.AnythingOfType("handlers.InputGetter")
	h.On("Execute", getter).Return(response).Once()
	h.On("Input", mock.AnythingOfType("*handlers.MockInputRequest")).Return(input).Once()

	ih.On("NewInputRequest", mock.AnythingOfType("*http.Request")).Return(&mMockInputRequest)
	ih.On("Input").Return(input, response)
	ih.On(
		"SetInputRequest",
		mock.AnythingOfType("*handlers.MockInputRequest"),
		mock.AnythingOfType("*handlers.DummyInput"),
	)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/someurl?id=1,2", strings.NewReader("{}"))
	r = mux.SetURLVars(r, map[string]string{
		"id": "1, 2",
	})

	l.On("LogRequestStart", r)
	l.On("LogRequestEnd", r, response)

	fn := MakeJSONHandlerFunc(&h, &l, &ih, &mCache, false)
	fn(w, r)

	assert.Equal(t, 42, w.Code)
	assert.Equal(t, `{"Y":"That's some bad hat, Harry"}`+"\n", w.Body.String())
	h.AssertExpectations(t)
	ih.AssertExpectations(t)
	mMockInputRequest.AssertExpectations(t)
	mCache.AssertExpectations(t)
	l.AssertExpectations(t)
}

func TestJsonHandlerFuncParseError(t *testing.T) {
	h := MockHandler{}
	ih := MockInputHandler{}
	mMockInputRequest := MockInputRequest{}
	mCache := mockCacheDriver{}
	l := MockLogger{}
	input := &DummyInput{}
	getter := mock.AnythingOfType("handlers.InputGetter")
	response := &goutils.Response{
		Code: 400,
		Body: struct{ ErrorMessage string }{ErrorMessage: "unexpected EOF"},
	}
	h.On("Execute", getter)
	h.On("Input", mock.AnythingOfType("*handlers.MockInputRequest")).Return(input).Once()

	ih.On("NewInputRequest", mock.AnythingOfType("*http.Request")).Return(&mMockInputRequest)
	ih.On("Input").Return(input, response)
	ih.On(
		"SetInputRequest",
		mock.AnythingOfType("*handlers.MockInputRequest"),
		mock.AnythingOfType("*handlers.DummyInput"),
	)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "/someurl", strings.NewReader("{"))

	l.On("LogRequestStart", r)
	l.On("LogRequestEnd", r, mock.AnythingOfType("*goutils.Response"))

	fn := MakeJSONHandlerFunc(&h, &l, &ih, &mCache, false)
	fn(w, r)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, `{"ErrorMessage":"unexpected EOF"}`+"\n", w.Body.String())
	h.AssertExpectations(t)
	ih.AssertExpectations(t)
	mMockInputRequest.AssertExpectations(t)
	mCache.AssertExpectations(t)
	l.AssertExpectations(t)
}

func TestJsonHandlerFuncPanic(t *testing.T) {
	h := MockPanicHandler{}
	ih := MockInputHandler{}
	mMockInputRequest := MockInputRequest{}
	mCache := mockCacheDriver{}
	l := MockLogger{}
	getter := mock.AnythingOfType("handlers.InputGetter")
	input := &DummyInput{}
	h.On("Execute", getter)
	h.On("Input", mock.AnythingOfType("*handlers.MockInputRequest")).Return(input).Once()

	ih.On("NewInputRequest", mock.AnythingOfType("*http.Request")).Return(&mMockInputRequest)
	ih.On(
		"SetInputRequest",
		mock.AnythingOfType("*handlers.MockInputRequest"),
		mock.AnythingOfType("*handlers.DummyInput"),
	)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/someurl", strings.NewReader("{"))

	l.On("LogRequestStart", r)
	l.On("LogRequestPanic", r, mock.AnythingOfType("*goutils.Response"), "dead")

	fn := MakeJSONHandlerFunc(&h, &l, &ih, &mCache, false)
	fn(w, r)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Equal(t, "null"+"\n", w.Body.String())
	h.AssertExpectations(t)
	ih.AssertExpectations(t)
	mMockInputRequest.AssertExpectations(t)
	mCache.AssertExpectations(t)
	l.AssertExpectations(t)
}
