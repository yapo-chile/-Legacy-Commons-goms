package interfaces

import (
	"net/http"
)

type Handler interface {
	Run(w http.ResponseWriter, r *http.Request)
}
