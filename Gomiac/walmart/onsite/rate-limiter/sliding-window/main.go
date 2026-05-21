package main

import (
	"fmt"
	"time"
)

type RateLimiter struct {
	limit    int
	window   time.Duration
	requests map[string][]time.Time
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		limit:    limit,
		window:   window,
		requests: make(map[string][]time.Time),
	}
}

func (r *RateLimiter) Allow(clientID string, now time.Time) bool {
	// Edge case: if limit is 0 or negative, never allow requests
	if r.limit <= 0 {
		return false
	}

	// Retrieve existing timestamps for this client
	timestamp := r.requests[clientID]

	// Calculate the cutoff time (exclusive lower bound of the sliding window)
	cutOffWindow := now.Add(-r.window)

	// Filter out expired timestamps by reusing the underlying array for efficiency
	// This avoids unnecessary allocations while maintaining O(n) cleanup
	valid := timestamp[:0]
	for _, t := range timestamp {
		if t.After(cutOffWindow) {
			valid = append(valid, t)
		}
	}

	// Check if the client is under the rate limit
	if len(valid) < r.limit {
		valid = append(valid, now)
		r.requests[clientID] = valid
		return true
	}

	// Rate limit exceeded: update map with cleaned timestamps for memory management
	// Remove the client entry if no valid timestamps remain
	if len(valid) == 0 {
		delete(r.requests, clientID)
	} else {
		r.requests[clientID] = valid
	}
	return false
}

func main() {
	rl := NewRateLimiter(3, time.Minute)

	now := time.Now()
	fmt.Println(rl.Allow("user-1", now)) // true (1/3)
	fmt.Println(rl.Allow("user-1", now)) // true (2/3)
	fmt.Println(rl.Allow("user-1", now)) // true (3/3)
	fmt.Println(rl.Allow("user-1", now)) // false (limit exceeded)
	fmt.Println(rl.Allow("user-2", now)) // true (different client)

	// After window expires, user-1 can make requests again
	fmt.Println(rl.Allow("user-1", now.Add(61*time.Second))) // true
}
