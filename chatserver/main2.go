package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"strings"
)

type Client struct {
	conn net.Conn
	name string
}

var (
	clients   = make(map[net.Conn]Client)
	broadcast = make(chan string) 												// broadcast channel
	mutex     sync.Mutex      												    // handle concurrent access to client map
)


// goroutine = broadcasting messages to all clients
func handleBroadcast() {
    for {																		// always wait for new messages in broadcast channel
        msg := <-broadcast														// read message from broadcast and assign as msg

																				// access client map to read all available clients
        mutex.Lock()															// lock client map = shared resource
        for _, client := range clients {
																				// send message to the client(s)
																				// convert str to byt slices, client.conn = tcp conn
            client.conn.Write([]byte("\r" + strings.Repeat(" ", 100) + "\r"))	// clear current input line
            client.conn.Write([]byte(msg + "\n"))								// Write the broadcasted message
            client.conn.Write([]byte("> "))										// Reprint the prompt
        }
        mutex.Unlock()
    }
}


// goroutine = handle incoming connections
func handleConnection(conn net.Conn) {
	defer conn.Close()

	// create a new client
	conn.Write([]byte("Enter your name: "))							// send msg to client
	name, _ := bufio.NewReader(conn).ReadString('\n')				// read client's input, server reads the data
	name = strings.TrimSpace(name)									// remove the newline
	client := Client{conn: conn, name: name}

	// add client to the map
	mutex.Lock()
	clients[conn] = client // socket - name pair
	mutex.Unlock()

	broadcast <- fmt.Sprintf("%s has joined the chat", client.name)

	// always listen messages from the client
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')		// read client message
		if err != nil {
            broadcast <- fmt.Sprintf("%s has left the chat", client.name)
            mutex.Lock()
            delete(clients, conn)
            mutex.Unlock()
            conn.Close()
            return
        }
		message = strings.TrimSpace(message)						// remove the newline
		if len(message) > 0 {										// Only broadcast non-empty messages
            broadcast <- fmt.Sprintf("%s: %s", client.name, message)
        }
	}
}

func main() {

	listener, err := net.Listen("tcp", ":8080")				// listen tcp conn's on port 8080
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
	defer listener.Close()									// release network port
	fmt.Println("Server started on: 8080")

	// always runs, reads message from broadcast and sends to all clients
	go handleBroadcast()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}


// telnet localhost 8080