package infrastructure

import (
	"net/http"

	"github.com/newrelic/go-agent"
)

// NewRelicHandler struct representing a NewRelic handler with the new relic app
type NewRelicHandler struct {
	Appname string
	Key     string
	app     newrelic.Application
}

// Start initializes the NewRelicHandler
func (n *NewRelicHandler) Start() error {
	config := newrelic.NewConfig(n.Appname, n.Key)
	app, err := newrelic.NewApplication(config)
	if err != nil {
		return err
	}
	n.app = app
	return err
}

// TrackHandler instruments a http.Handler function
func (n *NewRelicHandler) TrackHandler(pattern string, handler http.Handler) http.Handler {
	_, newHandler := newrelic.WrapHandle(n.app, pattern, handler)
	return newHandler
}
