package server

import (
	"sync/atomic"
)

// GetNextServer returns the next backend server using round-robin scheduling
func (s *Server) GetNextServer() *BackendServer {
	index := atomic.AddInt32(&s.currentIndex, 1) % int32(len(s.backendServers))
	return &s.backendServers[index]
}
