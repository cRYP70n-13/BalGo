package main

import (
	"balance-ot/config"
	"balance-ot/server"
	"os"
)

func main() {
	file, err := os.Open("./config.yaml")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	_, err = config.Load(file)
	if err != nil {
		panic(err)
	}

	s := server.NewServer("localhost:4000")
    s.Start()
}
