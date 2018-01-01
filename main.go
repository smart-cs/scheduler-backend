package main

const port = 8080

func main() {
	s := NewServer(port)
	s.Start()
}
