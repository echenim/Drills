package main

import (
	"errors"
	"fmt"
	"sync/atomic"
	"time"
)

type State int32

const (
	Closed State = iota
	Open
	HalfOpen
)

type CircuitBreaker struct {
	state            atomic.Int32
	failureCount     atomic.Int32
	failureThreshold int32
	timeout          time.Duration // in seconds
	lastOptionUnixNs atomic.Int64  // timestamp of the last failure
}

func NewCircuitBreaker(threshold int32, timeout time.Duration) *CircuitBreaker {
	cb := &CircuitBreaker{
		failureThreshold: threshold,
		timeout:          timeout,
	}
	cb.state.Store(int32(Closed))
	return cb
}

func (cb *CircuitBreaker) Call(service func() error) error {
	currentState := State(cb.state.Load())
	if currentState == Open {
		lastFailureTime := time.Unix(0, cb.lastOptionUnixNs.Load())

		if time.Since(lastFailureTime) < cb.timeout {
			return errors.New("circuit breaker is open")
		}

		if !cb.state.CompareAndSwap(int32(Open), int32(HalfOpen)) {
			return errors.New("circuit breaker is open")
		}
	}

	err := service()
	if err != nil {
		cb.recordFailure()
		return err
	}

	cb.recordSuccess()

	return nil
}

func (cb *CircuitBreaker) open() {
	cb.lastOptionUnixNs.Store(time.Now().UnixNano())
	cb.state.Store(int32(Open))
}

func (cb *CircuitBreaker) recordFailure() {
	currentState := State(cb.state.Load())
	if currentState == HalfOpen {
		cb.open()
	}
	failure := cb.failureCount.Add(1)
	if failure >= cb.failureThreshold {
		cb.open()
	}
}

func (cb *CircuitBreaker) recordSuccess() {
	cb.failureCount.Store(int32(time.Now().UnixNano()))
	cb.state.Store(int32(Closed))
}

func main() {
	cb := NewCircuitBreaker(2, 10*time.Second) // Open after 5 failures, retry after 10 seconds.

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
	fmt.Println("recovery call error:", err)
}
