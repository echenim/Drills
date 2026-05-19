Here is the same content organized into cleaner sections:

---

# Walmart Onsite Interview Preparation Plan

## 1. Onsite Overview

The recruiter update is very good news. Walmart confirmed positive feedback from your code assessment and is moving you to an in-person onsite with three rounds:

1. **Go coding interview — 60 minutes**
2. **Domain fit — 60 minutes**
3. **Hiring manager — 30 minutes**

They also asked whether you are available to come onsite Friday, **May 22**.

This means you have already cleared the first technical screen. The onsite is now about confirming three things:

- You can write clean Go under time pressure.
- You deeply understand traffic infrastructure.
- You can operate at Staff level with ownership, judgment, communication, leadership, and production accountability.

---

# Section 1: Go Coding Interview — 60 Minutes

## What to Expect

Since the onsite includes a dedicated Go coding interview and they explicitly told you to bring your laptop, expect live coding in Go.

This round will likely be practical problem solving, not pure academic puzzles. Walmart will look for clean implementation, edge-case handling, and clear communication.

## What They Will Evaluate

- Clean, readable Go
- Good use of maps, slices, structs, sorting, heaps, and queues
- Correct edge-case handling
- Clear communication while coding
- Simple production-style organization
- Ability to explain complexity and tradeoffs
- Practical judgment, not cleverness

Your goal is not fancy code. Your goal is **boring, correct, readable Go**.

## Go Coding Style to Use

- Define simple structs.
- Keep state explicit.
- Avoid clever one-liners.
- Talk through edge cases before coding.
- Write small helper functions.
- Test with at least two or three examples out loud.

## Category A: Standard Coding Patterns

Be fast and clean with:

- Hash map counting
- Two pointers
- Sliding window
- Sorting + greedy
- Heap / priority queue
- BFS / DFS
- Interval merging
- Stack
- Prefix sum

## Category B: Traffic-Infrastructure Flavored Coding

These are more likely for this team:

- LRU cache
- Rate limiter
- Weighted round-robin load balancer
- Consistent hashing
- Top K endpoints from logs
- Route matcher by host/path prefix
- Health-check based backend picker
- Request counter over rolling time window
- Circuit breaker state machine
- Worker pool / bounded queue
- Detect unhealthy backend based on rolling errors
- Merge intervals for deploy windows
- Shortest path / BFS for network hops

## Strong Opening Line

> I’ll first clarify the input size and edge cases, then I’ll aim for a simple correct solution. Once that is working, I’ll discuss complexity and any production considerations.

That line sets the right tone: calm, structured, and production-minded.

---

# Section 2: Domain Fit Interview — 60 Minutes

## Why This Round Matters

This is probably the most important round for you.

Walmart Traffic Foundation builds foundational traffic technology that helps application teams reach customers quickly and efficiently, handles millions of requests per second, scales systems during extremely high traffic, and distributes traffic to the right places.

They will test whether your experience is real in traffic infrastructure.

## Topics to Speak Fluently About

- L4 vs. L7 routing
- TCP connection lifecycle
- HTTP request lifecycle
- Load balancing policies
- Connection pooling and keepalive
- Retries, timeouts, backoff, and circuit breaking
- Proxy and cache behavior
- Canarying traffic changes
- Observability: p50, p95, p99, error rate, saturation, request volume, connection churn
- Incident RCA
- Scaling under traffic spikes

## Your Core Domain-Fit Message

> My strongest experience is close to the traffic path: routing correctness, proxy behavior, ingress control, load balancing, connection management, retries/timeouts, and safe rollouts. I’ve worked in environments where the cost of a bad traffic change is high, so I focus heavily on observability, deterministic behavior, staged rollout, and fast rollback.

You should sound like someone who has operated real systems, debugged real failures, and understands what happens when traffic systems behave badly in production.

---

# Section 3: Production Stories to Prepare

## Story 1: L4 Load Balancing Skew

### When to Use It

Use this when asked about:

- Load balancing
- Production debugging
- Traffic imbalance
- Backend saturation
- TCP behavior

### Story Structure

> We saw uneven backend utilization even though request volume looked balanced at a higher level. The root cause was long-lived TCP connections causing load to stick to a subset of backends. I broke the problem down by looking at active connections, request volume, CPU/memory, p99 latency, and backend saturation. The fix was not just changing an algorithm blindly; we tuned balancing behavior, health-check sensitivity, and connection lifecycle behavior, then rolled it out through canary validation.

### Key Terms to Use

- TCP connection lifetime
- Connection reuse
- Backend saturation
- Least request / round robin / consistent hashing tradeoff
- Health checks
- Canary
- Rollback criteria

---

## Story 2: L7 Misrouting

### When to Use It

Use this when asked about:

- Routing correctness
- Request classification
- Route matching
- Blast radius
- Config validation

### Story Structure

> We had a broad L7 rule matching more traffic than intended. The issue came down to route precedence across host, path, and header conditions. I isolated the affected request classes, compared expected versus actual route decisions, tightened the matcher specificity, reordered rules, and validated using shadow traffic and canary rollout before expanding the change.

### Key Terms to Use

- Host/path/header precedence
- Route specificity
- Deterministic routing
- Request classification
- Config validation
- Blast radius
- Canary validation

---

## Story 3: Proxy Connection Reuse and p99 Latency

### When to Use It

Use this when asked about:

- Latency
- Proxy behavior
- Service performance
- Tail latency
- Connection pooling

### Story Structure

> We saw elevated p99 latency, but app-level latency alone did not explain it. I separated proxy overhead from service latency and found excessive upstream connection churn, which meant more TCP/TLS handshakes and worse tail latency. We tuned keepalive, idle timeout, pool limits, and retry behavior, then tracked connection reuse ratio, handshake rate, upstream latency, and error rates.

### Key Terms to Use

- Upstream connection pool
- Keepalive
- Idle timeout
- TCP/TLS handshake
- p99 latency
- Retry budget
- Proxy latency vs. application latency
- Connection churn

---

## Story 4: Brooklin Migration / Event Consumption Reliability

### When to Use It

Use this when asked about:

- Cross-team project leadership
- Platform migration
- Reliability
- Preserving behavior during system change
- Rollout safety

### Story Structure

> We migrated a production event-consumption path from legacy Databus behavior to Brooklin. The core challenge was preserving business behavior while changing the consumption model from callback-style handling to a poll-loop model. I focused on compatibility, configuration correctness across environments, observability, canary rollout, and avoiding silent failure modes.

### Key Terms to Use

- Migration safety
- Compatibility
- Environment-specific configuration
- Rollout gates
- Observability
- Failure visibility
- Silent failure prevention
- Canary rollout

---

## Story 5: Replicationdelay Service Failure Semantics

### When to Use It

Use this when asked about:

- Judgment
- Production correctness
- Quality
- API contracts
- Caller safety
- Ambiguous system behavior

### Story Structure

> The service returned HTTP 200 with empty results for multiple different failure modes. That forced callers to guess whether there was truly no data or whether the system had failed. I pushed the boundary behavior to be explicit: real failures should surface as failures, while true unknown states should remain distinguishable. That made downstream behavior safer and easier to reason about.

### Key Terms to Use

- Failure semantics
- API contract
- Caller correctness
- Production debugging
- Explicit failure over silent ambiguity
- Safer downstream behavior

---

# Section 4: Hiring Manager Interview — 30 Minutes

## What to Expect

This round is short, so your answers need to be concise and Staff-level.

Walmart’s JD emphasizes technical vision, roadmap influence, discovery for large projects, cross-functional engineering leadership, RCA for production issues, and communication with decision-makers.

This is where they will test maturity, not just technical depth.

## Likely Questions

- Why Walmart?
- Why Traffic Foundation?
- How do you lead without authority?
- How do you handle disagreement?
- Tell me about a production issue you owned.
- Tell me about a cross-team project.
- Tell me about a time you improved reliability.
- How do you mentor or influence engineers?
- How do you handle ambiguous production issues?
- How do you make tradeoffs between speed, reliability, and simplicity?
- What kind of environment helps you do your best work?

## Manager-Round Theme

> I’m at my best in infrastructure roles where the work has real production consequences: traffic correctness, reliability, latency, and safe rollout. I like building systems that application teams can trust without needing to understand every low-level networking detail.

---

# Section 5: Prepared Hiring Manager Answers

## Why Walmart?

> Walmart is one of the few places where traffic infrastructure has direct real-world impact at massive scale — e-commerce, stores, fulfillment, and customer-facing systems. The Traffic Foundation role is especially interesting because it sits at the layer where routing, load balancing, resilience, and developer productivity meet. That is exactly the kind of infrastructure work I’ve been doing and want to keep growing in.

## Why This Team?

> Traffic Foundation is close to the request path. It requires deep systems fundamentals — TCP, HTTP, proxies, routing, load balancing — but also strong product thinking because internal application teams depend on the platform. That mix of low-level engineering and broad platform impact is what attracts me.

## How Do You Lead Without Authority?

> I try to make the problem concrete first: what is broken, what is the risk, what are the options, and what evidence supports the decision. Once teams agree on the failure mode and success criteria, alignment becomes easier. I’ve found that clear docs, small rollout steps, and objective metrics usually reduce opinion-based disagreement.

## Tell Me About a Production Issue

Use one of these three, depending on the interviewer’s angle:

- **Load balancing / traffic imbalance:** L4 load balancing skew from long-lived TCP connections
- **Routing correctness:** L7 misrouting caused by broad route matching
- **Latency:** Proxy connection reuse causing elevated p99 latency

---

# Section 6: Resume Risk Areas to Prepare For

Your resume is strong, but these areas may invite pressure testing.

## Risk Area 1: The 5M+ RPS Claim

You list **100K–5M+ RPS** at StrataLinks and **5M+ RPS** at Dell. Be ready to explain:

- Was this per service, per platform, per region, or aggregate?
- Was it request-per-second, event-per-second, packet-per-second, or bid-request volume?
- How was it measured?
- What telemetry showed it?
- What was the peak versus sustained number?
- What was your direct ownership?

### Strong Answer

> That number was aggregate platform traffic at peak. My direct work was on the routing/ingress policies and backend request path behavior that supported that traffic, not every individual service behind it.

## Risk Area 2: “Managed Systems With >1K Servers”

The JD asks for this. Your resume says more than 1K nodes at StrataLinks. Be ready to explain:

- Kubernetes clusters?
- Nodes or servers?
- Multi-region?
- What did you personally manage?
- What were the operational controls?
- What did you own directly versus support indirectly?

## Risk Area 3: “Staff” Title vs. Current LinkedIn Title

Your resume profile says **Staff Traffic & Infrastructure Engineer**, while the LinkedIn role section says **Senior Software Engineer**. That is okay if positioned as role scope, but do not let it sound inflated.

### Strong Answer

> My title has been Senior Software Engineer, but the scope has included Staff-level infrastructure work: cross-team design, production reliability, rollout strategy, and technical alignment across service owners.

## Risk Area 4: Azure / Go Preferred

The JD says cloud, with Azure and Go preferred. Your resume has Go primary, Azure under Dell, and AWS elsewhere. That is enough.

### Strong Answer

> Go is my primary systems language, and I’ve worked across Kubernetes-based cloud environments including Azure and AWS.

---

# Section 7: Core Interview Positioning

## Interview Positioning Statement

Use this early in the domain or hiring manager round:

> My background is in production traffic infrastructure and distributed systems. I’ve worked on L4/L7 routing, ingress, proxy behavior, retries, timeouts, connection management, and rollout safety across Go, C++, and Java systems. The common thread in my work is making request paths deterministic, observable, and resilient under high traffic. I’m strongest in environments where correctness and reliability matter as much as raw throughput.

## Overall Walmart Signal

Your message across all three rounds should be consistent:

> I am a production infrastructure engineer with strong experience in distributed systems, traffic routing, reliability, and safe rollout of critical changes. I care about correctness, determinism, observability, and operational ownership because traffic infrastructure leaves very little room for vague behavior. My strength is taking complex infrastructure problems, reducing them to first principles, and driving practical, production-safe solutions.

That is the right Walmart signal: grounded, senior, practical, and credible.
