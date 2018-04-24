package interfaces

import (
	"github.com/Yapo/goutils"
	"gopkg.in/stretchr/testify.v1/assert"
	"net/http"
	"testing"
)

func TestHealthHandlerInput(t *testing.T) {
	var h HealthHandler
	input := h.Input()
	var expected *healthHandlerInput
	assert.IsType(t, expected, input)
}

func TestHealthHandlerRun(t *testing.T) {
	var h HealthHandler
	var input HandlerInput
	r := h.Execute(input)

	expected := &goutils.Response{
		Code: http.StatusOK,
		Body: healthRequestOutput{"OK"},
	}

	assert.Equal(t, expected, r)
}