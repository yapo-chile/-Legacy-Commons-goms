package infrastructure

import (
	"github.schibsted.io/Yapo/goms/pkg/interfaces/loggers"
	"github.schibsted.io/Yapo/goms/pkg/interfaces/repository"
)

// HTTPCircuitBreakerHandler struct to implements http repository operations with circuit breaker
type HTTPCircuitBreakerHandler struct {
	circuitBreaker CircuitBreaker
	logger         loggers.Logger
	httpHandler    repository.HTTPHandler
}

// NewHTTPCircuitBreakerHandler will create a new instance of a custom http request handler
func NewHTTPCircuitBreakerHandler(circuitBreaker CircuitBreaker, logger loggers.Logger, h repository.HTTPHandler) repository.HTTPHandler {
	return &HTTPCircuitBreakerHandler{
		circuitBreaker: circuitBreaker,
		logger:         logger,
		httpHandler:    h,
	}
}

// Send will execute the sending of a http request
// but in this case it will wait until it obtains a succesful response
// in order to continue it's execution
func (h *HTTPCircuitBreakerHandler) Send(req repository.HTTPRequest) (interface{}, error) {
	h.logger.Debug("HTTP - %s - Sending HTTP with circuit breaker request to: %+v", req.GetMethod(), req.GetPath())

	var response interface{}
	var err error
	// do-while: try once or retry until circuit breaker closes
	for ok := true; ok; ok = (err == ErrOpenState || err == ErrTooManyRequests) {
		response, err = h.circuitBreaker.Execute(func() (interface{}, error) {
			return h.httpHandler.Send(req)
		})
	}
	return response, err
}

// NewRequest returns an initialized struct that can be used to make a http request
func (h *HTTPCircuitBreakerHandler) NewRequest() repository.HTTPRequest {
	return h.httpHandler.NewRequest()
}
