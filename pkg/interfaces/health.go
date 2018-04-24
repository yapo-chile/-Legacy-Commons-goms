package interfaces

import (
	"github.com/Yapo/goutils"
	"net/http"
)

// HealthHandler implements the handler interface and responds to /healthcheck
// requests with an OK message. Expected response format:
// { Status: string - Always set to OK }
type HealthHandler struct{}

type healthHandlerInput struct{}
type healthRequestOutput struct {
	Status string `json:"status"`
}

// Input returns a fresh, empty instance of healthHandlerInput
func (*HealthHandler) Input() HandlerInput {
	return &healthHandlerInput{}
}

// Execute returns the service health status.
// Expected response format:
//   { Status: string - Always "OK" }
func (*HealthHandler) Execute(input HandlerInput) *goutils.Response {
	return &goutils.Response{
		Code: http.StatusOK,
		Body: healthRequestOutput{
			Status: "OK",
		},
	}
}
