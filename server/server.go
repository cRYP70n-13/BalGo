package server

import (
	"log"
	"net/http"
	"os"
	"syscall"
	"time"
)

type Server struct {
	Addr string
}

func NewServer(Addr string) *Server {
	return &Server{
		Addr: Addr,
	}
}

// interruptSignal are the signals that we rely on to gracefully shut down the system.
var interruptSignal = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func (s *Server) Start() {
	mux := http.NewServeMux()
    handler := NewHandler()
	mux.HandleFunc("/hello", handler.ServeHTTP)

	server := http.Server{
		Addr:         s.Addr,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		Handler:      Middleware{mux: mux},
	}

	log.Fatal(server.ListenAndServe())
}
