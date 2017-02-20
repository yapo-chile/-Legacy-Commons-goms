package main

import (
	"encoding/json"
	"fmt"
	"github.com/Yapo/logger"
	"net/http"
)

// Response is a struct to generate a response from POST/PUT requests
type Response struct {
	Code int
	Body interface{}
}

// WriteHTTPResponse write to te response stream
func WriteHTTPResponse(w http.ResponseWriter, response *Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.Code)
	fmt.Fprintf(w, "%s", response.Body)
}

// CreateJSON convert Body to json format
func CreateJSON(response *Response) {
	jsonResponse, err := json.Marshal(response.Body)
	if err != nil {
		logger.Info("CAN'T ENCODE \"%+v\" TO JSON", response.Body)
		response.Body = ""
		response.Code = http.StatusInternalServerError
		return
	}
	response.Body = jsonResponse
}

// EditCreditsHandler is the handler for the PUT /api/v1/edit endpoint
func MyGOMSHandler(w http.ResponseWriter, r *http.Request) {
	logger.Info("EditCredits [%s] %s from: %s", r.Method, r.URL, r.RemoteAddr)
	var response Response
	defer WriteHTTPResponse(w, &response)
	defer CreateJSON(&response)

	response.Body = struct {
		Response string
	}{Response: "HOLA MUNDO!"}
	response.Code = http.StatusOK
}
