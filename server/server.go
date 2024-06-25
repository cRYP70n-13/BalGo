package server

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

type BackendServer struct {
	address   string
	isHealthy atomic.Value
	ID        int
}

type Server struct {
	Addr           string
	backendServers []BackendServer
	currentIndex   int32
}

// TODO: Remove the hard coded list from here and get it from the config.
func NewServer(Addr string) *Server {
	return &Server{
		Addr: Addr,
		backendServers: []BackendServer{
			createHealthyBackend("http://localhost:8081"),
			createHealthyBackend("http://localhost:8082"),
			createHealthyBackend("http://localhost:8083"),
			createHealthyBackend("http://localhost:8084"),
		},
	}
}

func createHealthyBackend(address string) BackendServer {
	var healthStatus atomic.Value
	healthStatus.Store(true) // Initializing the health status to true
	return BackendServer{
		address:   address,
		isHealthy: healthStatus,
	}
}

var interruptSignal = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
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

	ctx, cancel := signal.NotifyContext(context.Background(), interruptSignal...)
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
	ticker := time.NewTicker(time.Second * 2)
	for range ticker.C {
		fmt.Println("healh checking ...")
		go ping(s.backendServers)
	}
}

func ping(servers []BackendServer) {
	for _, server := range servers {
		res, err := http.Get(server.address)
		if err != nil || res.StatusCode != http.StatusOK {
			fmt.Printf("This fucking server: %s is fucking unhealthy\n", server.address)
			server.isHealthy.Store(false)
		} else {
			server.isHealthy.Store(true)
		}
	}
}

// FIXME: In case something went down we have to go to the next server and in case all of them are down we have to return an error
// Beause atm we have a bug here
func (s *Server) handleLBRequest(w http.ResponseWriter, r *http.Request) {
    logRequest(r)
	nextUrl := s.GetNextServer().address
	url, err := url.Parse(nextUrl+"/healthz")
	if err != nil {
		panic(err)
	}

	client := http.Client{}
	request := &http.Request{
		Method: r.Method,
		URL:    url,
	}
	res, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println("******** RESPONSE *********", string(body))
	w.WriteHeader(res.StatusCode)
	_, _ = w.Write(body)
}

// GetNextServer returns the next backend server using round-robin scheduling
func (s *Server) GetNextServer() *BackendServer {
	index := atomic.AddInt32(&s.currentIndex, 1) % int32(len(s.backendServers))
	return &s.backendServers[index]
}

func logRequest(r *http.Request) {
	slog.Info("Request",
		"origin", r.RemoteAddr,
		"method", r.Method,
		"url", r.URL,
		"host", r.Host,
		"user-agent", r.UserAgent(),
	)
}
