package infrastructure

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLoggerRepository simulate Logger Repo
type MockLoggerInfrastructure struct {
	mock.Mock
}

// Info simulate Info
func (m *MockLoggerInfrastructure) Info(message string, params ...interface{}) {
	m.Called()
}

// Debug simulate Debug Logger
func (m *MockLoggerInfrastructure) Debug(message string, params ...interface{}) {
	m.Called()
}

// Crit simulate Crit Logger
func (m *MockLoggerInfrastructure) Crit(message string, params ...interface{}) {
	m.Called()
}

// Error simulate Error Logger
func (m *MockLoggerInfrastructure) Error(message string, params ...interface{}) {
	m.Called()
}

type expectedResponse struct {
	Status int
	Body   []byte
}

func testingServer(exp expectedResponse) *httptest.Server {
	// generate a test server so we can capture and inspect the request
	testServer := httptest.NewServer(
		http.HandlerFunc(
			func(res http.ResponseWriter, req *http.Request) {
				res.WriteHeader(exp.Status)
				res.Write(exp.Body) // nolint
			},
		),
	)
	return testServer
}

func TestRconfSuccess(t *testing.T) {
	mlogger := MockLoggerInfrastructure{}
	mlogger.On("Info").Return(nil)

	content := &EtcdContent{
		Action: "get",
		Node: EtcdNode{
			Key:           "/payment-schedule/language/es.json",
			Value:         "{ \"ErrorExistsSubs\": \"Subscripción no existe en la base de datos\"}",
			ModifiedIndex: 51,
			CreatedIndex:  51,
		},
	}
	bodyJSON, _ := json.Marshal(content)
	t.Log(string(bodyJSON))

	testServer := testingServer(
		expectedResponse{
			Status: 200,
			Body:   bodyJSON,
		},
	)
	defer testServer.Close()

	createRconf, err := NewRconf(testServer.URL, "/payment-schedule/language/es.json", "/v2/keys", &mlogger)
	if err != nil {
		t.Logf("Error %s", err.Error())
	}

	assert.NoError(t, err)
	assert.Equal(t, content, createRconf.Content)
	mlogger.AssertExpectations(t)
}

func TestRconfErrToGet(t *testing.T) {
	mlogger := MockLoggerInfrastructure{}
	mlogger.On("Error").Return(nil)

	content := &EtcdContent{}
	createRconf, err := NewRconf(
		"http://www.thisdirdoesnotexist.com",
		"/payment-schedule/language/es.json",
		"/v2/keys",
		&mlogger,
	)
	if err != nil {
		t.Logf("Error %s", err.Error())
	}

	expectedError := "Get " +
		"http://www.thisdirdoesnotexist.com/v2/keys/payment-schedule/language/es.json:" +
		" dial tcp: lookup www.thisdirdoesnotexist.com"
	assert.Error(t, err)
	assert.Equal(t, content, createRconf.Content)
	assert.Contains(t, err.Error(), expectedError)
	mlogger.AssertExpectations(t)
}

func TestRconfStatusNotFoundWithBody(t *testing.T) {
	mlogger := MockLoggerInfrastructure{}
	mlogger.On("Error").Return(nil)

	content := &EtcdContent{
		Action: "get",
		Node: EtcdNode{
			Key:           "/payment-schedule/language/es.json",
			Value:         "{}",
			ModifiedIndex: 51,
			CreatedIndex:  51,
		},
	}
	bodyJSON, _ := json.Marshal(content)
	t.Log(string(bodyJSON))
	testServer := testingServer(
		expectedResponse{
			Status: 404,
			Body:   bodyJSON,
		},
	)
	defer testServer.Close()

	createRconf, err := NewRconf(
		testServer.URL,
		"/payment-schedule/language/es.json",
		"/v2/keys",
		&mlogger,
	)
	if err != nil {
		t.Logf("Error %s", err.Error())
	}

	assert.Error(t, err)
	assert.Equal(t, &EtcdContent{}, createRconf.Content)
	assert.EqualError(t, err, "response code invalid: 404 Not Found")
	mlogger.AssertExpectations(t)
}
func TestRconfWrongJson(t *testing.T) {
	mlogger := MockLoggerInfrastructure{}
	mlogger.On("Error").Return(nil)

	testServer := testingServer(
		expectedResponse{
			Status: 200,
			Body:   []byte("{\"action\":what?}"),
		},
	)
	defer testServer.Close()

	createRconf, err := NewRconf(
		testServer.URL,
		"/payment-schedule/language/es.json",
		"/v2/keys",
		&mlogger,
	)
	if err != nil {
		t.Logf("Error %s", err.Error())
	}

	assert.Error(t, err)
	assert.Equal(t, &EtcdContent{}, createRconf.Content)
	assert.EqualError(t, err, "invalid character 'w' looking for beginning of value")
	mlogger.AssertExpectations(t)
}

type errReader int

func (errReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("test error")
}
func (errReader) Close() error {
	return nil
}
func TestReadAllErrOnStatusNotFound(t *testing.T) {
	r := &http.Response{
		StatusCode: http.StatusNotFound,
		Status:     "404 Not Found",
		Body:       errReader(0),
	}
	mlogger := MockLoggerInfrastructure{}
	mlogger.On("Error").Return(nil)

	response, err := readAll(r, &mlogger)

	assert.Error(t, err)
	assert.EqualError(t, err, "response code invalid: 404 Not Found")
	assert.Nil(t, response)
	mlogger.AssertExpectations(t)
}
func TestReadAllErrOnStatusOk(t *testing.T) {
	r := &http.Response{
		StatusCode: http.StatusOK,
		Body:       errReader(0),
	}
	mlogger := MockLoggerInfrastructure{}
	mlogger.On("Error").Return(nil)

	response, err := readAll(r, &mlogger)

	assert.Error(t, err)
	assert.Nil(t, response)
	mlogger.AssertExpectations(t)
}

func TestGetSuccess(t *testing.T) {
	mlogger := MockLoggerInfrastructure{}
	rconf := Rconf{
		Log: &mlogger,
		Content: &EtcdContent{
			Action: "get",
			Node: EtcdNode{
				Key:   "my-key",
				Value: "{\"super_json\": \"my-value\"}",
				IsDir: false,
			},
		},
	}

	callToGet := rconf.Get("super_json")
	t.Logf("Result %s", callToGet)
	assert.Equal(t, "my-value", callToGet)

	callToGet = rconf.Get("key_not_exist")
	assert.Equal(t, "", callToGet)
	mlogger.AssertExpectations(t)
}

func TestEmptyRconf(t *testing.T) {
	mlogger := MockLoggerInfrastructure{}
	mlogger.On("Error").Return(nil)
	rconf := Rconf{
		Log:     &mlogger,
		Content: nil,
	}

	callToGet := rconf.Get("super_json")
	t.Logf("Result %s", callToGet)
	assert.Equal(t, "", callToGet)

	mlogger.AssertExpectations(t)
}

func TestGetIsDir(t *testing.T) {
	mlogger := MockLoggerInfrastructure{}
	mlogger.On("Error").Return(nil)
	rconf := Rconf{
		Log: &mlogger,
		Content: &EtcdContent{
			Action: "get",
			Node: EtcdNode{
				Key:   "my-key",
				Value: "",
				IsDir: true,
			},
		},
	}

	callToGet := rconf.Get("super_json")
	assert.Equal(t, "", callToGet)
	mlogger.AssertExpectations(t)
}
func TestTranslateSuccess(t *testing.T) {
	mlogger := MockLoggerInfrastructure{}
	mlogger.On("Debug").Return(nil)
	rconf := Rconf{
		Log: &mlogger,
		Content: &EtcdContent{
			Action: "get",
			Node: EtcdNode{
				Key:   "my-key",
				Value: "{\"super_json\": \"my-value\"}",
				IsDir: false,
			},
		},
	}

	translate := rconf.Translate("super_json")
	assert.Equal(t, "my-value", translate)
	mlogger.AssertExpectations(t)
}
