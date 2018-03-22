package interfaces

import (
	"net/http"
)

// Handler is the interface for the objects that should process
// web requests and use the ResponseWriter as output
type Handler interface {
	Run(w http.ResponseWriter, r *http.Request)
}
