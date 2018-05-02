package handlers

import (
	"github.com/Yapo/goutils"
	"github.com/Yapo/logger"
	"net/http"
)

// HandlerInput is a placeholder for whatever input a handler may need.
type HandlerInput interface{}

// InputGetter defines a function that, when called will attempt to retrieve
// the input of a request and return it. Should any error happen, a
// goutils.Response will be filled with an adequate message and error code.
type InputGetter func() (HandlerInput, *goutils.Response)

// Handler is the interface for the objects that should process web requests.
// Input() must return a fresh struct to be filled with the request input
// Execute(input) receives a filled input struct to handle the request
type Handler interface {
	// Input should return a pointer to the struct that this handler will need
	// to be filled with the user input for a request
	Input() HandlerInput
	// Execute is the actual handler code. The InputGetter can be used to retrieve
	// the user input at any time (or not at all).
	Execute(input InputGetter) *goutils.Response
}

// MakeJSONHandlerFunc wraps a handler on a json over http context.
func MakeJSONHandlerFunc(h Handler) http.HandlerFunc {
	jh := jsonHandler{Handler: h}
	return jh.run
}

// jsonHandler is an http.Handler that reads its input and formats its output
// as json
type jsonHandler struct {
	Handler Handler
}

// run will prepare the input for the actual handler and format the response
// as json. Also, request information will be logged.
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
