package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

func searchLogFile(key string) ([]string, error) {
	var results []string

	// ppen the logs
	file, err := os.Open("log.txt")
	if err != nil {
		return nil, fmt.Errorf("Error opening log file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, key) {		// check if the line contains the key
			results = append(results, line)
		}
	}

	// return scanning errors
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error reading log file: %v", err)
	}

	// if key is not available
	if len(results) == 0 {
		return nil, fmt.Errorf("No logs found for key: %s", key)
	}

	return results, nil
}

func getLogHandler(w http.ResponseWriter, r *http.Request) {
	
	// parse the key from url
	keys, ok := r.URL.Query()["key"]
	if !ok || len(keys[0]) < 1 {
		http.Error(w, "URL parameter 'key' is missing...", http.StatusBadRequest)
		return
	}

	key := keys[0]
	fmt.Printf("Searching logs for key: %s...\n", key)

	// Search the log file for the specified key
	results, err := searchLogFile(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Print the results in the HTTP response
	for _, result := range results {
		fmt.Fprintf(w, "%s\n", result)
	}
}

func main() {
	// HTTP Routers
	http.HandleFunc("/get", getLogHandler)

	// Start the HTTP server
	selectedPort := ":8085"
	fmt.Printf("Starting HTTP server on port %s...\n", selectedPort)
	err := http.ListenAndServe(selectedPort, nil)
	if err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}


// curl "http://localhost:8085/get?key=key5"