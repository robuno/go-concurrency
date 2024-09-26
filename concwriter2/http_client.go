package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	// parse the cli arguments
	if len(os.Args) < 3 {
		fmt.Println("Required command usage: go run http_client.go get <key>")
		return
	}

	// parse format: "get key5"
	command := os.Args[1]
	key := os.Args[2]

	// ensure the command is valid
	supportedFormats := "[get]"
	if strings.ToLower(command) != "get" {
		fmt.Printf("Invalid command. Use one of the supported formats: %s\n", supportedFormats)
		return
	}

	// construct the url
	url := fmt.Sprintf("http://localhost:8085/get?key=%s", key)

	// send HTTP request
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error sending request: %v...\n", err)
		return
	}
	defer resp.Body.Close()

	// check the response
	if resp.StatusCode == http.StatusOK {
		// read the response body
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("Error reading response: %v...\n", err)
			return
		}

		// print the content of the response
		fmt.Printf("Content:\n%s\n", string(body))
	} else {
		// handle non-200 status
		fmt.Printf("Error: StatusCode: %d, StatusDesc: %s\n", resp.StatusCode, resp.Status)
	}
}
