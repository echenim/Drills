// Write a component that accepts an array of tasks (functions that return an error).
// It must execute these tasks concurrently using a fixed maximum number of worker goroutines ($N$).
// If any task returns a fatal error, the remaining unstarted tasks should be canceled immediately,
// and the component should return the first error encountered.

package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

func executeTasks(tasks []func() error, maxWorkers int) error {
	if maxWorkers <= 0 {
		maxWorkers = 1
	}

	var (
		wg       sync.WaitGroup
		firstErr error
		errOnce  sync.Once
		done     = make(chan struct{})
		sem      = make(chan struct{}, maxWorkers)
	)
DispatchLoop:
	for _, task := range tasks {
		select {
		case <-done:
			break DispatchLoop
		default:
		}

		select {
		case sem <- struct{}{}:
		case <-done:
			break DispatchLoop
		}

		wg.Add(1)
		go func(t func() error) {
			defer wg.Done()
			defer func() { <-sem }()
			if err := t(); err != nil {
				errOnce.Do(func() {
					firstErr = err
					close(done)
				})
			}
		}(task)
	}
	wg.Wait()
	return firstErr
}

func main() {
	tasks := []func() error{
		func() error { time.Sleep(100 * time.Millisecond); return nil },
		func() error { time.Sleep(50 * time.Millisecond); return errors.New("fatal failure in task 2") },
		func() error { time.Sleep(200 * time.Millisecond); return nil }, // Will run concurrently
		func() error { time.Sleep(10 * time.Millisecond); return nil },  // Will likely be skipped
	}

	err := executeTasks(tasks, 2)
	if err != nil {
		fmt.Printf("Aborted with error: %v\n", err)
	} else {
		fmt.Println("All tasks completed successfully")
	}
}
