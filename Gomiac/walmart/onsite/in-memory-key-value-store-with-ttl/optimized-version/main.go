package main

import (
	"container/heap"
	"fmt"
	"sync"
	"time"
)

type entry struct {
	value    string
	expireAt time.Time
}

type expirationItem struct {
	key      string
	expireAt time.Time
}

type expirationHeap []expirationItem

func (h expirationHeap) Len() int {
	return len(h)
}

func (h expirationHeap) Less(i, j int) bool {
	return h[i].expireAt.Before(h[j].expireAt)
}

func (h expirationHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *expirationHeap) Push(x any) {
	*h = append(*h, x.(expirationItem))
}

func (h *expirationHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

func (h expirationHeap) Peek() expirationItem {
	return h[0]
}

type TTLStore struct {
	mu       sync.RWMutex
	items    map[string]entry
	expHeap  expirationHeap
	stopChan chan struct{}
}

func NewTTLStore(cleanupInterval time.Duration) *TTLStore {
	store := &TTLStore{
		items:    make(map[string]entry),
		expHeap:  make(expirationHeap, 0),
		stopChan: make(chan struct{}),
	}
	heap.Init(&store.expHeap)
	go store.cleanupExpiredEntries(cleanupInterval)
	return store
}

func (s *TTLStore) Set(key, value string, ttl time.Duration) {
	expireAt := time.Now().Add(ttl)
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[key] = entry{
		value:    value,
		expireAt: expireAt,
	}
	heap.Push(&s.expHeap, expirationItem{
		key:      key,
		expireAt: expireAt,
	})
}

func (s *TTLStore) Get(key string) (string, bool) {
	s.mu.RLock()
	item, exists := s.items[key]
	s.mu.RUnlock()
	if !exists {
		return "", false
	}

	if time.Now().After(item.expireAt) {
		s.Remove(key)
		return "", false
	}
	return item.value, true
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

	for s.expHeap.Len() > 0 {
		next := s.expHeap.Peek()

		// Since this is a min-heap, if the earliest item has not expired,
		// nothing else after it has expired either.
		if next.expireAt.After(now) {
			break
		}

		expired := heap.Pop(&s.expHeap).(expirationItem)

		current, exist := s.items[expired.key]
		if !exist {
			continue
		}

		// Lazy deletion:
		// If the map has a newer expiration time, this heap item is stale.
		if exist && current.expireAt.Equal(next.expireAt) {
			continue
		}

		delete(s.items, expired.key)
	}
}

func (s *TTLStore) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return len(s.items)
}

func main() {
	store := NewTTLStore(500 * time.Millisecond)
	defer store.Stop()

	fmt.Println("starting TTL store")

	store.Set("a", "first", 2*time.Second)
	store.Set("b", "second", 5*time.Second)

	fmt.Println("added: a TTL=2s")
	fmt.Println("added: b TTL=5s")
	fmt.Println("store size:", store.Len())

	value, ok := store.Get("a")
	fmt.Println("get a:", value, ok)

	time.Sleep(3 * time.Second)

	// At this point, key "a" should have expired automatically.
	// We are checking store size to show that cleanup removed it.
	fmt.Println("\nafter 3 seconds")
	fmt.Println("store size:", store.Len())

	value, ok = store.Get("a")
	fmt.Println("get a:", value, ok) // "" false

	value, ok = store.Get("b")
	fmt.Println("get b:", value, ok) // second true

	store.Set("c", "third", 2*time.Second)
	fmt.Println("\nadded: c TTL=2s")
	fmt.Println("store size:", store.Len())

	time.Sleep(3 * time.Second)

	// Now c should be gone, but b may also be gone depending on timing.
	// b had a 5s TTL from the beginning.
	fmt.Println("\nafter another 3 seconds")
	fmt.Println("store size:", store.Len())

	value, ok = store.Get("b")
	fmt.Println("get b:", value, ok)

	value, ok = store.Get("c")
	fmt.Println("get c:", value, ok)

	store.Set("d", "fourth", 4*time.Second)
	fmt.Println("\nadded: d TTL=4s")
	fmt.Println("store size:", store.Len())

	for i := 1; i <= 5; i++ {
		time.Sleep(1 * time.Second)

		value, ok := store.Get("d")
		fmt.Printf("after %ds, get d: value=%q ok=%v storeSize=%d\n",
			i,
			value,
			ok,
			store.Len(),
		)
	}

	fmt.Println("\ndone")
}
