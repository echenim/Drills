package threadsafe

import (
	"errors"
	"sync"
	"time"
)

type State int // Define the states of the circuit breaker

const (
	Closed State = iota
	Open
	HalfOpen
)

type CircuitBreaker struct {
	mu               sync.Mutex
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
	cb.mu.Lock()
	if cb.state == Open {
		if time.Since(cb.lastFailureTime) > cb.timeout {
			cb.state = HalfOpen // Move to half-open after timeout
		} else {
			cb.mu.Unlock() // Unlock before returning
			return errors.New("circuit is open")
		}
	}
	cb.mu.Unlock() // Unlock before calling the service

	err := service()

	cb.mu.Lock()
	defer cb.mu.Unlock()

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
