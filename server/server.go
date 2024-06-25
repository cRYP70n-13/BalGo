package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

type BackendServer struct {
	Address   string
	IsHealthy atomic.Value
	ID        int
}

type Server struct {
	Addr           string
	backendServers []BackendServer
	currentIndex   int32
}

func NewServer(Addr string, backendServers []BackendServer) *Server {
	return &Server{
		Addr:           Addr,
		backendServers: backendServers,
	}
}

func (s *Server) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", s.handleLBRequest)

	server := http.Server{
		Addr:         s.Addr,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		Handler:      Middleware{mux: mux},
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT)
	defer cancel()

	log.Println("Server is up and running on port", s.Addr)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen and serve returned err: %v", err)
		}
	}()

	go s.healthChecker()

	<-ctx.Done()
	log.Println("HTTP server shutting down gracefully...")
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("server shutdown returned an err: %v\n", err)
	}
}

func (s *Server) healthChecker() {
	ticker := time.NewTicker(time.Second * 8)
	for range ticker.C {
		fmt.Println("health checking ...")
		go ping(s.backendServers)
	}
}

func (s *Server) handleLBRequest(w http.ResponseWriter, r *http.Request) {
	for i := 0; i < len(s.backendServers); i++ {
		nextServer := s.GetNextServer()
		if !nextServer.IsHealthy.Load().(bool) {
			log.Printf("Skipping unhealthy server: %s\n", nextServer.Address)
			continue
		}

		proxyURL, err := url.Parse(nextServer.Address + "/healthz")
		if err != nil {
			http.Error(w, "Error parsing URL", http.StatusInternalServerError)
			return
		}

		client := http.Client{}
		request := &http.Request{
			Method: r.Method,
			URL:    proxyURL,
			Header: r.Header,
		}

		res, err := client.Do(request)
		if err != nil {
			log.Printf("Error contacting backend server: %v\n", err)
			continue
		}
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)
		if err != nil {
			http.Error(w, "Error reading response", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(res.StatusCode)
		_, _ = w.Write(body)
		return
	}
	http.Error(w, "Service Unavailable", http.StatusServiceUnavailable)
}
