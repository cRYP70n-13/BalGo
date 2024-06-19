package server

import (
	"log"
	"net/http"
)

type Server struct {
	Addr string
}

func NewServer(Addr string) *Server {
	return &Server{
		Addr: Addr,
	}
}

func (s *Server) Start() {
    mux := http.NewServeMux()
    mux.HandleFunc("/hello", helloHandler)

	server := http.Server{
		Addr:    s.Addr,
		Handler: Middleware{mux: mux},
	}

	log.Fatal(server.ListenAndServe())
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("Hello World\n")); err != nil {
		log.Println("error when writing response for /hello request")
		w.WriteHeader(http.StatusInternalServerError)
	}
}
