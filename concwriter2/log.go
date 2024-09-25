package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var mu sync.Mutex

const maxRetries = 2
const retryDelay = 1 * time.Millisecond

func logAndRetry(clientID, action, key, value string) error {
	retryCount := 0
	var err error

	for retryCount < maxRetries {
		err = logWriter(clientID, action, key, value)
		if err == nil {
			return nil // successful writing
		}

		log.Printf("Error writing to file: %v. Retrying %d/%d...", err, retryCount+1, maxRetries)
		retryCount++

		// wait before retry
		time.Sleep(retryDelay)
	}

	return errors.New("Max retries exceeded! Failed to write to logs!")
}

func logWriter(clientID, action, key, value string) error {
	mu.Lock() // only one goroutine can accesss file
	defer mu.Unlock()

	// open txt
	f, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close() // ensure file is closed

	timestamp := time.Now().Format("02-01-2006 15:04:05")
	logLine := fmt.Sprintf("%s, %s, %s, %s, time: %s\n", clientID, action, key, value, timestamp)
	if _, err := f.WriteString(logLine); err != nil {
		return err
	}

	return nil // successful writing
}
