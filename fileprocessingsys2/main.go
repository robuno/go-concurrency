// shared one source

package main

import (
	"fmt"
	"sync"

	// "path/filepath"
	"os"
	"io/ioutil"
	"time"
)

func readFile(filePath string) (string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	openedTime := time.Now().Format("15:04:05.000")
	newData := fmt.Sprintf("%s[File Opened At: %s]\n", string(data), openedTime)

	return newData, nil
}

func appendToFile(filePath string, content string, mu *sync.Mutex) error {
	mu.Lock()
	defer mu.Unlock()

	// open in append mode
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// write the content 
	if _, err := file.WriteString(content); err != nil {
		return err
	}

	return nil
}

// func modifyFile(filePath string, content string, goroutineID int) error {
// 	updatedTime := time.Now().Format("15:04:05.000")
// 	modifiedData := fmt.Sprintf("%sGoroutine [%d] accessed and wrote to this file!\n[File Updated At: %s]\n",
// 		content,
// 		goroutineID,
// 		updatedTime)

// 	// Write the modified content
// 	err := ioutil.WriteFile(filePath, []byte(modifiedData), 0644)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func processFile(filePath string, ch chan<- string, wg *sync.WaitGroup, mu *sync.Mutex, goroutineID int) {
	defer wg.Done()

	// log the start
	startTime := time.Now().Format("15:04:05.000")
	mu.Lock()
	fmt.Printf("Goroutine [%d] Processing file: %s at %s\n",
		goroutineID,
		filePath,
		startTime)
	mu.Unlock()

	// read file and add current time to the file
	_, err := readFile(filePath)
	if err != nil {
		ch <- fmt.Sprintf("Goroutine [%d] Error reading file %s: %v",
			goroutineID,
			filePath,
			err)
		return
	}

	// modify file and add current time to the file
    modifyTime := time.Now().Format("15:04:05.000")
	modifiedContent := fmt.Sprintf("Goroutine [%d] accessed the file!  [File Updated At: %s]\n", goroutineID, modifyTime)
	err = appendToFile(filePath, modifiedContent, mu)
	if err != nil {
		ch <- fmt.Sprintf("Goroutine [%d] Error writing to file %s: %v", 
                            goroutineID, 
                            filePath, 
                            err)
		return
	}

	// log the end
	endTime := time.Now().Format("15:04:05.000")
	mu.Lock()
	fmt.Printf("Goroutine [%d] Completed processing file: %s at %s\n",
		goroutineID,
		filePath,
		endTime)
	mu.Unlock()

	ch <- fmt.Sprintf("Goroutine [%d] Successfully processed %s",
		goroutineID,
		filePath)
}

func main() {
	const filePath = "files/shared_file.txt"
	const numGoroutines = 10

	statusChannel := make(chan string)
	// statusChannel := make(chan string, len(filePaths))

	var mu sync.Mutex
	var wg sync.WaitGroup
	// a goroutine for the same file
	for i := 1; i <= numGoroutines; i++ {
		wg.Add(1)
		go processFile(filePath, statusChannel, &wg, &mu, i)
	}

	go func() {
		wg.Wait()
		close(statusChannel)
	}()

	for status := range statusChannel {
		fmt.Println(status)
	}
}
