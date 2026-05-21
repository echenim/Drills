Below is a **15-question domain-fit answer bank for Nathan Brunson — Director, Software Engineering — Domain Fit**.

For Nathan, the goal is to prove you are not just “backend experienced.” You need to sound like someone who understands **traffic infrastructure as a production discipline**: routing correctness, proxy behavior, load balancing, latency, failure isolation, rollout safety, and operational ownership.

Walmart’s JD is explicit about **L4/L7 orchestration, routing/load balancing, high-performance HTTP proxies and caches, TCP/HTTP lifecycle, cloud/DevOps, and production RCA**. Your resume already maps well to that through LinkedIn Traffic Infrastructure, StrataLinks traffic foundation ownership, Envoy/Istio, Go/C++/Java, Kubernetes, Azure/AWS, p99 latency, retry storms, and regional failover work.

---

# 1. “Walk me through your background and how it fits Traffic Foundation.”

### Strong answer

“My background is in traffic and infrastructure systems where routing behavior, latency, and failure isolation directly affect production users.

At LinkedIn, I worked in Traffic Infrastructure, focused on routing determinism, cross-language traffic library behavior, rollout safety, and correctness across Go, C++, and Java service stacks. Before that, at StrataLinks, I owned traffic foundation architecture for a latency-sensitive real-time platform running across multi-region Kubernetes clusters. That involved ingress strategy, L4/L7 routing, Envoy/Istio policies, retries, timeouts, connection management, circuit breaking, and incident response during retry storms and regional failover events.

So when I look at Walmart Traffic Foundation, I see a very familiar problem space: application teams need a reliable platform layer that helps them route traffic safely, absorb scale, protect backends, and make production changes without introducing traffic regressions.”

### Why this works

It maps your experience directly to the JD without sounding like you are reciting buzzwords.

---

# 2. “What does Traffic Foundation mean to you?”

### Strong answer

“To me, Traffic Foundation is the platform layer that sits between user demand and application services. Its job is to make sure requests reach the right destination quickly, safely, and predictably.

That includes L4 and L7 routing, load balancing, proxy behavior, ingress control, retries, timeouts, circuit breaking, observability, failover, and rollout safety. The platform should hide a lot of traffic complexity from application teams, but it should not hide failure. It should give teams strong defaults, clear visibility, and safe mechanisms to change routing behavior.

The hard part is that traffic systems are shared infrastructure. A small routing or retry mistake can amplify across many services, so correctness, gradual rollout, and operational visibility matter as much as raw performance.”

---

# 3. “Explain L4 vs L7 load balancing.”

### Strong answer

“L4 load balancing operates at the transport layer. It routes based on information like IP, port, TCP connection, and sometimes connection-level health. It is fast and relatively simple, but it does not understand HTTP semantics.

L7 load balancing operates at the application layer. It can route based on host, path, headers, method, cookies, request attributes, or service-specific rules. It gives much more control, but it is also easier to make mistakes because route order, header normalization, retries, and timeouts can change behavior.

In production, I think of L4 as connection distribution and L7 as request-aware routing. For a large traffic platform, you usually need both: L4 to distribute connections efficiently and L7 to make intelligent routing decisions.”

### Follow-up line

“The danger with L7 is accidental broad matching. A route that looks harmless can capture traffic it was never meant to handle.”

---

# 4. “Walk me through the lifecycle of an HTTP request in a large-scale traffic system.”

### Strong answer

“A typical request starts with DNS resolving to an edge or regional entry point. From there it may hit CDN, WAF, or an external load balancer. Then it reaches an L4 or L7 proxy layer where TLS may terminate, routing rules are applied, headers are normalized, and policies like retries, timeouts, auth, rate limits, and circuit breakers may execute.

After that, the request is routed to a regional service endpoint, Kubernetes ingress, service mesh sidecar, or directly to a backend service. The backend may call downstream services or databases. The response then flows back through the same proxy path, where metrics, traces, access logs, and response codes are captured.

The important thing is to understand where latency and failure can be introduced: DNS, connection setup, TLS handshake, queueing at the proxy, poor upstream connection reuse, retry amplification, backend saturation, or downstream dependency delay.”

---

# 5. “Tell me about a production traffic incident you handled.”

### Strong answer

“At StrataLinks, we had a high-severity event during a traffic spike where retries and regional failover behavior started amplifying load instead of protecting the system. The first symptom was p95 staying somewhat manageable while p99 started bending upward, which told us the long tail was under pressure.

We separated proxy latency from application latency, looked at request volume, upstream errors, retry counts, connection churn, queueing, and regional distribution. The issue was not one bad service alone; it was the interaction between retry policy, backend saturation, and failover behavior.

The immediate mitigation was to reduce retry amplification, tighten timeouts, isolate unhealthy paths, and shift traffic more conservatively. The follow-up was better guardrails: clearer retry budgets, circuit-breaking thresholds, better dashboarding, and safer failover behavior.

The lesson for me was that traffic systems need to fail deliberately. If every layer retries blindly, the platform can multiply the original failure.”

---

# 6. “How do you design retry and timeout policies?”

### Strong answer

“I start by treating retries as a limited budget, not a default reflex. A retry is useful only when the failure is likely transient and the backend can absorb the extra work.

For timeouts, I work backward from the user-facing latency budget. If the overall request budget is, say, 200ms, each hop needs a smaller budget. You cannot allow every downstream call to consume the full request timeout. That creates tail latency and queue buildup.

For retries, I look at method safety, idempotency, error type, retry count, jitter, backoff, and whether the request is already near its deadline. I also avoid retrying on failures that indicate overload unless there is a clear alternate backend or region.

The production rule is: retries should improve availability without hiding systemic failure or amplifying load.”

### Good phrase

“Retries are useful medicine, but in traffic infrastructure the wrong dose becomes poison.”

---

# 7. “How would you design circuit breaking for a traffic platform?”

### Strong answer

“I would design it as a state machine with closed, open, and half-open states.

In the closed state, traffic flows normally while we track rolling failure rate, timeout rate, and maybe latency. If failures exceed a threshold over a meaningful window, the circuit opens and we fail fast instead of continuing to send traffic to a failing dependency.

After a cooldown, we move to half-open and allow a limited number of probe requests. If probes succeed, we close the circuit gradually. If they fail, we reopen.

For a traffic platform, the details matter: the circuit should be per backend or per route, not too global. It should use rolling windows, avoid flapping, expose metrics, and integrate with load balancing so unhealthy endpoints are not selected. I would also separate local process-level circuit breaking from fleet-level health decisions, because one instance may see a different failure pattern than the whole fleet.”

---

# 8. “How do you think about p95 and p99 latency?”

### Strong answer

“p50 tells me the normal path. p95 tells me whether a meaningful minority of users are starting to feel pain. p99 tells me where the system breaks under stress.

When p50 is stable but p95 creeps, I look for growing contention, queueing, uneven load, slow downstream dependencies, or partial saturation. When p99 bends upward sharply, I look for retry amplification, connection churn, cold paths, overloaded instances, GC pauses, bad shards, or regional imbalance.

For traffic infrastructure, I care more about tail latency than average latency because proxies and routing layers sit on the hot path. A small amount of added latency at this layer affects many services.”

---

# 9. “How would you debug a tail-latency regression in a proxy layer?”

### Strong answer

“I would first isolate where the latency is coming from: client to proxy, inside proxy, proxy to upstream, upstream service, or downstream dependency.

Then I would compare before and after across p50, p95, p99, error rate, retry count, connection reuse, upstream connection churn, TLS handshakes, queue depth, CPU, memory, GC, and per-backend distribution. I would also check whether the regression is global, regional, route-specific, or tied to a subset of backends.

If p99 moved but p50 stayed flat, I would suspect uneven load, connection pool issues, retries, overloaded backends, or a small number of slow paths. If all percentiles moved, I would suspect a broader configuration or capacity issue.

The key is not to guess from one graph. You want to narrow the blast radius and separate proxy behavior from application behavior.”

---

# 10. “Tell me about routing correctness.”

### Strong answer

“Routing correctness means the request consistently reaches the intended destination based on the agreed contract: host, path, headers, identity, region, service version, or policy.

At LinkedIn, one area I worked on was routing determinism across heterogeneous service stacks. The challenge was making sure routing behavior stayed consistent across Go, C++, and Java paths as identity and routing logic evolved. The risk in these systems is correctness drift, where two implementations look equivalent but behave differently under edge cases.

My approach is to define clear invariants, normalize inputs, test edge cases, compare behavior across implementations, and roll out gradually with canaries and observability. In traffic infrastructure, an incorrect route can be worse than a failed request because it may silently send traffic to the wrong place.”

---

# 11. “How would you design a request router?”

### Strong answer

“I would start with the matching model. For L7 routing, common dimensions are host, path prefix, exact path, method, headers, and sometimes tenant or region.

The router needs deterministic precedence rules. For example, exact host before wildcard host, longest path prefix before shorter prefix, explicit header match before default route, and a safe fallback for no match.

I would store routes in a structure that supports fast lookup and predictable ordering. For simple cases, sorted route lists by specificity may be enough. At larger scale, a trie for path prefixes or compiled route tables may be better.

Operationally, I would add validation before deployment to catch ambiguous routes, shadow evaluation to compare old and new route decisions, and metrics showing route match counts, misses, and fallback usage.”

---

# 12. “How do you prevent bad config or routing changes from taking down production?”

### Strong answer

“I use layered safety.

First, validate config before it reaches production: schema validation, route conflict checks, ownership checks, and test cases for expected route decisions.

Second, roll out gradually: dev, staging, canary, small production slice, then wider rollout. For traffic systems, I prefer canarying by region, route, or percentage of traffic.

Third, observe the right metrics: route match rate, 4xx/5xx, p95/p99, retry volume, backend distribution, connection churn, and fallback behavior.

Fourth, have fast rollback. A traffic change should be reversible without waiting for a full application deployment.

The principle is old-fashioned but still the best one: don’t make a global traffic change until the small version has proven itself.”

---

# 13. “How do you think about load balancing strategy?”

### Strong answer

“It depends on what you are optimizing for.

Round robin is simple and works when backends are homogeneous. Least connections can help when request duration varies. Weighted balancing helps when backends have different capacities. Consistent hashing is useful when you want request affinity, cache locality, or shard stability. Latency-aware or outlier-aware balancing can help avoid slow backends, but it must be designed carefully to avoid overreacting to noisy signals.

In production, I care about fairness, backend health, connection reuse, load skew, and failure behavior. The best algorithm on paper can still behave badly if long-lived connections pin too much traffic to a small set of instances or if health checks lag reality.”

---

# 14. “What metrics would you put on a traffic foundation platform?”

### Strong answer

“I would separate platform-level, route-level, backend-level, and dependency-level metrics.

At platform level: total RPS, error rate, saturation, CPU, memory, open connections, connection churn, TLS handshakes, queue depth, and proxy latency.

At route level: requests by route, p50/p95/p99, status codes, retries, timeouts, circuit breaker opens, rate-limit rejects, fallback usage, and route misses.

At backend level: backend selection count, active connections, health status, outlier ejection, per-backend latency, and per-backend errors.

For rollout safety: compare canary versus baseline on latency, error rate, retry rate, and traffic distribution. I also want logs and traces that let me answer: why did this request go to this backend?”

---

# 15. “Why Walmart Traffic Foundation?”

### Strong answer

“What interests me is the scale and the centrality of the platform. Walmart’s traffic foundation is not a side system. It affects e-commerce, stores, distribution, and customer-facing experiences. The JD talks about millions of requests per second, high-traffic situations, and helping application teams reach customers in the fastest and most efficient way possible. That is exactly the kind of infrastructure layer I enjoy working on.

My strongest experience is in production traffic systems: routing behavior, proxy policies, reliability, failover, p99 latency, and safe rollouts. I like systems where engineering discipline matters because the cost of mistakes is real. Walmart has the scale where small improvements in traffic reliability and efficiency can have very large impact.”

---

# Nathan round strategy

For this interview, keep coming back to these five themes:

1. **Correctness before cleverness**
   Bad routing is worse than no routing.

2. **Retries and failover must be controlled**
   Otherwise the platform amplifies failure.

3. **Tail latency matters more than averages**
   Traffic layers sit on the hot path.

4. **Safe rollout is part of the design**
   Canary, metrics, rollback, and validation are not afterthoughts.

5. **Traffic infrastructure is shared infrastructure**
   Your job is to make the safe path the easy path for application teams.

## Best closing line for Nathan

> “The way I think about traffic foundation is simple: route correctly, fail safely, observe clearly, and roll out gradually. At scale, those basics matter more than fancy architecture.”
