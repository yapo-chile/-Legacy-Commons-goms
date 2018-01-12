package core

import (
	"gopkg.in/gorilla/mux.v1"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routeGroups struct {
	Prefix string
	Groups []Route
}

type Routes []routeGroups

// Routes var contains the NAME, VERB, PATH and Handler function for each route
// NewRouter creates a new instance of a mux route and load rules in config
func NewRouter(routes []routeGroups) *mux.Router {
	router := mux.NewRouter()
	for _, routeGroup := range routes {
		subRouter := router.PathPrefix(routeGroup.Prefix).Subrouter()
		for _, route := range routeGroup.Groups {
			subRouter.
				Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Handler(route.HandlerFunc)
		}
	}
	return router
}
