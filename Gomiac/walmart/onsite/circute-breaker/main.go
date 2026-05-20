// Build a circuit breaker that wraps a function call.
// It should have 3 states:

// Closed: allow requests normally
// Open: block requests immediately
// HalfOpen: after cooldown, allow one trial request

// If failures exceed a threshold, open the circuit. After a timeout, move to half-open.
// If the trial succeeds, close it. If it fails, open again.

package main

import (
	"errors"
	"fmt"
	"time"

	nonthread "gomaic/walmart/onsite/circute-breaker/non-thread"
	threadsafe "gomaic/walmart/onsite/circute-breaker/thread-safe"
)

func noneThreadBreaker() {
	cb := nonthread.NewCircuitBreaker(10, 10*time.Second) // Open after 5 failures, retry after 10 seconds.

	failingService := func() error {
		return errors.New("service failed")
	}

	for i := 0; i < 100; i++ {
		err := cb.Call(failingService)
		if err != nil {
			println("Request failed:", err.Error())
		} else {
			println("Request succeeded")
		}
	}

	time.Sleep(11 * time.Second) // Wait for the timeout to expire

	healthCheck := func() error {
		return nil // Simulate a successful health check
	}

	err := cb.Call(healthCheck)
	if err != nil {
		println("Health check failed:", err.Error())
	} else {
		println("Health check succeeded")
	}
}

func threadBreaker() {
	cb := threadsafe.NewCircuitBreaker(10, 10*time.Second) // Open after 5 failures, retry after 10 seconds.
	failingService := func() error {
		return errors.New("service failed")
	}

	for i := 0; i < 100; i++ {
		err := cb.Call(failingService)
		if err != nil {
			println("Request failed:", err.Error())
		} else {
			println("Request succeeded")
		}
	}

	time.Sleep(11 * time.Second) // Wait for the timeout to expire

	healthCheck := func() error {
		return nil // Simulate a successful health check
	}

	err := cb.Call(healthCheck)
	if err != nil {
		println("Health check failed:", err.Error())
	} else {
		println("Health check succeeded")
	}
}

func main() {
	noneThreadBreaker()
	fmt.Println("---------------------------")
	threadBreaker()
}
