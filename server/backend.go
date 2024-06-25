package server

import "sync"

var (
	backendServersMutex sync.Mutex
)

// AddBackendServer adds a new backend server
func (s *Server) AddBackendServer(server BackendServer) {
	backendServersMutex.Lock()
	defer backendServersMutex.Unlock()
	s.backendServers = append(s.backendServers, server)
}

// RemoveBackendServer removes a backend server by ID
func (s *Server) RemoveBackendServer(id int) {
	backendServersMutex.Lock()
	defer backendServersMutex.Unlock()
	for i, server := range s.backendServers {
		if server.ID == id {
			s.backendServers = append(s.backendServers[:i], s.backendServers[i+1:]...)
			break
		}
	}
}
