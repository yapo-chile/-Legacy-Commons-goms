package infrastructure

import (
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

// Routes is an array of routes with a common prefix
type Routes []routeGroups

// NewRouter setups a Router based on the provided routes.
func NewRouter(routes []routeGroups) *mux.Router {
	router := mux.NewRouter()
	for _, routeGroup := range routes {
		subRouter := router.PathPrefix(routeGroup.Prefix).Subrouter()
		for _, route := range routeGroup.Groups {
			subRouter.
				Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Handler(handlers.MakeJSONHandlerFunc(route.Handler))
		}
	}
	return router
}
