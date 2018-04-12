package interfaces

import (
	"encoding/json"
	"github.com/Yapo/goutils"
	"github.com/Yapo/logger"
	"github.schibsted.io/Yapo/goms/pkg/domain"
	"github.schibsted.io/Yapo/goms/pkg/usecases"
	"net/http"
)

// FibonacciHandler implements the handler interface and responds to
// /fibonacci requests using an interactor. It's purpose is just to
// demonstrate Clean Architecture with a practical scenario
type FibonacciHandler struct {
	Interactor usecases.FibonacciInteractor
}

type fibonacciRequestInput struct {
	N int `json:n`
}

// Run executes an /fibonacci request. Uses the given interactor to carry out
// the operation and get the desired value. Expected body format:
//	{
//		n: int - Number of fibonacci to retrieve (1 based)
//	}
// Expected response format:
//   { Result: int - Operation result }
// Expected error format:
//   { Error: string - Error detail }
func (h *FibonacciHandler) Run(w http.ResponseWriter, r *http.Request) {
	logger.Info("Fibonacci Request: [%s] %s from: %s", r.Method, r.URL, r.RemoteAddr)
	var response goutils.Response
	defer goutils.WriteJSONResponse(w, &response)
	defer goutils.CreateJSON(&response)

	var input fibonacciRequestInput
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		response.Code = http.StatusBadRequest
		response.Body = struct {
			Error string
		}{Error: err.Error()}
		return
	}
	defer r.Body.Close() // nolint: errcheck

	f, err := h.Interactor.GetNth(input.N)

	if err != nil {
		response.Code = http.StatusBadRequest
		response.Body = struct {
			Error string
		}{Error: err.Error()}
		return
	}

	response.Body = struct {
		Result domain.Fibonacci
	}{Result: f}

	response.Code = http.StatusOK
}
