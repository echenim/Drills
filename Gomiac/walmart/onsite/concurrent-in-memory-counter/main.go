package main

import (
	"fmt"
	"sync"
)

type RequestCounter struct {
	counts map[string]int64
	mu     sync.RWMutex
}

func NewRequestCounter() *RequestCounter {
	return &RequestCounter{
		counts: make(map[string]int64),
	}
}

func (rc *RequestCounter) Increment(endpoint string) {
	rc.mu.Lock()
	defer rc.mu.Unlock()
	rc.counts[endpoint]++
}

func (rc *RequestCounter) Get(endpoint string) int64 {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	return rc.counts[endpoint]
}

func (rc *RequestCounter) Snapshot() map[string]int64 {
	rc.mu.RLock()
	defer rc.mu.RUnlock()
	copy := make(map[string]int64, len(rc.counts))
	for endpoint, count := range rc.counts {
		copy[endpoint] = count
	}
	return copy
}

func main() {
	counter := NewRequestCounter()

	var wg sync.WaitGroup

	endpoints := []string{
		"/api/users",
		"/api/orders",
		"/api/users",
		"/api/payments",
		"/api/users",
	}

	for _, endpoint := range endpoints {
		wg.Add(1)

		go func(ep string) {
			defer wg.Done()
			counter.Increment(ep)
		}(endpoint)
	}

	wg.Wait()

	fmt.Println(counter.Get("/api/users"))
	fmt.Println(counter.Snapshot())
}
