package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Yapo/goutils"
)

// HandlerInput is a placeholder for whatever input a handler may need.
type HandlerInput interface{}

// InputGetter defines a type for all functions that, when called, will attempt
// to retrieve and parse the input of a request and return it. Should any error
// happen, a goutils.Response must be filled with an adequate message and code
type InputGetter func() (HandlerInput, *goutils.Response)

// Handler is the interface for the objects that should process web requests.
// Input() must return a fresh struct to be filled with the request input
// Execute(input) receives a filled input struct to handle the request
type Handler interface {
	// Input should return a pointer to the struct that this handler will need
	// to be filled with the user input for a request
	Input(InputRequest) HandlerInput
	// Execute is the actual handler code. The InputGetter can be used to retrieve
	// the request's input at any time (or not at all).
	Execute(InputGetter) *goutils.Response
}

// InputHandler defines what methods an input handler should have
type InputHandler interface {
	NewInputRequest(*http.Request) InputRequest
	SetInputRequest(InputRequest, HandlerInput)
	Input() (HandlerInput, *goutils.Response)
}

// InputRequest defines what methods an input handler should have
type InputRequest interface {
	Set(interface{}) InputRequest
	FromJSONBody() InputRequest
	FromRawBody() InputRequest
	FromPath() InputRequest
	FromQuery() InputRequest
	FromHeaders() InputRequest
}

// MakeJSONHandlerFunc wraps a Handler on a json-over-http context, returning
// a standard http.HandlerFunc
func MakeJSONHandlerFunc(h Handler, l JSONHandlerLogger, ih InputHandler,
	cacheDriver CacheDriver, cacheable bool) http.HandlerFunc {
	jh := jsonHandler{handler: h, logger: l, inputHandler: ih, cacheDriver: cacheDriver, cacheable: cacheable}
	return jh.run
}

// JSONHandlerLogger defines all the events a jsonHandler can report
type JSONHandlerLogger interface {
	LogRequestStart(r *http.Request)
	LogRequestEnd(*http.Request, *goutils.Response)
	LogRequestPanic(*http.Request, *goutils.Response, interface{})
	LogResponseFromCache(*http.Request, *goutils.Response)
	LogErrorSettingCache(r *http.Request, err error)
}

// CacheDriver implements cache control operations
type CacheDriver interface {
	SetCache(input interface{}, data json.RawMessage) error
	GetCache(input interface{}) (json.RawMessage, error)
}

// jsonHandler provides an http.HandlerFunc that reads its input and formats
// its output as json
type jsonHandler struct {
	handler      Handler
	logger       JSONHandlerLogger
	inputHandler InputHandler
	cacheDriver  CacheDriver
	cacheable    bool
}

// run will prepare the input for the actual handler and format the response
// as json. Also, request information will be logged. It's an instance of
// http.HandlerFunc
func (jh *jsonHandler) run(w http.ResponseWriter, r *http.Request) {
	jh.logger.LogRequestStart(r)
	// Default response
	response := &goutils.Response{
		Code: http.StatusInternalServerError,
	}
	// Function the request can call to retrieve its input
	ri := jh.inputHandler.NewInputRequest(r)
	input := jh.handler.Input(ri)
	jh.inputHandler.SetInputRequest(ri, input)
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
	if jh.cacheable {
		cacheKey, _ := jh.inputHandler.Input()
		if cache, e := jh.cacheDriver.GetCache(cacheKey); e == nil {
			if e := json.Unmarshal(cache, response); e == nil {
				jh.logger.LogResponseFromCache(r, response)
				return
			}
		}
	}
	// Do the Harlem Shake
	response = jh.handler.Execute(jh.inputHandler.Input)
	if jh.cacheable {
		if dataRaw, e := json.Marshal(response); e == nil {
			if e := jh.cacheDriver.SetCache(input, dataRaw); e != nil {
				jh.logger.LogErrorSettingCache(r, e)
			}
		}
	}
	jh.logger.LogRequestEnd(r, response)
}
