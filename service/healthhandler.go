package service

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/yangb8/webservice/common/api"
)

// HealthHandler ...
type HealthHandler struct {
	api.Handler
}

// NewHealthHandler ...
func NewHealthHandler() http.Handler {
	res := &HealthHandler{}

	res.Handler.Router = mux.NewRouter()

	res.Methods("GET").Path("/_health").
		Name("get_health").
		HandlerFunc(healthHandler)

	return res
}

// healthHandler ...
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}
