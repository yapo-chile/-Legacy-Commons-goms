package service

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

// Routes var contains the NAME, VERB, PATH and Handler function for each route
var Routes = routes{
	{
		//this is the base path, all routes will start with this
		"/api/v{version:[1-9][0-9]*}",
		[]route{
			{
				"Check service health",
				"GET",
				"/healthcheck",
				HealthHandler,
			},
			{
				"injecttest",
				"GET",
				"/inject",
				MyInjectHandler,
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
			handler := route.HandlerFunc
			subRouter.
				Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Handler(handler)
		}
	}
	return router
}
