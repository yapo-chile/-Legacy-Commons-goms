package handlers

import (
	"github.com/Yapo/goutils"
	"github.com/Yapo/logger"
	"net/http"
)

// HandlerInput is a placeholder for whatever input a handler may need.
type HandlerInput interface{}
type InputGetter func() (HandlerInput, *goutils.Response)

// Handler is the interface for the objects that should process web requests.
// Input() must return a fresh struct to be filled with the request input
// Execute(input) receives a filled input struct to handle the request
type Handler interface {
	Input() HandlerInput
	Execute(input InputGetter) *goutils.Response
}

// MakeJSONHandlerFunc wraps a handler on a json over http context.
func MakeJSONHandlerFunc(h Handler) http.HandlerFunc {
	jh := jsonHandler{Handler: h}
	return jh.run
}

type jsonHandler struct {
	Handler Handler
}

func (jh *jsonHandler) run(w http.ResponseWriter, r *http.Request) {
	logger.Info("< %s %s %s", r.RemoteAddr, r.Method, r.URL)
	inputGetter := func() (HandlerInput, *goutils.Response) {
		input := jh.Handler.Input()
		response := goutils.ParseJSONBody(r, input)
		return input, response
	}
	response := jh.Handler.Execute(inputGetter)
	logger.Info("> %s %s %s (%d)", r.RemoteAddr, r.Method, r.URL, response.Code)

	goutils.CreateJSON(response)
	goutils.WriteJSONResponse(w, response)
}
