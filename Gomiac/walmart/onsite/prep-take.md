Received. This confirms the onsite structure and gives us the exact interview map.

From the recruiter thread, the onsite was triggered after **positive feedback** from the prior code assessment, and Walmart is bringing you onsite for three rounds: **Go coding, domain fit, and hiring manager**.

## Walmart onsite schedule

| Time                                       | Interviewer        | Role                                  | Focus                        |
| ------------------------------------------ | ------------------ | ------------------------------------- | ---------------------------- |
| **Fri 5/29/2026, 9:30 AM – 10:30 AM PST**  | **Brian Merrell**  | Principal Software Engineer           | **Go coding interview**      |
| **Fri 5/29/2026, 10:30 AM – 11:00 AM PST** | **Kiran Reddy**    | Senior Director, Software Engineering | **Hiring manager interview** |
| **Fri 5/29/2026, 11:00 AM – 12:00 PM PST** | **Nathan Brunson** | Director, Software Engineering        | **Domain fit**               |

Location: **Excellence Building, 1345 Crossman Ave, Sunnyvale, CA 94089**

At the entrance, ask for **Madhu Patra**.

## What each round likely means

### 1. Brian Merrell — Principal Software Engineer — Go Coding

This will likely be practical and implementation-heavy. Since Brian is a Principal Engineer, expect him to care about more than just getting the answer accepted.

He will likely evaluate:

- clean Go implementation
- edge cases
- time and space complexity
- how you structure code
- how you think aloud
- whether your solution would survive production-style input
- whether you can debug calmly

For this Traffic Foundation role, I would prepare coding problems around:

- rate limiter
- LRU cache
- consistent hashing
- weighted round robin load balancer
- request router by host/path
- top K endpoints from logs
- rolling error-rate detector
- backend health checker
- concurrency-safe map/cache
- interval merge
- BFS/service dependency traversal

Your goal in this round is not to sound clever. Your goal is to look like a steady engineer who writes correct, readable Go under pressure.

---

### 2. Kiran Reddy — Senior Director — Hiring Manager

This is only 30 minutes, so expect a focused leadership and fit conversation.

He will likely test:

- why Walmart
- why Traffic Foundation
- what level of ownership you have had
- whether you can operate at Staff level
- how you influence teams
- how you handle ambiguity
- how you respond during production issues
- whether you are practical and business-aware

Your strongest message here should be:

> “My background is in production traffic and distributed infrastructure where routing correctness, reliability, rollout safety, and operational visibility matter. I’m comfortable owning systems where mistakes affect real users, and I know how to lead changes carefully across teams.”

Keep this round crisp. Senior Directors usually do not want deep code detail unless they ask for it. They want judgment, maturity, and confidence.

---

### 3. Nathan Brunson — Director — Domain Fit

This is probably the most important round for this role.

The JD specifically asks for experience with **Layer 4 and Layer 7 orchestration/routing/load balancing systems**, **high-performance HTTP proxies and caches**, **TCP/HTTP lifecycle**, systems with **more than 1K servers**, and large-scale distributed systems.

Expect discussion around:

- L4 vs L7 load balancing
- traffic routing
- proxy behavior
- retries and timeouts
- circuit breakers
- TCP connection reuse
- HTTP request lifecycle
- p95/p99 latency
- cache behavior
- production incidents
- canary rollout
- observability
- root cause analysis
- failure isolation

This is where your LinkedIn Traffic Infrastructure experience should be front and center.

Your anchor story should be:

> “I worked on production infrastructure supporting deterministic routing and reliability across large-scale distributed systems. My focus was safe traffic behavior: correctness, routing determinism, failure isolation, rollout safety, and observability.”
