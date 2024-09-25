package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func HandleConnection(conn net.Conn, clientID string) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("%s disconnected", clientID)
			return
		}

		// get client message
		message = strings.TrimSpace(message)
		parts := strings.Split(message, " ")

		// check whether request format okay
		if len(parts) < 2 {
			conn.Write([]byte("Invalid command\n"))
			continue
		}

		// get msg parts
		command := parts[0]
		key := parts[1]
		value := ""

		if command == "set" && len(parts) >= 3 {
			value = strings.Join(parts[2:], " ")
		}

		// log file
		switch command {
		case "get":
			err := logAndRetry(clientID, "get", key, "")
			if err != nil {
				log.Printf("Error: %v", err)
				conn.Write([]byte("Error in logging [GET]\n"))
			} else {
				conn.Write([]byte(fmt.Sprintf("GET request for key: %s\n", key)))
			}

		case "set":
			err := logAndRetry(clientID, "set", key, value)
			if err != nil {
				log.Printf("Error: %v", err)
				conn.Write([]byte("Error in logging [SET]\n"))
			} else {
				conn.Write([]byte(fmt.Sprintf("SET request for key: %s, value: %s\n", key, value)))
			}

		case "rm":
			err := logAndRetry(clientID, "rm", key, "")
			if err != nil {
				log.Printf("Error: %v", err)
				conn.Write([]byte("Error in logging [RM]\n"))
			} else {
				conn.Write([]byte(fmt.Sprintf("RM request for key: %s\n", key)))
			}

		default:
			conn.Write([]byte("Unknown command\n"))
		}
	}
}
