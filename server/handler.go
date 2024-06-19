package server

import (
	"fmt"
	"net/http"
)

type Handler struct{}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/hello" {
		if _, err := w.Write([]byte("Hello World\n")); err != nil {
			fmt.Println("error when writing response for /hello request")
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusNotFound)
}
