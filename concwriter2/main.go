package main

import (
	"log"
	"net"
)

func main() {
	selectedPort := ":8080"
	listener, err := net.Listen("tcp", selectedPort)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	defer listener.Close()

	log.Printf("Server is listening on port %v\n", selectedPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		log.Printf("New connection from %s", conn.RemoteAddr())

		// call server
		go HandleConnection(conn)
	}
}

// go run main.go server.go log.go
// go run client.go
