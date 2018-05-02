package infrastructure

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewrelicStartError(t *testing.T) {
	nr := NewRelicHandler{
		Appname: "Test",
		Key:     "NotAValidKey",
	}
	err := nr.Start()
	assert.Error(t, err)
}

type MockHandler struct {
	mock.Mock
}

func MockHandlerFunc(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("been there"))
}

func TestNewrelicStartOk(t *testing.T) {
	var conf NewRelicConf
	LoadFromEnv(&conf)
	nr := NewRelicHandler{
		Appname: conf.Appname,
		Key:     conf.Key,
	}
	err := nr.Start()
	assert.NoError(t, err)

	m := MockHandlerFunc
	handler := nr.TrackHandlerFunc("test", m)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/someurl", strings.NewReader("{}"))
	handler(w, r)

	assert.Equal(t, "been there", w.Body.String())
}
