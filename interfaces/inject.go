package interfaces

import (
	"encoding/json"
	"github.com/Yapo/goutils"
	"github.com/Yapo/logger"
	"github.schibsted.io/Yapo/goms/usecases"
	"net/http"
)

// InjectHandler implements the handler interface and responds to /inject
// requests using an injected calculator. It's purpose is just to demonstrate
// Dependency Injection with a practical scenario.
type InjectHandler struct {
	Calculator usecases.Calculator `inject:""`
}

type injectRequestBody struct {
	Op   string
	A, B int
}

// Run executes a /inject request. Uses the given calculator to perform
// the operation specified on the body. Expected body format:
//   { Op: string - Operation to perform [add]
//   , A: int - First operand
//   , B: int - Second operand
//   }
// Expected response format:
//   { Result: int - Operation result }
// Expected error format:
//   { Error: string - Error detail }
func (h *InjectHandler) Run(w http.ResponseWriter, r *http.Request) {
	logger.Info("Inject Request: [%s] %s from: %s", r.Method, r.URL, r.RemoteAddr)
	var response goutils.Response
	defer goutils.WriteJSONResponse(w, &response)
	defer goutils.CreateJSON(&response)

	var t injectRequestBody
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
