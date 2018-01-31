package main

import (
	"os"
	"strconv"

	"github.com/nickwu241/schedulecreator-backend/server"
)

const defaultPort = 8080

func main() {
	var s server.Server
	// Use PORT from environment variables if it's set.
	if portEnv := os.Getenv("PORT"); portEnv != "" {
		port, err := strconv.Atoi(portEnv)
		if err != nil {
			panic("can't convert PORT environment is not an integer")
		}
		s = server.NewServer(port)
	} else {
		s = server.NewServer(defaultPort)
	}
	s.Start()
}
