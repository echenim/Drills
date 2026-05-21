Yes. Beyond the 15 we already covered, I would expect Nathan may ask **second-layer domain questions** that test whether your experience is real, not just keyword-aligned.

For this Walmart role, I would prepare these additional questions.

## 1. “Describe a traffic system you designed end-to-end.”

They may want to hear the full path:

> client → DNS/CDN/LB → L4/L7 proxy → ingress/service mesh → backend service → downstream dependency

Your answer should cover routing, load balancing, health checks, retries, timeouts, observability, rollback, and ownership.

---

## 2. “What happens when one region becomes unhealthy?”

Strong answer themes:

- detect regional degradation through latency, error rate, saturation, health checks
- avoid instant full failover if it will overload the next region
- shift gradually if possible
- protect downstreams with circuit breakers and retry budgets
- validate recovery before shifting traffic back

This maps well to your StrataLinks regional failover story.

---

## 3. “How do retries cause outages?”

Answer:

> “Retries are extra traffic. If a backend is already overloaded, blind retries multiply load and make recovery harder.”

Mention retry budgets, exponential backoff, jitter, deadlines, idempotency, and avoiding retries on clear overload signals.

---

## 4. “What is the difference between load balancing and traffic routing?”

Good framing:

- **Load balancing** chooses among equivalent backends.
- **Routing** decides the destination class: service, region, version, tenant, cluster, or path.
- Routing answers **where should this request go?**
- Load balancing answers **which healthy instance should handle it?**

---

## 5. “How would you safely roll out a new routing rule?”

Answer structure:

1. validate config offline
2. test route decisions against known examples
3. shadow old vs new route results
4. canary by route, region, or percentage
5. monitor p95/p99, 4xx/5xx, retries, backend distribution
6. rollback quickly if route misses or errors increase

This is directly aligned with Walmart’s expectation around complex cross-functional projects and production RCA.

---

## 6. “How do you debug uneven backend load?”

Strong signals:

- check connection distribution, not just request count
- long-lived TCP connections can pin traffic
- inspect per-backend RPS, active connections, CPU, latency, errors
- verify health checks and outlier ejection
- check whether the load balancer is balancing per connection or per request
- consider keepalive and connection pool behavior

This is a very likely traffic-foundation question.

---

## 7. “What can go wrong with connection pooling?”

Mention:

- stale connections
- too many idle connections
- too few connections causing queueing
- poor reuse causing TCP/TLS handshake overhead
- uneven backend load
- bad timeout settings
- connection leaks

A strong line:

> “Connection pooling is good when it reduces handshake cost, but dangerous when it hides backend imbalance or keeps sending traffic to degraded hosts.”

---

## 8. “How do you design health checks?”

Answer:

- shallow health checks prove process is alive
- deep health checks prove dependencies are usable
- too-deep health checks can cause false negatives
- health checks need thresholds to avoid flapping
- separate readiness from liveness
- use passive health from real traffic plus active probes

---

## 9. “How do you handle overload?”

Mention:

- rate limiting
- backpressure
- bounded queues
- circuit breaking
- load shedding
- priority traffic
- adaptive concurrency
- fast failure instead of slow timeout
- clear 429/503 behavior

Best phrase:

> “During overload, the system should degrade deliberately rather than collapse randomly.”

---

## 10. “What is your approach to root cause analysis?”

Answer with a calm production flow:

1. define the symptom
2. identify blast radius
3. compare before/after
4. isolate layer: proxy, app, dependency, network, region
5. mitigate first
6. preserve evidence
7. write follow-up actions
8. prevent recurrence with guardrails, alerts, and tests

The JD explicitly mentions directing RCA for critical business and production issues.

---

## 11. “How do you think about caches in traffic systems?”

Even if your deepest experience is routing/proxy, prepare this.

Talk about:

- cache hit ratio
- TTL
- invalidation
- stale data risk
- hot keys
- request coalescing
- cache stampede
- negative caching
- stale-while-revalidate
- protecting origin services

Be honest if needed:

> “My deepest hands-on experience is routing, proxy behavior, and ingress policy, but I understand cache behavior because it directly affects backend protection and tail latency.”

---

## 12. “What is the difference between hot path and cold path?”

Answer:

- **Hot path**: executed for most requests; must be fast, predictable, minimal allocation, low latency.
- **Cold path**: less frequent path, such as config load, cache miss, failover, initialization, or error handling.

Strong point:

> “In traffic infrastructure, hot-path changes need much stricter review because every extra allocation, lookup, or network call can affect many services.”

---

## 13. “How would you design observability for a routing platform?”

Mention:

- route match count
- route miss count
- fallback usage
- per-route latency
- per-backend selection
- retry count
- timeout count
- circuit breaker opens
- status codes
- canary vs baseline
- trace fields showing why a request routed where it did

Best line:

> “I want to answer one question quickly: why did this request go to this backend?”

---

## 14. “How do you prevent config drift across languages?”

This maps well to your LinkedIn story.

Answer:

- define invariants
- use canonical schema
- shared test vectors
- compatibility tests
- shadow comparison
- gradual rollout
- avoid multiple teams reimplementing behavior differently
- document edge cases

---

## 15. “What would you improve in an existing traffic platform in your first 90 days?”

Good answer:

> “First, I would understand the current architecture, the operational pain points, and the highest-risk traffic paths. I would look at recent incidents, routing-change process, latency dashboards, retry behavior, failover behavior, and config rollout safety. Then I would identify one or two practical improvements: better route validation, safer canarying, clearer ownership, or stronger observability around route decisions.”

Do not sound like you will rewrite everything.

---

## 16. “How do you balance innovation with stability?”

Answer:

> “For traffic infrastructure, innovation has to be incremental. The platform sits on the critical path, so the right approach is to make small, measurable changes behind guardrails. I’m comfortable pushing architecture forward, but I want validation, canaries, rollback, and clear ownership before broad rollout.”

---

## 17. “How do you influence teams that do not report to you?”

Staff-level question.

Answer:

- build trust through evidence
- write clear design docs
- show failure modes
- align on invariants
- make migration easy
- provide tooling
- avoid forcing unnecessary rewrites
- listen to application team constraints

---

## 18. “Tell me about a time you disagreed with another team technically.”

Use a calm story. Avoid blaming.

Structure:

> “The disagreement was about X. Their concern was Y. My concern was Z. I used data/design review/testing to align us. We chose a path that reduced risk and preserved the business goal.”

---

## 19. “What are the hardest parts of operating traffic infrastructure at Walmart scale?”

Strong answer:

- shared blast radius
- holiday/event traffic spikes
- regional failure
- route correctness
- tail latency
- config safety
- backward compatibility
- application-team self-service without unsafe freedom
- observability across many services

---

## 20. “What is your strongest traffic infrastructure story?”

Have this ready as your closing anchor.

Your best choice is probably:

> **retry storm / regional failover / p99 latency / safe mitigation**

That story shows real production scars, which is exactly what this role needs.

## My best guess

Nathan will likely ask a mix of:

- **L4/L7 routing**
- **load balancing**
- **production incident/RCA**
- **latency debugging**
- **safe rollout**
- **failure isolation**
- **how your past work maps to Walmart**

The biggest thing: do not answer like a textbook. Answer like someone who has watched traffic systems fail in production and knows how to make them boring, safe, and observable.
