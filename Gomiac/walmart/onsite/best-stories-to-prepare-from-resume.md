## Main concern

The resume is strong, but a few claims may invite deep probing.

The biggest one is the repeated **100K–5M+ RPS** and **5M+ RPS** language. If those numbers are real and defensible, keep them. But for onsite, you must be ready to explain:

- where the traffic entered
- whether that was aggregate platform traffic or service-owned traffic
- how it was measured
- what metrics backed it
- what your direct ownership was
- what happened during peak/failover
- what bottleneck you personally helped remove

At Staff level, interviewers will not just accept scale numbers. They will test whether you truly operated that system.

## What I would tighten before the onsite

### 1. Reduce any feeling of “too perfectly tailored”

The resume is very aligned, but it may feel almost too dense with Traffic Foundation keywords. That is not fatal, but in interview you need grounded stories behind every phrase.

For example, be ready for:

> “Tell me about the L4/L7 routing strategy you led.”
> “What exactly did you change in Envoy or Istio?”
> “How did circuit breaking work?”
> “What was the failure mode during the retry storm?”
> “How did you measure sub-100ms p99?”
> “What did you personally own versus what the team owned?”

### 2. Be careful with “Managed systems with >1K servers”

The JD asks for it, and your resume says **>1K nodes** under StrataLinks. Good. But prepare the wording carefully:

> “The platform ran across multi-region Kubernetes clusters with more than 1K nodes. My ownership was the traffic foundation layer: ingress, routing policy, retries, timeouts, circuit breaking, and rollout guardrails.”

That sounds more credible than saying you personally “managed 1K servers.”

### 3. Have one strong HTTP proxy/cache story

The JD explicitly says **“Built high-performance HTTP proxies and caches.”**

Your resume has proxy language, but not much cache depth. If asked, you should be ready with either:

- proxy routing / connection reuse / keepalive story, or
- cache behavior story: TTL, invalidation, hot keys, stale-while-revalidate, request coalescing, cache hit ratio, backend protection

If your real experience is stronger on proxy/routing than cache internals, say that cleanly:

> “My deepest experience is on the proxy/routing and ingress side. I’ve worked around cache behavior and backend protection, but the strongest ownership I bring is traffic routing, overload protection, retries/timeouts, and rollout safety.”

That is better than overclaiming.

## Best stories to prepare from this resume

For the onsite, I would prepare these five stories:

1. **Retry storm / regional failover incident**
   Use for domain fit and hiring manager. Show RCA, mitigation, observability, and prevention.

2. **L4/L7 routing strategy across Envoy/Istio**
   Use for domain fit. Explain routing rules, timeouts, retries, circuit breaking, health checks, and rollout.

3. **LinkedIn routing determinism / ID range evolution**
   Use to show Staff-level platform correctness, cross-language consistency, and safe migration.

4. **p99 latency protection under high RPS**
   Use for performance depth. Explain where latency came from: connection churn, retry amplification, backend saturation, queueing, GC, hot partitions, etc.

5. **Cross-team architecture leadership**
   Use for hiring manager. Show how you influenced teams, wrote docs, aligned standards, and avoided duplicate or unsafe implementations.

## Your positioning statement

Use this as your core interview frame:

> “My background is in traffic and infrastructure systems where request routing, proxy behavior, latency, and failure isolation directly affect production users. I’ve worked across Go, C++, Java, Envoy/Istio, Kubernetes, and cloud-native systems, with a focus on deterministic routing, safe rollouts, observability, and operational reliability. For Walmart Traffic Foundation, that maps directly to building and operating the platform layer that helps application teams route traffic safely and efficiently at very large scale.”

## Final verdict

This resume is **strong for the role**. I would not rewrite it heavily. The work now is to make sure your interview answers sound as grounded as the resume claims.

Next, we should build your **three-round prep plan**:

1. Go coding drill plan for Brian
2. Domain-fit answer bank for Nathan
3. Hiring-manager story bank for Kiran
