package main

import (
	"fmt"
	"sync"
	"time"
)

type Logger struct {
	// jobs is the bounded queue.
	// If the channel is full, sending to it will block.
	jobs chan string

	// wg tracks all worker goroutines.
	// shutdown waits for all workers to finish.
	wg sync.WaitGroup

	// mu protects the shutdown flag.
	mu sync.Mutex

	// closed tells us whether shutdown has started.
	closed bool
}

// NewLogger creates a logger.
//
// maxConcurrent controls how many messages can be processed at the same time.
// queueSize controls how many messages can wait in the queue.
func NewLogger(maxConcurrent int, queueSize int) *Logger {
	logger := &Logger{
		jobs: make(chan string, queueSize),
	}

	// Start maxConcurrent workers.
	// This guarantees that at most maxConcurrent log messages
	// are processed at the same time.
	for i := 0; i < maxConcurrent; i++ {
		logger.wg.Add(1)

		go func(workerID int) {
			defer logger.wg.Done()

			// Each worker keeps reading from the jobs channel.
			// When the jobs channel is closed and empty,
			// this loop exits automatically.
			for message := range logger.jobs {
				logger.process(workerID, message)
			}
		}(i + 1)
	}

	return logger
}

// Log submits a message to the logger.
//
// This method is thread-safe.
// If the queue is full, this method blocks until space is available.
// If shutdown has already started, this method returns immediately.
func (l *Logger) Log(message string) {
	// First check whether shutdown has started.
	l.mu.Lock()

	if l.closed {
		// After shutdown, new log calls should return immediately.
		l.mu.Unlock()
		return
	}

	// Keep the lock while sending to the channel.
	// This prevents shutdown from closing the channel while
	// this goroutine is trying to send.
	l.jobs <- message

	l.mu.Unlock()
}

// Shutdown stops the logger.
//
// It prevents future log messages from being accepted.
// It waits until all queued messages are processed.
func (l *Logger) Shutdown() {
	l.mu.Lock()

	if l.closed {
		// If shutdown was already called, do nothing.
		l.mu.Unlock()
		return
	}

	// Mark logger as closed.
	// Future calls to Log will return immediately.
	l.closed = true

	// Closing the jobs channel tells workers:
	// "finish whatever is already queued, then exit."
	close(l.jobs)

	l.mu.Unlock()

	// Wait for all workers to finish processing queued messages.
	l.wg.Wait()
}

// process simulates actual log processing.
// In a real system, this might write to stdout, a file, Kafka, etc.
func (l *Logger) process(workerID int, message string) {
	fmt.Printf("worker %d processing: %s\n", workerID, message)

	// Simulate slow processing.
	time.Sleep(500 * time.Millisecond)
}

func main() {
	// At most 3 messages can be processed concurrently.
	// At most 5 messages can wait in the queue.
	logger := NewLogger(10, 12)

	var wg sync.WaitGroup

	// Simulate many goroutines calling Log at the same time.
	for i := 1; i <= 120; i++ {
		wg.Add(1)

		go func(i int) {
			defer wg.Done()

			logger.Log(fmt.Sprintf("message-%d", i))
		}(i)
	}

	// Wait until all callers have submitted their messages.
	wg.Wait()

	// Shutdown waits for queued messages to finish.
	logger.Shutdown()

	// This message will not be accepted.
	// It returns immediately.
	logger.Log("message-after-shutdown")

	fmt.Println("logger shutdown complete")
}
