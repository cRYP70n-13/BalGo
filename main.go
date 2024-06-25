package main

import (
	"balance-ot/config"
	"balance-ot/server"
	"log"
	"os"
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

	s := server.NewServer("localhost:4000")
    s.Start()
}
