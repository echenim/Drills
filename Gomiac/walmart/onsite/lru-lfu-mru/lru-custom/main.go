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

		return
	}

	delete(c.keys, c.Data[c.Capacity-1])
	c.Data[c.Capacity-1] = key
	c.keys[key] = c.Capacity - 1
	if c.Capacity != 0 {
		tmp := c.Data[0]
		c.Data[0] = key
		c.Data[c.keys[key]] = tmp
	}
}

func main() {
	cache := NewLRUCache(3)
	cache.Put("apple")
	cache.Put("banana")
	cache.Put("cherry")

	fmt.Printf("\n Initial state: %v \n", cache.Data)

	println(cache.Get("banana")) // Output: banana
	fmt.Printf("\n Initial state: %v \n", cache.Data)
	println(cache.Get("apple")) // Output: apple
	fmt.Printf("\n Initial state: %v \n", cache.Data)

	cache.Put("date") // Evicts "cherry"
	fmt.Printf("\n Initial state: %v \n", cache.Data)

	println(cache.Get("cherry")) // Output: Data does not exist
	println(cache.Get("date"))   // Output: date
}
