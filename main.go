package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/smart-cs/scheduler-backend/server"
)

func main() {
	s := server.NewServer()
	port, present := os.LookupEnv("PORT")
	if present {
		if _, err := strconv.Atoi(port); err != nil {
			fmt.Printf("PORT environment variable is not an integer: %v\n", err)
			return
		}
	}
	s.Run()
}
