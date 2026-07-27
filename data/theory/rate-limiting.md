📖 Rate limiting & backpressure (system design)

Cap the request rate to protect a service and share it fairly.

▸ Algorithms
- Token bucket: tokens refill at a fixed rate; a request spends one. Allows bursts up to the bucket size. Most common.
- Leaky bucket: requests drain at a constant rate → smooths bursts.
- Fixed/sliding window counters: simple, but fixed windows allow 2x bursts at the boundary; sliding-window log/counter fixes that.

▸ Distributed
A shared counter (Redis INCR + expiry, or a Lua script) coordinates limits across instances — trade accuracy vs coordination cost.

▸ Backpressure
When overloaded, signal upstream to slow down (429 + Retry-After, bounded queues, load shedding) instead of collapsing.

▸ Interview probes
Token vs leaky bucket; distributed rate limiting with Redis; what backpressure is and why bounded queues matter.
