package main

type CacheItem struct {
	value string
	key   int
	prev  *CacheItem
	next  *CacheItem
}

type LRUCache struct {
	capacity int
	items    map[int]*CacheItem
	head     *CacheItem
	tail     *CacheItem
}

func NewLRUCache(capacity int) *LRUCache {
	head := &CacheItem{}
	tail := &CacheItem{}

	head.next = tail
	tail.prev = head

	return &LRUCache{
		capacity: capacity,
		items:    make(map[int]*CacheItem),
		head:     head,
		tail:     tail,
	}
}

func (c *LRUCache) Get(key int) string {
	return ""
}

func (c *LRUCache) Put(key int, value string) {
	if node, exists := c.items[key]; exists {
		node.value = value
		c.remove(node)
		c.addToFront(node)
		return
	}

	node := &CacheItem{
		value: value,
		key:   key,
	}

	c.items[key] = node
	c.addToFront(node)

	if len(c.items) > c.capacity {
		lru := c.tail.prev
		c.remove(lru)
		delete(c.items, lru.key)

	}
}

func (c *LRUCache) remove(node *CacheItem) {
	prev := node.prev
	next := node.next

	prev.next = next
	next.prev = prev
}

func (c *LRUCache) addToFront(node *CacheItem) {
	first := c.head.next

	node.prev = c.head
	node.next = first

	c.head.next = node
	first.prev = node
}

func main() {
	cache := NewLRUCache(2)

	cache.Put(1, "one")
	cache.Put(2, "two")

	println(cache.Get(1)) // Output: "one"

	cache.Put(3, "three") // Evicts key 2

	println(cache.Get(2)) // Output: "Data does not exist"

	cache.Put(4, "four") // Evicts key 3

	println(cache.Get(3)) // Output: "Data does not exist"
	println(cache.Get(4)) // Output: "four"
}
