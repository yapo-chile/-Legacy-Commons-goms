package handlers

import (
	"net/http"

	"github.com/Yapo/goutils"
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
func MakeJSONHandlerFunc(h Handler, l JsonHandlerLogger) http.HandlerFunc {
	jh := jsonHandler{handler: h, logger: l}
	return jh.run
}

// JsonHandlerLogger defines all the events a jsonHandler can report
type JsonHandlerLogger interface {
	LogRequestStart(r *http.Request)
	LogRequestEnd(*http.Request, *goutils.Response)
	LogRequestPanic(*http.Request, *goutils.Response, interface{})
}

// jsonHandler is an http.Handler that reads its input and formats its output
// as json
type jsonHandler struct {
	handler Handler
	logger  JsonHandlerLogger
}

// run will prepare the input for the actual handler and format the response
// as json. Also, request information will be logged.
func (jh *jsonHandler) run(w http.ResponseWriter, r *http.Request) {
	jh.logger.LogRequestStart(r)
	// Default response
	response := &goutils.Response{
		Code: http.StatusInternalServerError,
	}
	// Function the request can call to retrieve its input
	inputGetter := func() (HandlerInput, *goutils.Response) {
		input := jh.handler.Input()
		response := goutils.ParseJSONBody(r, input)
		return input, response
	}
	// Format the output and send it down the writer
	outputWriter := func() {
		goutils.CreateJSON(response)
		goutils.WriteJSONResponse(w, response)
	}
	// Handle panicking handlers and report errors
	errorHandler := func() {
		if err := recover(); err != nil {
			jh.logger.LogRequestPanic(r, response, err)
		}
	}
	// Setup before calling the actual handler
	defer outputWriter()
	defer errorHandler()
	// Do the Harlem Shake
	response = jh.handler.Execute(inputGetter)
	jh.logger.LogRequestEnd(r, response)
}
