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

// TrackHandlerFunc instruments an http.HandlerFunc
func (n *NewRelicHandler) TrackHandlerFunc(pattern string, handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		txn := n.app.StartTransaction(pattern, w, r)
		defer txn.End()
		handler(txn, r)
	}
}
