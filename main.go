package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/smart-cs/scheduler-backend/server"
)

const defaultPort = 8080

func main() {
	var s server.Server
	// Use PORT from environment variables if it's set.
	if portEnv := os.Getenv("PORT"); portEnv != "" {
		port, err := strconv.Atoi(portEnv)
		if err != nil {
			fmt.Printf("ERROR: can't convert PORT environment is not an integer")
			return
		}
		s = server.NewServer(port)
	} else {
		s = server.NewServer(defaultPort)
	}
	s.Start()
}

// RunMain runs main, here for testing purposes
func RunMain() {
	main()
}
