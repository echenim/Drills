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
	"time"
)

type State int // Define the states of the circuit breaker

const (
	Closed State = iota
	Open
	HalfOpen
)

type CircuitBreaker struct {
	state            State
	failureCount     int
	failureThreshold int
	timeout          time.Duration // in seconds
	lastFailureTime  time.Time     // timestamp of the last failure
}

func NewCircuitBreaker(threshold int, timeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		failureThreshold: threshold,
		timeout:          timeout,
		state:            Closed,
	}
}

func (cb *CircuitBreaker) Call(service func() error) error {
	if cb.state == Open {
		if time.Since(cb.lastFailureTime) > cb.timeout {
			cb.state = HalfOpen // Move to half-open after timeout
		} else {
			return errors.New("circuit is open")
		}
	}

	err := service()
	if err != nil {
		cb.failureCount++

		if cb.state == HalfOpen || cb.failureCount >= cb.failureThreshold {
			cb.state = Open                 // Open the circuit if in half-open or threshold exceeded
			cb.lastFailureTime = time.Now() // Update the last failure time
		}
		return err
	}

	cb.failureCount = 0 // Reset failure count on success
	cb.state = Closed   // Ensure the circuit is closed on success

	return nil
}

func main() {
	cb := NewCircuitBreaker(5, 10*time.Second) // Open after 5 failures, retry after 10 seconds.

	failingService := func() error {
		return errors.New("service failed")
	}

	for i := 0; i < 10; i++ {
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
