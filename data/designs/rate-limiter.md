🏗 Design a distributed rate limiter

▸ Requirements
Functional: allow N requests per window per key (user/IP/API key); reject/queue the rest; return limit headers.
Non-functional: low added latency, accurate-ish across many nodes, fail-open vs fail-closed policy.

▸ Estimation
Runs in front of every API request → must add sub-ms overhead at high QPS.

▸ High-level design
At the gateway: on each request, check+increment a counter for the key in a shared store (Redis). Under the limit → allow; over → 429 + Retry-After.

▸ Deep dives
- Algorithm: token bucket (bursts, refill rate) is the common choice; sliding-window counter avoids fixed-window boundary bursts.
- Distributed counter: Redis with an atomic INCR+EXPIRE or a Lua script (check-and-increment atomically) so N gateway nodes share one limit.
- Accuracy vs cost: local per-node counters (fast, approximate) vs central Redis (accurate, a network hop); often a hybrid (local budget + periodic sync).
- Failure mode: if Redis is down, fail-open (allow) or fail-closed (reject) — a business decision.

▸ Trade-offs & bottlenecks
Central store latency/SPOF; fixed-window boundary bursts; clock skew; per-node vs global accuracy; hot keys.

🔗 Further reading
• Stripe — Scaling your API with rate limiters: https://stripe.com/blog/rate-limiters
• ByteByteGo — design a rate limiter (YouTube): https://www.youtube.com/@ByteByteGo
• Arpit Bhayani — rate limiting: https://arpitbhayani.me
