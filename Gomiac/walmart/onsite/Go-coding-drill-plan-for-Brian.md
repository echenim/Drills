Below is a **15-problem Go coding drill plan for Brian Merrell’s 60-minute coding round**. I would bias toward **traffic-infrastructure flavored problems**, because Brian is a Principal Engineer and the role is Traffic Foundation.

## How to practice each problem

For every drill, use this routine:

**5 minutes** — clarify requirements and edge cases
**30–35 minutes** — implement clean Go
**10 minutes** — test with examples and edge cases
**5 minutes** — explain complexity and production tradeoffs

Your goal is not just to solve. Your goal is to sound like someone who writes production Go.

---

# 15 Go Coding Drills for Walmart Traffic Foundation

## 1. Implement a Rate Limiter

**Interview framing:**

> Design a rate limiter that allows at most N requests per user within a time window.

**What Brian may test:**

- maps
- timestamps
- cleanup
- edge cases
- concurrency follow-up

**Expected approach:**

Use fixed window or sliding window. For interview speed, start with fixed window, then discuss sliding window improvement.

**Go concepts:**

```go
map[string][]time.Time
sync.Mutex
time.Now()
```

**Follow-up:**

How would you make this safe under high concurrency?

---

## 2. Implement LRU Cache

**Interview framing:**

> Build an LRU cache with Get and Put in O(1).

**What Brian may test:**

- hashmap + doubly linked list
- pointer correctness
- edge cases
- clean struct design

**Expected approach:**

Use `container/list` plus map.

**Go concepts:**

```go
container/list
map[int]*list.Element
```

**Production angle:**

LRU is common in proxy/cache systems. Be ready to discuss TTL, memory pressure, and eviction metrics.

---

## 3. Weighted Round-Robin Load Balancer

**Interview framing:**

> Given a list of backends with weights, return the next backend for each request.

**What Brian may test:**

- load balancing basics
- stateful iteration
- fairness
- clean implementation

**Example:**

```text
A weight 3
B weight 1
C weight 2
```

Expected rough distribution: A gets 3/6, B gets 1/6, C gets 2/6.

**Production follow-up:**

What happens when a backend is unhealthy?

---

## 4. Consistent Hashing

**Interview framing:**

> Implement consistent hashing for routing keys to backend servers.

**What Brian may test:**

- sorted ring
- binary search
- hash function
- minimal remapping

**Go concepts:**

```go
sort.Search
hash/fnv
map[int]string
```

**Production angle:**

This maps very well to traffic routing and cache sharding.

---

## 5. Request Router by Host and Path

**Interview framing:**

> Build a router that maps requests to services based on host and path prefix.

**Example:**

```text
api.walmart.com /cart     -> cart-service
api.walmart.com /checkout -> checkout-service
```

**What Brian may test:**

- prefix matching
- longest-prefix wins
- clean route table
- edge cases

**Important rule:**

More specific routes should win over broad routes.

**Production angle:**

This is directly relevant to L7 routing.

---

## 6. Top K Busiest Endpoints from Logs

**Interview framing:**

> Given request logs, return the top K endpoints by request count.

**What Brian may test:**

- hashmap counting
- heap
- sorting tradeoffs
- memory usage

**Expected approach:**

Use map for counts, then min-heap of size K.

**Go concepts:**

```go
container/heap
map[string]int
```

**Production angle:**

This feels like traffic observability.

---

## 7. Rolling Error-Rate Detector

**Interview framing:**

> Track success/failure events and alert when the error rate exceeds a threshold over the last N requests or last T seconds.

**What Brian may test:**

- sliding window
- queue behavior
- counters
- time-based cleanup

**Production angle:**

This maps to circuit breakers, health checks, and SLO alerting.

---

## 8. Backend Health Checker

**Interview framing:**

> Given backend health events, choose a healthy backend for each request.

**Requirements:**

- mark backend unhealthy after N failures
- recover after success or cooldown
- never route to unhealthy nodes unless all are unhealthy

**What Brian may test:**

- state transitions
- maps
- edge cases
- production thinking

**Production angle:**

This connects directly to load balancing and failover.

---

## 9. Circuit Breaker

**Interview framing:**

> Implement a circuit breaker with Closed, Open, and Half-Open states.

**What Brian may test:**

- state machine
- failure threshold
- timeout
- half-open probe
- concurrency safety

**States:**

```text
Closed    -> normal traffic
Open      -> reject fast
HalfOpen  -> allow limited probe
```

**Production angle:**

Excellent for traffic infrastructure. You should be very comfortable with this one.

---

## 10. Merge Deployment Windows

**Interview framing:**

> Given a list of deployment windows, merge overlapping intervals.

**Example:**

```text
[1,3], [2,6], [8,10] -> [1,6], [8,10]
```

**What Brian may test:**

- sorting
- interval logic
- clean edge cases

**Production angle:**

Can be framed around rollout windows, maintenance windows, or traffic drain periods.

---

## 11. Find Service Dependency Path

**Interview framing:**

> Given service dependencies, find whether service A can reach service B.

**Example:**

```text
frontend -> api
api -> cart
cart -> pricing
```

Can `frontend` reach `pricing`?

**What Brian may test:**

- graph traversal
- BFS/DFS
- visited set
- cycles

**Production angle:**

Useful for dependency graphs, blast-radius analysis, and routing dependencies.

---

## 12. Detect Cycle in Service Dependency Graph

**Interview framing:**

> Given service dependencies, detect if there is a cycle.

**What Brian may test:**

- DFS states
- recursion or iterative stack
- graph correctness

**States:**

```text
unvisited
visiting
visited
```

**Production angle:**

Service dependency cycles are dangerous for startup ordering, retries, cascading failure, and deployment planning.

---

## 13. Parse Logs and Compute p95 Latency

**Interview framing:**

> Given request logs with endpoint and latency, compute p95 latency per endpoint.

**What Brian may test:**

- grouping
- sorting
- percentile calculation
- edge cases

**Example input:**

```text
/cart 12
/cart 20
/checkout 40
/cart 100
```

**Production angle:**

This is highly relevant because Traffic Foundation cares about p95/p99 behavior, not just averages.

---

## 14. Bounded Worker Pool

**Interview framing:**

> Implement a worker pool that processes jobs with limited concurrency.

**What Brian may test:**

- goroutines
- channels
- wait groups
- graceful shutdown
- error handling

**Go concepts:**

```go
chan Job
sync.WaitGroup
context.Context
```

**Production angle:**

This tests whether your Go concurrency is practical and safe.

---

## 15. Concurrent In-Memory Counter

**Interview framing:**

> Implement a thread-safe request counter per endpoint.

**Requirements:**

- increment endpoint count
- get count
- return top endpoint
- safe for concurrent access

**What Brian may test:**

- `sync.Mutex`
- `sync.RWMutex`
- atomics follow-up
- map concurrency safety

**Production angle:**

This maps to traffic metrics, rate tracking, and request accounting.

---

# Priority Order

Given the Walmart role, I would prioritize them this way:

## Must master before interview

1. Rate limiter
2. LRU cache
3. Weighted round-robin load balancer
4. Request router by host/path
5. Circuit breaker
6. Rolling error-rate detector
7. Top K endpoints from logs
8. Bounded worker pool

## Strong second layer

1. Consistent hashing
2. Backend health checker
3. p95 latency from logs
4. Service dependency BFS
5. Cycle detection
6. Merge deployment windows
7. Concurrent request counter

---

# What to say before coding

Use this structure:

> “Let me first clarify the behavior and edge cases. Then I’ll start with a simple correct design, implement it cleanly in Go, and after that I’ll discuss how I would make it production-ready.”

That sounds mature and gives Brian confidence you are not rushing.

---

# What Brian will likely reward

He will likely care about:

- simple, readable Go
- correct edge-case handling
- clear struct design
- no over-engineering
- good test examples
- calm debugging
- complexity explanation
- production awareness

For this role, every answer should have a small production comment at the end:

> “In production, I’d add metrics around rejected requests, latency, backend selection count, and error rates. I’d also make this concurrency-safe and test it under race detector.”

That is the Walmart Traffic Foundation signal.

| #   | Problem Area                   | Likelihood      |
| --- | ------------------------------ | --------------- |
| 1   | Log parsing / top K / counting | High            |
| 2   | Rate limiter / sliding window  | High            |
| 3   | LRU cache                      | Medium-high     |
| 4   | Request router / prefix match  | Medium-high     |
| 5   | Load balancer / backend picker | Medium-high     |
| 6   | Rolling error-rate detector    | Medium-high     |
| 7   | Circuit breaker                | Medium          |
| 8   | Worker pool / Go concurrency   | Medium          |
| 9   | Consistent hashing             | Medium          |
| 10  | p95/p99 latency from logs      | Medium          |
| 11  | Backend health checker         | Medium          |
| 12  | Service dependency BFS         | Medium-low      |
| 13  | Cycle detection                | Medium-low      |
| 14  | Merge intervals                | Low-medium      |
| 15  | Generic array/string problem   | Always possible |

## My honest bet

If I had to bet on Brian’s round, I would expect one of these five:

1. **Rate limiter**
2. **Top K from logs**
3. **Request router by path/host**
4. **LRU cache**
5. **Load balancer / backend selection**

Those are the best overlap between **Go coding**, **60-minute interview**, and **Walmart Traffic Foundation**.

Your prep should not be “memorize 15.” It should be: master those five, then be comfortable adapting the patterns.
