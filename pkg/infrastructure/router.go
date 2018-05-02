package infrastructure

import (
	"net/http"

	"github.schibsted.io/Yapo/goms/pkg/interfaces/handlers"
	"gopkg.in/gorilla/mux.v1"
)

// Route stands for an http endpoint description
type Route struct {
	Name    string
	Method  string
	Pattern string
	Handler handlers.Handler
}

type routeGroups struct {
	Prefix string
	Groups []Route
}

type WrapperFunc func(pattern string, handler http.HandlerFunc) http.HandlerFunc

// Routes is an array of routes with a common prefix
type Routes []routeGroups

type RouterMaker struct {
	Routes      Routes
	WrapperFunc WrapperFunc
}

// NewRouter setups a Router based on the provided routes.
func (maker *RouterMaker) NewRouter() *mux.Router {
	router := mux.NewRouter()
	for _, routeGroup := range maker.Routes {
		subRouter := router.PathPrefix(routeGroup.Prefix).Subrouter()
		for _, route := range routeGroup.Groups {
			handler := handlers.MakeJSONHandlerFunc(route.Handler)
			if maker.WrapperFunc != nil {
				handler = maker.WrapperFunc(route.Name, handler)
			}
			subRouter.
				Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Handler(handler)
		}
	}
	return router
}
