package main

import (
	"os"
	"strconv"
)

const defaultPort = 8080

func main() {
	var s Server
	// Use PORT from environment variables if it's set.
	if portEnv := os.Getenv("PORT"); portEnv != "" {
		port, err := strconv.Atoi(portEnv)
		if err != nil {
			panic("can't convert PORT env into integer")
		}
		s = NewServer(port)
	} else {
		s = NewServer(defaultPort)
	}
	s.Start()
}
