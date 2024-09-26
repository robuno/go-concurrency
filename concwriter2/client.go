package main

import (
	"fmt"
	"math/rand"
	"net"
	"strings"
	"sync"
	"time"
)

var actions = []string{"get", "set", "rm"}
// var actions = []string{"get"}

// client generator with random requests
func simulateClient(clientID int) {
	// // increase val = spread load
	// time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)

	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("Client %d: Error connecting to server: %v\n", clientID, err)
		return
	}
	defer conn.Close()

	// set random action, key[0-9], and value[0-99]
	action := actions[rand.Intn(len(actions))]
	key := fmt.Sprintf("key%d", rand.Intn(10)) // key[0-9]

	var request string
	if action == "set" {
		value := fmt.Sprintf("value%d", rand.Intn(100)) // value[0-99]
		request = fmt.Sprintf("%d %s %s %s\n", clientID, action, key, value)
	} else {
		request = fmt.Sprintf("%d %s %s\n", clientID, action, key)
	}

	//fmt.Printf("Client %d: Sending request: %s", clientID, request)

	// send the request to the server
	_, err = conn.Write([]byte(request))
	if err != nil {
		fmt.Printf("Client %d: Error writing to server: %v\n", clientID, err)
		return
	}

	// read the server's response
	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("Client %d: Error reading from server: %v\n", clientID, err)
		return
	}

	// print response
	response := strings.TrimSpace(string(buffer[:n]))
	fmt.Printf("Client %d: Received response: %s\n", clientID, response)
	// if strings.HasPrefix(response, "[ERROR]") {
	// 	fmt.Printf("## Error Client %d: Received response: %s\n", clientID, response)
	// } else {
	// 	fmt.Printf("Client %d: Received response: %s\n", clientID, response)
	// }
}

func main() {
	var wg sync.WaitGroup
	rand.Seed(time.Now().UnixNano())

	clientCount := 50

	for i := 1; i <= clientCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			simulateClient(i) // each client = a new goroutine
		}()
	}

	wg.Wait()
}
