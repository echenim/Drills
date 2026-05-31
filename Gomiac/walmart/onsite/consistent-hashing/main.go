package main

import (
	"fmt"
	"hash/crc32"
	"sort"
	"sync"
)

// ConsistentHash implements a thread-safe consistent hashing ring
type ConsistentHash struct {
	mu       sync.RWMutex
	replicas int
	hashFunc func([]byte) uint32
	ring     []uint32          // Sorted positions of virtual nodes
	ringMap  map[uint32]string // Maps virtual node position -> physical node name
}

// New creates a new ConsistentHash ring with given replica count and hash function.
// If hashFunc is nil, crc32.ChecksumIEEE is used (fast & standard for load balancing).
func NewConsistentHash(replicas int, hashFunc func([]byte) uint32) *ConsistentHash {
	if hashFunc == nil {
		hashFunc = crc32.ChecksumIEEE
	}
	return &ConsistentHash{
		replicas: replicas,
		hashFunc: hashFunc,
		ring:     make([]uint32, 0),
		ringMap:  make(map[uint32]string),
	}
}

// Add registers one or more physical nodes to the ring
func (c *ConsistentHash) Add(nodes ...string) {
	for _, node := range nodes {
		for i := 0; i < c.replicas; i++ {
			// Virtual node naming: "node#replicaIndex"
			hash := c.hashFunc([]byte(fmt.Sprintf("%s#%d", node, i)))
			c.ring = append(c.ring, hash)
			c.ringMap[hash] = node
		}
	}

	// Keep ring sorted for binary search
	sort.Slice(c.ring, func(i, j int) bool { return c.ring[i] < c.ring[j] })
}

// Get returns the physical node responsible for the given key
func (c *ConsistentHash) Get(key string) string {
	if len(c.ring) == 0 {
		return ""
	}
	hash := c.hashFunc([]byte(key))

	// Find first virtual node >= key hash
	dx := sort.Search(len(c.ring), func(i int) bool { return c.ring[i] >= hash })
	if dx == len(c.ring) {
		dx = 0 // Wrap around to the first node
	}
	return c.ringMap[c.ring[dx]]
}

// Remove unregisters one or more physical nodes from the ring
func (c *ConsistentHash) Remove(node ...string) {
	for _, n := range node {
		for i := 0; i < c.replicas; i++ {
			hash := c.hashFunc([]byte(fmt.Sprintf("%s#%d", n, i)))
			delete(c.ringMap, hash)

			// Binary search to find position, then remove from slice
			dx := sort.Search(len(c.ring), func(i int) bool { return c.ring[i] >= hash })
			if dx < len(c.ring) && c.ring[dx] == hash {
				c.ring = append(c.ring[:dx], c.ring[dx+1:]...)
			}
		}
	}
}

func main() {
	ch := NewConsistentHash(300, nil) // 3 replicas, default hash function

	ch.Add("NodeA", "NodeB", "NodeC")
	fmt.Println(ch.Get("myKey1")) // Should consistently return the same node for "myKey1"
	fmt.Println(ch.Get("myKey2")) // Should consistently return the same node for "myKey2"

	keys := []string{"myKey1", "myKey2", "myKey3", "myKey4", "myKey5"}
	for _, key := range keys {
		fmt.Printf("Key: %s -> Node: %s\n", key, ch.Get(key))
	}

	ch.Remove("NodeB")
	fmt.Println(ch.Get("myKey1")) // May return a different node after removal
	fmt.Println(ch.Get("myKey2")) // May return a different node after removal

	for _, key := range keys {
		fmt.Printf("Key: %s -> Node: %s\n", key, ch.Get(key))
	}
}
