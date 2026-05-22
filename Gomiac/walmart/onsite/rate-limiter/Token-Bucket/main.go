package main

import (
	"sync"
	"time"
)

type TokenBucket struct {
	rate           int64
	capacity       int64
	tokens         int64
	lastRefillTime int64
	mu             sync.Mutex
}

func NewTokenBucket(rate, capacity int64) *TokenBucket {
	return &TokenBucket{
		rate:           rate,
		capacity:       capacity,
		tokens:         capacity,
		lastRefillTime: time.Now().UnixNano(),
	}
}

func (tb *TokenBucket) refill() {
	now := time.Now().UnixNano()
	elapsed := now - tb.lastRefillTime
	newToken := (elapsed * tb.rate) / int64(time.Second)
	if newToken > 0 {
		tb.tokens = min(tb.capacity, tb.tokens+newToken)
		tb.lastRefillTime = now
	}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	tb.refill()
	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	tb := NewTokenBucket(5, 10) // 5 tokens per second, capacity of 10 tokens
	for i := 0; i < 105; i++ {
		if tb.Allow() {
			println("Request allowed")
		} else {
			println("Request denied")
		}
		time.Sleep(10 * time.Millisecond) // Simulate requests every 200ms
	}
}
