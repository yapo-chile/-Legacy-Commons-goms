package infrastructure

import (
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/gorilla/mux.v1"
)

func TestQueryParamsOK(t *testing.T) {
	type input struct {
		Id string `query:"id"`
	}

	result := input{}
	expected := input{"1"}
	r := httptest.NewRequest("GET", "/api/v1?id=1", nil)

	inputHandler := NewInputHandler()
	ri := inputHandler.NewInputRequest(r)
	ri.Set(&result).FromQuery()

	inputHandler.SetInputRequest(ri, &result)
	result2, err := inputHandler.Input()
	assert.Nil(t, err)
	assert.Equal(t, &expected, result2)
}

func TestPathOK(t *testing.T) {
	type input struct {
		Id string `path:"id"`
	}

	result := input{}
	expected := input{"1"}
	r := httptest.NewRequest("GET", "/api/v1/1", nil)
	r = mux.SetURLVars(r, map[string]string{
		"id": "1",
	})

	inputHandler := NewInputHandler()
	ri := inputHandler.NewInputRequest(r)
	ri.Set(&result).FromPath()

	inputHandler.SetInputRequest(ri, &result)
	result2, err := inputHandler.Input()
	assert.Nil(t, err)
	assert.Equal(t, &expected, result2)
}

func TestJsonBodyOK(t *testing.T) {
	type input struct {
		Id string `json:"id"`
	}

	result := input{}
	expected := input{"edgar"}
	r := httptest.NewRequest("POST", "/api/v1/", strings.NewReader(`{"id": "edgar"}`))

	inputHandler := NewInputHandler()
	ri := inputHandler.NewInputRequest(r)
	ri.Set(&result).FromJsonBody()

	inputHandler.SetInputRequest(ri, &result)
	result2, err := inputHandler.Input()
	assert.Nil(t, err)
	assert.Equal(t, &expected, result2)
}

func TestHeadersOK(t *testing.T) {
	type input struct {
		Id string `headers:"Id"`
	}

	result := input{}
	expected := input{"edgar"}
	r := httptest.NewRequest("POST", "/api/v1/", nil)
	r.Header.Add("Id", "edgar")

	inputHandler := NewInputHandler()
	ri := inputHandler.NewInputRequest(r)
	ri.Set(&result).FromHeaders()

	inputHandler.SetInputRequest(ri, &result)
	result2, err := inputHandler.Input()
	assert.Nil(t, err)
	assert.Equal(t, &expected, result2)
}

func TestAllSourcesAtSameTimeOK(t *testing.T) {
	type input struct {
		HeaderId string `headers:"Id"`
		BodyId   string `json:"id"`
		PathId   string `path:"id"`
		QueryId  string `query:"id"`
	}

	result := input{}
	expected := input{"edgar", "edgod", "edgugu", "edgarda"}
	r := httptest.NewRequest("POST", "/api/v1/edgugu?id=edgarda", strings.NewReader(`{"id": "edgod"}`))
	r.Header.Add("Id", "edgar")
	r = mux.SetURLVars(r, map[string]string{
		"id": "edgugu",
	})

	inputHandler := NewInputHandler()
	ri := inputHandler.NewInputRequest(r)
	ri.Set(&result).FromHeaders().FromJsonBody().FromPath().FromQuery()

	inputHandler.SetInputRequest(ri, &result)
	result2, err := inputHandler.Input()
	assert.Nil(t, err)
	assert.Equal(t, &expected, result2)
}

func TestOverlapSourcesOK(t *testing.T) {
	type input struct {
		BodyId string `json:"id" path:"id"`
		PathId string `path:"id"`
	}

	result := input{}
	expected := input{"edgugu", "edgugu"}
	r := httptest.NewRequest("POST", "/api/v1/edgugu", strings.NewReader(`{"id": "edgod"}`))
	r = mux.SetURLVars(r, map[string]string{
		"id": "edgugu",
	})

	inputHandler := NewInputHandler()
	ri := inputHandler.NewInputRequest(r)
	ri.Set(&result).FromJsonBody().FromPath()

	inputHandler.SetInputRequest(ri, &result)
	result2, err := inputHandler.Input()
	assert.Nil(t, err)
	assert.Equal(t, &expected, result2)
}

func TestEmptySourceOK(t *testing.T) {
	type input struct {
		BodyId string `json:"id"`
		PathId string `path:"id"`
	}

	result := input{}
	expected := input{BodyId: "edgod"}
	r := httptest.NewRequest("POST", "/api/v1/edgugu", strings.NewReader(`{"id": "edgod"}`))

	inputHandler := NewInputHandler()
	ri := inputHandler.NewInputRequest(r)
	ri.Set(&result).FromJsonBody().FromPath()

	inputHandler.SetInputRequest(ri, &result)
	result2, err := inputHandler.Input()
	assert.Nil(t, err)
	assert.Equal(t, &expected, result2)
}
