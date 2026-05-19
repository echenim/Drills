Below is a **Go coding drill plan with 20 likely interview problems** tailored for a **traffic infrastructure / backend systems** role like Walmart Traffic Foundation.

## 1. Drill Strategy

For each problem, practice in this order:

1. Restate the problem clearly.
2. Identify the pattern.
3. Write brute force verbally.
4. Implement optimal solution in Go.
5. Walk through edge cases.
6. State time and space complexity.
7. Refactor names and helper functions.

For a 60-minute Go interview, your goal is:

- First 5 minutes: clarify and design.
- Next 25–35 minutes: implement.
- Next 10 minutes: test edge cases.
- Final 5–10 minutes: improve and discuss tradeoffs.

---

# 20 Likely Go Coding Problems

## Section A — Traffic / Infra-Flavored Problems

### 1. Weighted Round Robin Load Balancer

**Why likely:** Traffic infrastructure, request routing, backend selection.

```go
type Server struct {
	Name   string
	Weight int
}

type WeightedRoundRobin struct {
	servers []Server
	index   int
	remain  int
}

func NewWeightedRoundRobin(servers []Server) *WeightedRoundRobin {
	return &WeightedRoundRobin{
		servers: servers,
		index:   -1,
		remain:  0,
	}
}

func (w *WeightedRoundRobin) Next() string {
	if len(w.servers) == 0 {
		return ""
	}

	for {
		if w.remain == 0 {
			w.index = (w.index + 1) % len(w.servers)
			w.remain = w.servers[w.index].Weight
		}

		if w.remain > 0 {
			w.remain--
			return w.servers[w.index].Name
		}
	}
}
```

Test:

```go
servers := []Server{
	{Name: "A", Weight: 2},
	{Name: "B", Weight: 1},
}
lb := NewWeightedRoundRobin(servers)

for i := 0; i < 6; i++ {
	fmt.Println(lb.Next())
}
// A A B A A B
```

---

### 2. Route Matching by Host and Path Prefix

**Pattern:** Trie / prefix matching / sorting by specificity.

```go
type Route struct {
	Host    string
	Prefix  string
	Backend string
}

func MatchRoute(routes []Route, host, path string) string {
	best := ""
	bestLen := -1

	for _, r := range routes {
		if r.Host != host {
			continue
		}

		if strings.HasPrefix(path, r.Prefix) && len(r.Prefix) > bestLen {
			best = r.Backend
			bestLen = len(r.Prefix)
		}
	}

	return best
}
```

Example:

```go
routes := []Route{
	{"api.example.com", "/v1", "backend-v1"},
	{"api.example.com", "/v1/payments", "payments"},
	{"api.example.com", "/", "default"},
}

fmt.Println(MatchRoute(routes, "api.example.com", "/v1/payments/123"))
// payments
```

---

### 3. Fixed Window Rate Limiter

**Pattern:** Hash map + timestamp window.

```go
type FixedWindowLimiter struct {
	limit  int
	window time.Duration
	counts map[string]int
	starts map[string]time.Time
}

func NewFixedWindowLimiter(limit int, window time.Duration) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		limit:  limit,
		window: window,
		counts: make(map[string]int),
		starts: make(map[string]time.Time),
	}
}

func (l *FixedWindowLimiter) Allow(key string, now time.Time) bool {
	start, exists := l.starts[key]
	if !exists || now.Sub(start) >= l.window {
		l.starts[key] = now
		l.counts[key] = 1
		return true
	}

	if l.counts[key] >= l.limit {
		return false
	}

	l.counts[key]++
	return true
}
```

---

### 4. Sliding Window Rate Limiter

**Pattern:** Queue of timestamps per key.

```go
type SlidingWindowLimiter struct {
	limit  int
	window time.Duration
	events map[string][]time.Time
}

func NewSlidingWindowLimiter(limit int, window time.Duration) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		limit:  limit,
		window: window,
		events: make(map[string][]time.Time),
	}
}

func (l *SlidingWindowLimiter) Allow(key string, now time.Time) bool {
	q := l.events[key]
	cutoff := now.Add(-l.window)

	i := 0
	for i < len(q) && q[i].Before(cutoff) {
		i++
	}
	q = q[i:]

	if len(q) >= l.limit {
		l.events[key] = q
		return false
	}

	q = append(q, now)
	l.events[key] = q
	return true
}
```

---

### 5. LRU Cache

**Pattern:** Hash map + doubly linked list.

```go
type LRUCache struct {
	capacity int
	items    map[int]*list.Element
	order    *list.List
}

type entry struct {
	key   int
	value int
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		items:    make(map[int]*list.Element),
		order:    list.New(),
	}
}

func (c *LRUCache) Get(key int) int {
	if elem, ok := c.items[key]; ok {
		c.order.MoveToFront(elem)
		return elem.Value.(entry).value
	}
	return -1
}

func (c *LRUCache) Put(key, value int) {
	if elem, ok := c.items[key]; ok {
		elem.Value = entry{key, value}
		c.order.MoveToFront(elem)
		return
	}

	if c.order.Len() == c.capacity {
		back := c.order.Back()
		old := back.Value.(entry)
		delete(c.items, old.key)
		c.order.Remove(back)
	}

	elem := c.order.PushFront(entry{key, value})
	c.items[key] = elem
}
```

Import:

```go
import "container/list"
```

---

### 6. Consistent Hashing

**Pattern:** Ring + sorted keys + binary search.

```go
type ConsistentHash struct {
	replicas int
	ring     []uint32
	nodes    map[uint32]string
}

func NewConsistentHash(replicas int) *ConsistentHash {
	return &ConsistentHash{
		replicas: replicas,
		nodes:    make(map[uint32]string),
	}
}

func (c *ConsistentHash) AddNode(node string) {
	for i := 0; i < c.replicas; i++ {
		key := hash(fmt.Sprintf("%s-%d", node, i))
		c.ring = append(c.ring, key)
		c.nodes[key] = node
	}
	sort.Slice(c.ring, func(i, j int) bool {
		return c.ring[i] < c.ring[j]
	})
}

func (c *ConsistentHash) GetNode(key string) string {
	if len(c.ring) == 0 {
		return ""
	}

	h := hash(key)
	idx := sort.Search(len(c.ring), func(i int) bool {
		return c.ring[i] >= h
	})

	if idx == len(c.ring) {
		idx = 0
	}

	return c.nodes[c.ring[idx]]
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}
```

Imports:

```go
import (
	"fmt"
	"hash/fnv"
	"sort"
)
```

---

### 7. Top K Busiest Endpoints From Logs

**Pattern:** Hash map + min heap.

```go
type EndpointCount struct {
	Endpoint string
	Count    int
}

type MinHeap []EndpointCount

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].Count < h[j].Count }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *MinHeap) Push(x any) {
	*h = append(*h, x.(EndpointCount))
}

func (h *MinHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func TopKEndpoints(logs []string, k int) []EndpointCount {
	counts := make(map[string]int)

	for _, endpoint := range logs {
		counts[endpoint]++
	}

	h := &MinHeap{}
	heap.Init(h)

	for endpoint, count := range counts {
		heap.Push(h, EndpointCount{endpoint, count})
		if h.Len() > k {
			heap.Pop(h)
		}
	}

	result := make([]EndpointCount, 0, h.Len())
	for h.Len() > 0 {
		result = append(result, heap.Pop(h).(EndpointCount))
	}

	// optional: highest count first
	sort.Slice(result, func(i, j int) bool {
		return result[i].Count > result[j].Count
	})

	return result
}
```

Imports:

```go
import (
	"container/heap"
	"sort"
)
```

---

### 8. Detect Unhealthy Backend From Rolling Errors

**Pattern:** Sliding window.

```go
type HealthChecker struct {
	window    time.Duration
	limit     int
	failures  map[string][]time.Time
}

func NewHealthChecker(limit int, window time.Duration) *HealthChecker {
	return &HealthChecker{
		window:   window,
		limit:    limit,
		failures: make(map[string][]time.Time),
	}
}

func (h *HealthChecker) RecordFailure(server string, now time.Time) bool {
	q := h.failures[server]
	cutoff := now.Add(-h.window)

	i := 0
	for i < len(q) && q[i].Before(cutoff) {
		i++
	}
	q = q[i:]

	q = append(q, now)
	h.failures[server] = q

	return len(q) >= h.limit
}
```

---

### 9. Merge Deployment Windows

**Pattern:** Sort intervals.

```go
type Interval struct {
	Start int
	End   int
}

func MergeIntervals(intervals []Interval) []Interval {
	if len(intervals) == 0 {
		return nil
	}

	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i].Start < intervals[j].Start
	})

	merged := []Interval{intervals[0]}

	for _, curr := range intervals[1:] {
		last := &merged[len(merged)-1]

		if curr.Start <= last.End {
			if curr.End > last.End {
				last.End = curr.End
			}
		} else {
			merged = append(merged, curr)
		}
	}

	return merged
}
```

---

### 10. Shortest Network Path

**Pattern:** BFS.

```go
func ShortestPath(graph map[string][]string, start, target string) int {
	if start == target {
		return 0
	}

	visited := map[string]bool{start: true}
	queue := []string{start}
	distance := 0

	for len(queue) > 0 {
		size := len(queue)
		distance++

		for i := 0; i < size; i++ {
			node := queue[0]
			queue = queue[1:]

			for _, next := range graph[node] {
				if visited[next] {
					continue
				}

				if next == target {
					return distance
				}

				visited[next] = true
				queue = append(queue, next)
			}
		}
	}

	return -1
}
```

---

## Section B — Classic High-Probability Go Interview Problems

### 11. Two Sum

**Pattern:** Hash map.

```go
func TwoSum(nums []int, target int) []int {
	seen := make(map[int]int)

	for i, n := range nums {
		if j, ok := seen[target-n]; ok {
			return []int{j, i}
		}
		seen[n] = i
	}

	return nil
}
```

---

### 12. Best Time to Buy and Sell Stock

**Pattern:** Single pass, track minimum.

```go
func MaxProfit(prices []int) int {
	minPrice := int(^uint(0) >> 1)
	best := 0

	for _, price := range prices {
		if price < minPrice {
			minPrice = price
		}

		profit := price - minPrice
		if profit > best {
			best = profit
		}
	}

	return best
}
```

---

### 13. Longest Substring Without Repeating Characters

**Pattern:** Sliding window.

```go
func LengthOfLongestSubstring(s string) int {
	lastSeen := make(map[byte]int)
	left := 0
	best := 0

	for right := 0; right < len(s); right++ {
		ch := s[right]

		if idx, ok := lastSeen[ch]; ok && idx >= left {
			left = idx + 1
		}

		lastSeen[ch] = right

		if right-left+1 > best {
			best = right - left + 1
		}
	}

	return best
}
```

---

### 14. Valid Parentheses

**Pattern:** Stack.

```go
func IsValid(s string) bool {
	stack := []rune{}
	match := map[rune]rune{
		')': '(',
		']': '[',
		'}': '{',
	}

	for _, ch := range s {
		switch ch {
		case '(', '[', '{':
			stack = append(stack, ch)
		case ')', ']', '}':
			if len(stack) == 0 || stack[len(stack)-1] != match[ch] {
				return false
			}
			stack = stack[:len(stack)-1]
		}
	}

	return len(stack) == 0
}
```

---

### 15. Group Anagrams

**Pattern:** Hash map with normalized key.

```go
func GroupAnagrams(strs []string) [][]string {
	groups := make(map[[26]int][]string)

	for _, s := range strs {
		var key [26]int
		for _, ch := range s {
			key[ch-'a']++
		}
		groups[key] = append(groups[key], s)
	}

	result := make([][]string, 0, len(groups))
	for _, group := range groups {
		result = append(result, group)
	}

	return result
}
```

---

### 16. Kth Largest Element

**Pattern:** Min heap.

```go
type IntMinHeap []int

func (h IntMinHeap) Len() int           { return len(h) }
func (h IntMinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h IntMinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntMinHeap) Push(x any) {
	*h = append(*h, x.(int))
}

func (h *IntMinHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func FindKthLargest(nums []int, k int) int {
	h := &IntMinHeap{}
	heap.Init(h)

	for _, n := range nums {
		heap.Push(h, n)
		if h.Len() > k {
			heap.Pop(h)
		}
	}

	return (*h)[0]
}
```

---

### 17. Product of Array Except Self

**Pattern:** Prefix and suffix product.

```go
func ProductExceptSelf(nums []int) []int {
	n := len(nums)
	result := make([]int, n)

	prefix := 1
	for i := 0; i < n; i++ {
		result[i] = prefix
		prefix *= nums[i]
	}

	suffix := 1
	for i := n - 1; i >= 0; i-- {
		result[i] *= suffix
		suffix *= nums[i]
	}

	return result
}
```

---

### 18. Binary Search

**Pattern:** Sorted array search.

```go
func BinarySearch(nums []int, target int) int {
	left, right := 0, len(nums)-1

	for left <= right {
		mid := left + (right-left)/2

		if nums[mid] == target {
			return mid
		}

		if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	return -1
}
```

---

### 19. Number of Islands

**Pattern:** DFS / BFS grid traversal.

```go
func NumIslands(grid [][]byte) int {
	if len(grid) == 0 {
		return 0
	}

	rows, cols := len(grid), len(grid[0])
	count := 0

	var dfs func(r, c int)
	dfs = func(r, c int) {
		if r < 0 || c < 0 || r >= rows || c >= cols || grid[r][c] != '1' {
			return
		}

		grid[r][c] = '0'

		dfs(r+1, c)
		dfs(r-1, c)
		dfs(r, c+1)
		dfs(r, c-1)
	}

	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if grid[r][c] == '1' {
				count++
				dfs(r, c)
			}
		}
	}

	return count
}
```

---

### 20. Clone Graph

**Pattern:** DFS + map old node to new node.

```go
type Node struct {
	Val       int
	Neighbors []*Node
}

func CloneGraph(node *Node) *Node {
	if node == nil {
		return nil
	}

	seen := make(map[*Node]*Node)

	var dfs func(*Node) *Node
	dfs = func(curr *Node) *Node {
		if clone, ok := seen[curr]; ok {
			return clone
		}

		clone := &Node{Val: curr.Val}
		seen[curr] = clone

		for _, neighbor := range curr.Neighbors {
			clone.Neighbors = append(clone.Neighbors, dfs(neighbor))
		}

		return clone
	}

	return dfs(node)
}
```

---

# Go Templates You Should Memorize

## 1. Main Function Template

```go
package main

import (
	"fmt"
)

func main() {
	fmt.Println("start")
}
```

---

## 2. Hash Map Counting Template

```go
counts := make(map[string]int)

for _, item := range items {
	counts[item]++
}
```

---

## 3. Sliding Window Template

```go
left := 0

for right := 0; right < len(nums); right++ {
	// add nums[right]

	for /* window invalid */ {
		// remove nums[left]
		left++
	}

	// update answer
}
```

---

## 4. BFS Template

```go
func BFS(graph map[int][]int, start int) {
	visited := map[int]bool{start: true}
	queue := []int{start}

	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]

		for _, next := range graph[node] {
			if visited[next] {
				continue
			}

			visited[next] = true
			queue = append(queue, next)
		}
	}
}
```

---

## 5. DFS Template

```go
func DFS(graph map[int][]int, node int, visited map[int]bool) {
	if visited[node] {
		return
	}

	visited[node] = true

	for _, next := range graph[node] {
		DFS(graph, next, visited)
	}
}
```

---

## 6. Heap Template

```go
type MinHeap []int

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[i], h[j] }

func (h *MinHeap) Push(x any) {
	*h = append(*h, x.(int))
}

func (h *MinHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}
```

Small correction:

```go
func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}
```

Use:

```go
h := &MinHeap{}
heap.Init(h)
heap.Push(h, 10)
x := heap.Pop(h).(int)
```

---

## 7. Sort Template

```go
sort.Slice(items, func(i, j int) bool {
	return items[i].Value < items[j].Value
})
```

---

## 8. Two Pointer Template

```go
left, right := 0, len(nums)-1

for left < right {
	sum := nums[left] + nums[right]

	if sum == target {
		return []int{left, right}
	} else if sum < target {
		left++
	} else {
		right--
	}
}
```

---

# 7-Day Drill Plan

## Day 1 — Go Fundamentals + Arrays/Maps

Practice:

1. Two Sum
2. Best Time to Buy and Sell Stock
3. Group Anagrams

Focus:

- Map usage
- Slices
- Clean function signatures
- Edge cases

---

## Day 2 — Sliding Window

Practice:

1. Longest Substring Without Repeating Characters
2. Sliding Window Rate Limiter
3. Detect Unhealthy Backend From Rolling Errors

Focus:

- Window validity
- Shrinking from the left
- Timestamp cleanup

---

## Day 3 — Sorting + Intervals

Practice:

1. Merge Deployment Windows
2. Route Matching
3. Product of Array Except Self

Focus:

- `sort.Slice`
- Prefix/suffix logic
- Choosing most-specific match

---

## Day 4 — Heap / Top K

Practice:

1. Top K Busiest Endpoints
2. Kth Largest Element
3. Heap template from memory

Focus:

- `container/heap`
- Pointer receiver for Push/Pop
- Min heap vs max heap

---

## Day 5 — Graphs

Practice:

1. Shortest Network Path
2. Number of Islands
3. Clone Graph

Focus:

- BFS queue
- DFS recursion
- Visited map

---

## Day 6 — System-Coding Problems

Practice:

1. LRU Cache
2. Fixed Window Rate Limiter
3. Consistent Hashing

Focus:

- Struct design
- Constructor functions
- Encapsulation
- Production-style naming

---

## Day 7 — Mock Interview Day

Do these under a timer:

1. Route Matching — 25 minutes
2. Top K Endpoints — 30 minutes
3. LRU Cache — 40 minutes
4. Shortest Path — 25 minutes

Then verbally explain:

- Complexity
- Edge cases
- Production concerns
- How you would make it thread-safe

---

# Go Interview Talking Points

Use these phrases naturally:

> I’ll start with the simplest correct design, then optimize if needed.

> I’m using a map here because lookup needs to be constant time.

> For production, I would also consider concurrency safety around this structure.

> The main edge cases are empty input, duplicate values, and boundary conditions.

> This solution is O(n) time and O(n) space.

> I’ll keep the code simple and readable first, then we can discuss tradeoffs.

---

# Problems to Prioritize for Walmart Traffic Foundation

Highest priority:

1. Route Matching
2. Rate Limiter
3. LRU Cache
4. Weighted Round Robin
5. Top K Endpoints
6. Consistent Hashing
7. Shortest Path
8. Rolling Health Checker
9. Merge Intervals
10. Longest Substring

These are the most likely to connect directly to **traffic routing, proxy behavior, backend selection, observability, failure detection, and infrastructure ownership**.
