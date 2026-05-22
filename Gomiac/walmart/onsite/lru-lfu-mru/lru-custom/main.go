package main

import "fmt"

type LRUCache struct {
	Data            []string
	keys            map[string]int
	Capacity        int
	CurrentCapacity int
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		Data:            make([]string, capacity),
		keys:            make(map[string]int),
		Capacity:        capacity,
		CurrentCapacity: 0,
	}
}

func (c *LRUCache) Get(key string) string {
	_, exists := c.keys[key]
	if exists {
		if c.keys[key] != 0 {
			tmp := c.Data[0]
			c.Data[0] = key
			c.Data[c.keys[key]] = tmp

			c.keys[tmp] = c.keys[key]
			c.keys[key] = 0

		}
		return key
	}
	return "Data does not exist"
}

func (c *LRUCache) Put(key string) {
	if c.CurrentCapacity < c.Capacity {
		c.Data[c.CurrentCapacity] = key
		c.keys[key] = c.CurrentCapacity
		c.CurrentCapacity++

		if c.keys[key] != 0 {
			c.moveToFront(key)
		}

		return
	}

	// Evict LRU (tail)
	lruKey := c.Data[c.Capacity-1]
	delete(c.keys, lruKey)

	// Insert new key at tail position first
	c.Data[c.Capacity-1] = key
	c.keys[key] = c.Capacity - 1

	// Then swap to front (and update BOTH map entries)
	if c.Capacity > 0 {
		tmp := c.Data[0]
		c.Data[0] = key
		c.Data[c.Capacity-1] = tmp // Use Capacity-1 directly, not c.keys[key]

		c.keys[key] = 0
		c.keys[tmp] = c.Capacity - 1 // ✅ Critical: update the swapped element
	}
}

func (c *LRUCache) moveToFront(key string) {
	if c.keys[key] != 0 {
		tmp := c.Data[0]
		idx := c.keys[key]
		c.Data[0] = key
		c.Data[idx] = tmp
		c.keys[key] = 0
		c.keys[tmp] = idx
	}
}

func main() {
	cache := NewLRUCache(3)
	cache.Put("apple")
	cache.Put("banana")
	cache.Put("cherry")

	fmt.Printf("\n Initial state: %v \n %v \n", cache.Data, cache.keys)

	println(cache.Get("banana")) // Output: banana
	fmt.Printf("\n Initial state: %v \n %v \n", cache.Data, cache.keys)
	println(cache.Get("apple")) // Output: apple
	fmt.Printf("\n Initial state: %v \n %v \n", cache.Data, cache.keys)

	cache.Put("date") // Evicts "cherry"
	fmt.Printf("\n Initial state: %v \n %v \n", cache.Data, cache.keys)
	println(cache.Get("cherry")) // Output: Data does not exist
	println(cache.Get("date"))   // Output: date
}
