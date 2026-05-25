package main

import (
	"container/heap"
	"fmt"
	"strings"
)

type Item struct {
	Data  string
	Count int
}

type MinHeap []Item

func (h MinHeap) Len() int {
	return len(h)
}

func (h MinHeap) Less(i, j int) bool {
	// Primary sort by count ascending.
	if h[i].Count != h[j].Count {
		return h[i].Count < h[j].Count
	}

	// Tie-breaker: reverse lexicographic order.
	// This makes output deterministic when counts are equal.
	return h[i].Data > h[j].Data
}

func (h MinHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinHeap) Push(x any) {
	*h = append(*h, x.(Item))
}

func (h *MinHeap) Pop() any {
	old := *h
	n := len(old)

	item := old[n-1]
	*h = old[0 : n-1]

	return item
}

func TopKItem(logs []string, k int) []Item {
	if k <= 0 || len(logs) == 0 {
		return nil
	}

	counts := make(map[string]int)

	for _, log := range logs {
		fields := strings.Fields(log)
		if len(fields) < 3 {
			continue // Skip malformed log entries.
		}
		log = fields[2] // Extract the log message (3rd field).
		counts[log]++
	}
	h := &MinHeap{}
	heap.Init(h)

	for log, count := range counts {
		heap.Push(h, Item{Data: log, Count: count})
		if h.Len() > k {
			heap.Pop(h)
		}
	}
	result := make([]Item, h.Len())
	for i := len(result) - 1; i >= 0; i-- {
		result[i] = heap.Pop(h).(Item)
	}
	return result
}

func main() {
	logs := []string{
		"2026-05-23T10:01:00Z GET /api/v1/users 200 34ms",
		"2026-05-23T10:01:01Z POST /api/v1/orders 201 82ms",
		"2026-05-23T10:01:02Z GET /api/v1/users 200 28ms",
		"2026-05-23T10:01:03Z GET /api/v1/products 200 41ms",
		"2026-05-23T10:01:04Z GET /api/v1/users 500 120ms",
	}

	result := TopKItem(logs, 2)

	for _, item := range result {
		fmt.Printf("%s: %d\n", item.Data, item.Count)
	}
}
