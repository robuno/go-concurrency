package main

import (
	"fmt"
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

	clientCounter := 0
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}

		clientCounter++
		clientID := fmt.Sprintf("client%d", clientCounter)
		log.Printf("New connection from %s", clientID)

		// call server
		go HandleConnection(conn, clientID)
	}
}


// go run main.go server.go log.go
// go run client.go