package main

type CacheItem struct {
	value string
	key   int
	prev  *CacheItem
	next  *CacheItem
}

type LRUCache struct {
	Data            CacheItem
	keys            map[int]CacheItem
	Capacity        int
	Currentcapacity int
	Head            *CacheItem
	Tail            *CacheItem
}

func NewLRUCache(capacity int) *LRUCache {
	head := &CacheItem{}
	tail := &CacheItem{}

	head.next = tail
	tail.prev = head

	return &LRUCache{
		Capacity:        capacity,
		Currentcapacity: 0,
		keys:            make(map[int]CacheItem),
		Head:            head,
		Tail:            tail,
	}
}

func (c *LRUCache) Put(node *CacheItem) {
	if c.Currentcapacity == c.Capacity {
	}
}

func (c *LRUCache) Get(key int) string {
	return ""
}

func (c *LRUCache) remove(node *CacheItem) {
	prev := node.prev
	next := node.next

	prev.next = next
	next.prev = prev
}

func (c *LRUCache) addToFront(node *CacheItem) {
	first := c.Head.next

	node.next = first
	node.prev = c.Head

	first.prev = node
	c.Head.next = node
}

// func main() {
// 	cache := NewLRUCache(3)
// 	cache.Put("apple")
// 	cache.Put("banana")
// 	cache.Put("cherry")

// 	fmt.Printf("\n Initial state: %v \n %v \n", cache.Data, cache.keys)

// 	println(cache.Get("banana")) // Output: banana
// 	fmt.Printf("\n Initial state: %v \n %v \n", cache.Data, cache.keys)
// 	println(cache.Get("apple")) // Output: apple
// 	fmt.Printf("\n Initial state: %v \n %v \n", cache.Data, cache.keys)

// 	cache.Put("date") // Evicts "cherry"
// 	fmt.Printf("\n Initial state: %v \n %v \n", cache.Data, cache.keys)
// 	println(cache.Get("cherry")) // Output: Data does not exist
// 	println(cache.Get("date"))   // Output: date
// }
