package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("pong"))
		w.WriteHeader(http.StatusOK)
	})

	log.Fatal(http.ListenAndServe(":8081", mux))
}
