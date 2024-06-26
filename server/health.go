package server

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// HealthCheckInterval is the interval for health checking
const HealthCheckInterval = 2 * time.Second

var (
	healthCheckTicker *time.Ticker
	healthCheckMutex  sync.Mutex
)

// StartHealthChecker starts the health checker for backend servers
func (s *Server) StartHealthChecker() {
	healthCheckMutex.Lock()
	defer healthCheckMutex.Unlock()

	if healthCheckTicker != nil {
		return // Already running
	}

	healthCheckTicker = time.NewTicker(HealthCheckInterval)
	go func() {
		for range healthCheckTicker.C {
			ping(s.backendServers)
		}
	}()
}

// StopHealthChecker stops the health checker for backend servers
func (s *Server) StopHealthChecker() {
	healthCheckMutex.Lock()
	defer healthCheckMutex.Unlock()

	if healthCheckTicker != nil {
		healthCheckTicker.Stop()
		healthCheckTicker = nil
	}
}

// ping checks the health of each backend server
func ping(servers []BackendServer) {
	for _, server := range servers {
		res, err := http.Get(server.Address)
		if err != nil || res.StatusCode != http.StatusOK {
			fmt.Printf("Server %s is unhealthy\n", server.Address)
			server.IsHealthy.Store(false)
		} else {
			server.IsHealthy.Store(true)
		}
	}
}
