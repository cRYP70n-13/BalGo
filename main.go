package main

import (
	"balance-ot/config"
	"balance-ot/server"
	"log"
	"os"
	"sync/atomic"
)

func main() {
	file, err := os.Open("./config.yaml")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	conf, err := config.Load(file)
	if err != nil {
		panic(err)
	}
	log.Println(conf)

	backendServers := []server.BackendServer{
		createHealthyBackend("http://localhost:8081"),
		createHealthyBackend("http://localhost:8082"),
		createHealthyBackend("http://localhost:8083"),
		createHealthyBackend("http://localhost:8084"),
	}

	s := server.NewServer("localhost:4000", backendServers)
	s.Start()
}

func createHealthyBackend(address string) server.BackendServer {
	var healthStatus atomic.Value
	healthStatus.Store(true) // Initializing the health status to true
	return server.BackendServer{Address: address, IsHealthy: healthStatus}
}
