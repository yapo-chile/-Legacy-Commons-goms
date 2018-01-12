package interfaces

import (
	"encoding/json"
	"github.com/Yapo/goutils"
	"github.com/Yapo/logger"
	"github.schibsted.io/Yapo/goms/usecases"
	"net/http"
)

type InjectHandler struct {
	Calculator usecases.Calculator `inject:""`
}

type InjectRequestBody struct {
	Op   string
	A, B int
}

// HealthHandler retrieves service health status.
func (h *InjectHandler) Run(w http.ResponseWriter, r *http.Request) {
	logger.Info("Inject Request: [%s] %s from: %s", r.Method, r.URL, r.RemoteAddr)
	var response goutils.Response
	defer goutils.WriteJSONResponse(w, &response)
	defer goutils.CreateJSON(&response)

	var t InjectRequestBody
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&t)
	if err != nil {
		response.Code = http.StatusBadRequest
		response.Body = struct {
			Error string
		}{Error: err.Error()}
		return
	}
	defer r.Body.Close()

	res, err := h.Calculator.Execute(t.Op, t.A, t.B)

	if err != nil {
		response.Code = http.StatusBadRequest
		response.Body = struct {
			Error string
		}{Error: err.Error()}
		return
	}

	response.Body = struct {
		Result int
	}{Result: res}

	response.Code = http.StatusOK
}
