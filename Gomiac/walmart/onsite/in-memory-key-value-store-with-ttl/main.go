package main

import (
	"fmt"
	"sync"
	"time"
)

type Entry struct {
	value    string
	expireAt time.Time
}

type TTLStore struct {
	items    map[string]Entry
	mu       sync.RWMutex
	stopChan chan struct{}
}

func NewTTLStore(cleanupInterval time.Duration) *TTLStore {
	store := &TTLStore{
		items:    make(map[string]Entry),
		stopChan: make(chan struct{}),
	}
	go store.cleanupExpiredEntries(cleanupInterval)
	return store
}

func (s *TTLStore) Set(key, value string, ttl time.Duration) {
	s.mu.Lock()
	s.items[key] = Entry{
		value:    value,
		expireAt: time.Now().Add(time.Duration(ttl) * time.Second),
	}
	s.mu.Unlock()
}

func (s *TTLStore) Get(key string) (string, bool) {
	s.mu.RLock()
	entry, exists := s.items[key]
	s.mu.RUnlock()
	if !exists {
		return "", false
	}

	if time.Now().After(entry.expireAt) {
		s.Remove(key)
		return "", false
	}
	return entry.value, true
}

func (s *TTLStore) Remove(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.items, key)
}

func (s *TTLStore) cleanupExpiredEntries(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.deleteExpiredEntries()
		case <-s.stopChan:
			return
		}
	}
}

func (s *TTLStore) Stop() {
	close(s.stopChan)
}

func (s *TTLStore) deleteExpiredEntries() {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	for key, entry := range s.items {
		if now.After(entry.expireAt) {
			delete(s.items, key)
		}
	}
}

func main() {
	store := NewTTLStore(1 * time.Second)
	defer store.Stop()

	store.Set("user:1", "William", 2*time.Second)
	store.Set("user:2", "Alice", 5*time.Second)
	store.Set("user:3", "Bob", 10*time.Second)
	store.Set("user:4", "Charlie", 1*time.Second)

	value, ok := store.Get("user:1")
	fmt.Println(value, ok) // William true

	time.Sleep(3 * time.Second)

	value, ok = store.Get("user:1")
	fmt.Println(value, ok) // "" false
}
