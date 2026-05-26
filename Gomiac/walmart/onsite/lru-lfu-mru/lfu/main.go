package main

type CacheItem struct {
	value string
	key   int
	prev  *CacheItem
	next  *CacheItem
}
type LFUCache struct {
	Capacity    int
	FreqCounter int
	Item        map[int]*CacheItem
	FreqList    map[int]int
	Head        *CacheItem
	Tail        *CacheItem
}

func Constructor(capacity int) *LFUCache {
	head := &CacheItem{}
	tail := &CacheItem{}

	head.next = tail
	tail.prev = head

	return &LFUCache{
		Capacity: capacity,
		Item:     make(map[int]*CacheItem),
		FreqList: make(map[int]int),
		Head:     head,
		Tail:     tail,
	}
}

func (this *LFUCache) Get(key int) int {
}

func (c *LFUCache) remove(node *CacheItem) {
	prev := node.prev
	next := node.next

	prev.next = next
	next.prev = prev
}
