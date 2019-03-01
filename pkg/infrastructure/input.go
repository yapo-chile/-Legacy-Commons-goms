package infrastructure

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/Yapo/goutils"
	"gopkg.in/gorilla/mux.v1"

	"github.schibsted.io/Yapo/goms/pkg/interfaces/handlers"
)

const (
	// BODY defines the constant for the body params
	BODY string = "body"
	// RAWBODY defines the constant for process the raw body in bytes
	RAWBODY string = "raw"
	// PATH defines the constant for the path params
	PATH string = "path"
	// QUERY defines the constant for the query params
	QUERY string = "query"
	// HEADERS defines the constant for the headers params
	HEADERS string = "headers"

	// NotSeteable defines the error string of this error
	NotSeteable string = "PROVIDED_INPUT_IS_NOT_SETEABLE"
)

// ErrNotSeteable describes an error that occurs when te var is not seteable
var ErrNotSeteable = errors.New(NotSeteable)

type inputRequest struct {
	output      interface{}
	httpRequest *http.Request
	sources     []string
}

func (ri *inputRequest) Set(output interface{}) handlers.InputRequest {
	ri.output = output
	return ri
}

func (ri *inputRequest) FromJSONBody() handlers.InputRequest {
	ri.sources = append(ri.sources, BODY)
	return ri
}

func (ri *inputRequest) FromRawBody() handlers.InputRequest {
	ri.sources = append(ri.sources, RAWBODY)
	return ri
}

func (ri *inputRequest) FromPath() handlers.InputRequest {
	ri.sources = append(ri.sources, PATH)
	return ri
}

func (ri *inputRequest) FromQuery() handlers.InputRequest {
	ri.sources = append(ri.sources, QUERY)
	return ri
}

func (ri *inputRequest) FromHeaders() handlers.InputRequest {
	ri.sources = append(ri.sources, HEADERS)
	return ri
}

type inputHandler struct {
	inputRequest *inputRequest
	output       handlers.HandlerInput
}

// NewInputHandler returns a new InputHandler
func NewInputHandler() handlers.InputHandler {
	return &inputHandler{}
}

// NewInputRequest returns a new InputRequest based on a http.request
func (ih *inputHandler) NewInputRequest(r *http.Request) handlers.InputRequest {
	return &inputRequest{httpRequest: r}
}

// SetInputRequest sets the input request and the request input
func (ih *inputHandler) SetInputRequest(ri handlers.InputRequest, hi handlers.HandlerInput) {
	ih.inputRequest = ri.(*inputRequest)
	ih.output = hi
}

// Input does the actual process of the input
func (ih *inputHandler) Input() (handlers.HandlerInput, *goutils.Response) {
	if ih.inputRequest.output == nil {
		return ih.output, &goutils.Response{
			Code: http.StatusBadRequest,
			Body: goutils.GenericError{
				ErrorMessage: "Output was not set correctly",
			},
		}
	}

	hasError := false
	reflectedOutput := reflect.ValueOf(ih.inputRequest.output)
	for _, source := range ih.inputRequest.sources {
		switch source {
		case BODY:
			hasError = hasError ||
				goutils.ParseJSONBody(
					ih.inputRequest.httpRequest,
					ih.inputRequest.output,
				) != nil
		case RAWBODY:
			rawBody, _ := ioutil.ReadAll(ih.inputRequest.httpRequest.Body)
			ih.inputRequest.httpRequest.Body = ioutil.NopCloser(bytes.NewBuffer(rawBody))
			hasError = hasError ||
				ih.parseInput(
					map[string]string{"body": string(rawBody)},
					source,
					reflectedOutput,
				) != nil
		case QUERY:
			hasError = hasError ||
				ih.parseInput(
					ih.httpValuesToMap(ih.inputRequest.httpRequest.URL.Query()),
					source,
					reflectedOutput,
				) != nil
		case PATH:
			hasError = hasError ||
				ih.parseInput(
					mux.Vars(ih.inputRequest.httpRequest),
					source,
					reflectedOutput,
				) != nil
		case HEADERS:
			hasError = hasError ||
				ih.parseInput(
					ih.httpValuesToMap(ih.inputRequest.httpRequest.Header),
					source,
					reflectedOutput,
				) != nil
		}
	}

	if hasError {
		return ih.output, &goutils.Response{
			Code: http.StatusBadRequest,
			Body: goutils.GenericError{
				ErrorMessage: "Output was not set correctly",
			},
		}
	}
	return ih.output, nil
}

func (ih *inputHandler) httpValuesToMap(values map[string][]string) map[string]string {
	outValues := make(map[string]string)
	for k, v := range values {
		outValues[k] = strings.Join(v, ",")
	}
	return outValues
}

func (ih *inputHandler) parseInput(vars map[string]string, inputTag string, input reflect.Value) error {
	if input.Kind() == reflect.Ptr {
		reflectedInput := reflect.Indirect(input)
		// We should only keep going if we can set values
		if reflectedInput.IsValid() && reflectedInput.CanSet() {
			if reflectedInput.Kind() == reflect.Struct {
				// Recursively load inner struct fields
				for i := 0; i < reflectedInput.NumField(); i++ {
					if tag, ok := reflectedInput.Type().Field(i).Tag.Lookup(inputTag); ok {
						switch reflectedInput.Field(i).Kind() {
						case reflect.Struct:
							if ih.parseInput(vars, inputTag, reflectedInput.Field(i).Addr()) != nil {
								continue
							}
						case reflect.String:
							reflectedInput.Field(i).SetString(vars[tag])
						case reflect.Int:
							if value, err := strconv.Atoi(vars[tag]); err == nil {
								reflectedInput.Field(i).Set(reflect.ValueOf(value))
							}
						case reflect.Slice:
							values := strings.Split(vars[tag], ",")
							switch reflectedInput.Field(i).Interface().(type) {
							case []string:
								trimValues := []string{}
								for _, value := range values {
									trimValues = append(trimValues, strings.TrimSpace(value))
								}
								reflectedInput.Field(i).Set(reflect.ValueOf(trimValues))
							case []int:
								intValues := []int{}
								for _, value := range values {
									if val, err := strconv.Atoi(strings.TrimSpace(value)); err == nil {
										intValues = append(intValues, val)
									}
								}
								reflectedInput.Field(i).Set(reflect.ValueOf(intValues))
							case []byte:
								reflectedInput.Field(i).SetBytes([]byte(vars[tag]))
							}
						}
					}
				}
			}
		} else {
			return ErrNotSeteable
		}
	}
	return nil
}
