package repository

import (
	"net/http"
	"time"
)

// HTTPRequest interface represents the request that is going to be sent via HTTP
type HTTPRequest interface {
	GetMethod() string
	SetMethod(string) HTTPRequest
	GetPath() string
	SetPath(string) HTTPRequest
	GetBody() interface{}
	SetBody(interface{}) HTTPRequest
	GetHeaders() map[string][]string
	SetHeaders(map[string]string) HTTPRequest
	GetQueryParams() map[string][]string
	SetQueryParams(map[string]string) HTTPRequest
	GetTimeOut() time.Duration
	SetTimeOut(int) HTTPRequest
}

// HTTPResponse struct that contains http response of
type HTTPResponse struct {
	Body    interface{} `json:"Body"`
	Code    int         `json:"Code"`
	Headers http.Header `json:"Headers"`
}

// HTTPHandler implements HTTP handler operations
type HTTPHandler interface {
	Send(HTTPRequest) (interface{}, error)
	NewRequest() HTTPRequest
}
