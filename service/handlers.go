package service

import (
	"io"
	"net/http"
	"strconv"
)

func FibHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, strconv.Itoa(fib(10)))
}
