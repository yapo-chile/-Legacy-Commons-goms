package main

import (
	"gopkg.in/gorilla/mux.v1"
	"net/http"
)

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type routeGroups struct {
	Prefix string
	Groups []route
}

type routes []routeGroups

// Routes var contains the VERB, PATH and Handler for each route
var Routes = routes{
	{
		"/api/v1",
		[]route{
			{
				"theendpoint",
				"GET",
				"/theendpoint",
				MyGOMSHandler,
			},
		},
	},
}

// NewRouter creates a new instance of a mux route and load rules in config
func NewRouter() *mux.Router {
	router := mux.NewRouter()
	for _, routeGroup := range Routes {
		subRouter := router.PathPrefix(routeGroup.Prefix).Subrouter()
		for _, route := range routeGroup.Groups {
			var handler http.Handler

			handler = route.HandlerFunc
			subRouter.
				Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Handler(handler)
		}
	}
	return router
}
