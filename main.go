package main

import "camellia/server"

func main() {
	s := server.Server{
		Port: 9090,
	}
	s.Start()
}
