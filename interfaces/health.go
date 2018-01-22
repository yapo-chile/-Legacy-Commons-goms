package interfaces

import (
	"github.com/Yapo/goutils"
	"github.com/Yapo/logger"
	"net/http"
)

// HealthHandler implements the handler interface and responds to /healthcheck
// requests with an OK message. Expected response format:
// { Status: string - Always set to OK }
type HealthHandler struct{}

// Run retrieves service health status.
func (*HealthHandler) Run(w http.ResponseWriter, r *http.Request) {
	logger.Info("Health Request: [%s] %s from: %s", r.Method, r.URL, r.RemoteAddr)
	var response goutils.Response
	defer goutils.WriteJSONResponse(w, &response)
	defer goutils.CreateJSON(&response)

	response.Body = struct {
		Status string
	}{Status: "OK"}

	response.Code = http.StatusOK
}
