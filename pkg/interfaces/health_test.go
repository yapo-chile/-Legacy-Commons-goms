package interfaces

import (
	"gopkg.in/stretchr/testify.v1/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthHandlerRun(t *testing.T) {
	h := HealthHandler{}
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/healthcheck", nil)
	h.Run(w, req)
	r := w.Result()
	assert.Equal(t, r.StatusCode, http.StatusOK)
	assert.Equal(t, w.Body.String(), `{"Status":"OK"}`)
}
