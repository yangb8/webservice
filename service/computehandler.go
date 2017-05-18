package service

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yangb8/webservice/common/api"
)

// ComputeHandler ...
type ComputeHandler struct {
	api.Handler
}

// NewComputeHandler ...
func NewComputeHandler() http.Handler {
	res := &ComputeHandler{}
	res.Handler.Router = mux.NewRouter()

	res.Methods("GET").Path("/fib").
		Name("get_fib").
		HandlerFunc(fibHandler)

	return res
}

// fibHandler ...
func fibHandler(w http.ResponseWriter, r *http.Request) {
	vals := r.URL.Query()
	if v, ok := vals["n"]; ok {
		if len(v) > 0 {
			if n, err := strconv.Atoi(v[0]); err == nil {
				io.WriteString(w, strconv.Itoa(fib(n)))
			}
		}
	}
}
