package main

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	id int
}

// run as a goroutine, processing tasks from the taskChannel
func worker(id int, wg *sync.WaitGroup, taskChannel <-chan Task) {
	defer wg.Done()
	for task := range taskChannel {
		fmt.Printf("Worker %d started task %d\n", id, task.id)
		time.Sleep(1 * time.Second) 
		fmt.Printf("Worker %d completed task %d\n", id, task.id)
	}
}

func main() {
	numWorkers := 3
	numTasks := 10

	taskChannel := make(chan Task, 2)

	var wg sync.WaitGroup
	for i := 1; i <= numWorkers; i++ {
		wg.Add(1)
		go worker(i, &wg, taskChannel)
	}

	for i := 1; i <= numTasks; i++ {
		taskChannel <- Task{id: i}
	}

	close(taskChannel)
	wg.Wait()

	fmt.Println("All tasks completed.")
}
