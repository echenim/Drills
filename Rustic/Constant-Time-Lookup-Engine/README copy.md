````markdown
# Design O(1) Lookup Data Structure

A practical hash map design exercise focused on building a data structure that supports **O(1) average-case insertion, lookup, retrieval, and deletion**.

---

## Table of Contents

- [Overview](#overview)
- [Problem Statement](#problem-statement)
- [Goals](#goals)
- [Expected API](#expected-api)
- [Examples](#examples)
- [Constraints](#constraints)
- [Performance Expectations](#performance-expectations)
- [Design Notes](#design-notes)
- [Rust Starter API](#rust-starter-api)
- [Go Starter API](#go-starter-api)
- [Interview Focus Areas](#interview-focus-areas)
- [Follow-Up Questions](#follow-up-questions)
- [Evaluation Rubric](#evaluation-rubric)

---

## Overview

This project explores the design of a fast key-value data structure with **constant average-time access**.

It is a classic interview problem that tests understanding of:

- hash tables
- collision handling
- API design
- resizing and rehashing
- time and space complexity tradeoffs

At its core, this is a simplified **hash map / dictionary** design problem.

---

## Problem Statement

Design and implement a data structure that stores key-value pairs and supports the following operations in **O(1) average time**:

- `put(key, value)` — insert a new key-value pair or update the value of an existing key
- `get(key)` — retrieve the value associated with a key
- `contains(key)` — check whether a key exists
- `remove(key)` — delete a key and its value if present
- `size()` — return the total number of stored entries

You may assume a reasonable hash function is available.

---

## Goals

Your implementation should:

- support fast key-based lookup
- allow updates to existing keys
- handle collisions correctly
- define clear behavior for missing keys
- maintain **O(1) average-case** performance for core operations

---

## Expected API

```text
put(key, value) -> void
get(key) -> value | null
contains(key) -> bool
remove(key) -> bool
size() -> int
```
````

---

## Examples

### Example 1

#### Input

```text
put("a", 10)
put("b", 20)
get("a")
contains("b")
remove("a")
get("a")
size()
```

#### Output

```text
10
true
true
null
1
```

---

### Example 2

#### Input

```text
put("user:1", "William")
put("user:2", "Alex")
put("user:1", "Will")
get("user:1")
contains("user:3")
remove("user:2")
size()
```

#### Output

```text
"Will"
false
true
1
```

---

## Constraints

Suggested baseline constraints:

- `1 <= number_of_operations <= 100000`
- keys may be strings or integers
- values may be strings, integers, or generic objects
- duplicate keys should update existing values
- deleting a missing key should return a safe, predictable result

---

## Performance Expectations

| Operation  | Target Complexity |
| ---------- | ----------------- |
| `put`      | `O(1)` average    |
| `get`      | `O(1)` average    |
| `contains` | `O(1)` average    |
| `remove`   | `O(1)` average    |
| `size`     | `O(1)`            |

### Space Complexity

- `O(n)` where `n` is the number of stored key-value pairs

> Note: true worst-case lookup in a hash table can degrade beyond `O(1)` if collisions are poorly handled.

---

## Design Notes

A strong implementation will usually rely on:

- a **hash table**
- an internal array of buckets
- a collision strategy such as:

  - **separate chaining**
  - **open addressing**

A complete solution should be able to explain:

- how keys are hashed
- how collisions are resolved
- when resizing should occur
- how rehashing preserves performance as the structure grows

---

## Rust Starter API

```rust
pub struct FastMap<K, V> {
    // internal fields
}

impl<K, V> FastMap<K, V> {
    pub fn new() -> Self {
        todo!()
    }

    pub fn put(&mut self, key: K, value: V) {
        todo!()
    }

    pub fn get(&self, key: &K) -> Option<&V> {
        todo!()
    }

    pub fn contains(&self, key: &K) -> bool {
        todo!()
    }

    pub fn remove(&mut self, key: &K) -> bool {
        todo!()
    }

    pub fn size(&self) -> usize {
        todo!()
    }
}
```

---

## Go Starter API

```go
type FastMap struct {
 // internal fields
}

func NewFastMap() *FastMap {
 return &FastMap{}
}

func (m *FastMap) Put(key string, value any) {
 // implement
}

func (m *FastMap) Get(key string) (any, bool) {
 // implement
 return nil, false
}

func (m *FastMap) Contains(key string) bool {
 // implement
 return false
}

func (m *FastMap) Remove(key string) bool {
 // implement
 return false
}

func (m *FastMap) Size() int {
 // implement
 return 0
}
```

---

## Interview Focus Areas

Interviewers typically use this problem to evaluate whether a candidate can:

- design a practical data structure
- reason clearly about complexity
- distinguish **average-case** from **worst-case**
- define behavior for edge cases
- write clean and maintainable code

---

## Follow-Up Questions

Common extensions include:

1. How do you handle collisions?
2. What is the worst-case lookup time?
3. When should the table resize?
4. How would you preserve performance as data grows?
5. How would you make the structure thread-safe?
6. How would you support ordered iteration?
7. How would you add TTL or eviction?
8. How would you optimize memory usage?

---

## Evaluation Rubric

### Strong Solution

- uses hashing correctly
- handles collisions safely
- updates existing keys cleanly
- explains tradeoffs honestly
- produces readable, testable code

### Weak Solution

- assumes true `O(1)` in every case without qualification
- ignores collision handling
- leaves duplicate or missing-key behavior undefined
- mixes interface and implementation carelessly

---

## Key Takeaway

The cleanest description of this problem is:

> Design a data structure that supports **O(1) average-case insertion, lookup, retrieval, and deletion**.

---

## Suggested Title Variants

- **Design O(1) Lookup Data Structure**
- **Constant-Time Key-Value Store**
- **Fast Lookup Map**
- **Implement a Hash-Based Retrieval System**
- **Design a Constant-Time Dictionary**
