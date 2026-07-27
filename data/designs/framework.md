🏗 The design-interview framework (RESHADED)

Walk these 8 steps out loud on every "design X". Interviewers grade structure, trade-offs, and communication — not a perfect diagram.

▸ 1. Requirements
Functional (what it does) + non-functional (scale, latency, availability, consistency). Ask clarifying questions; state assumptions. Lock scope before drawing.

▸ 2. Estimation
DAU → QPS (÷ ~86,400, then ×3-5 for peak), storage/year, bandwidth. Decide read-heavy vs write-heavy — it drives the whole design (read-heavy → caches/CDN/replicas; write-heavy → partitioned logs, batching, idempotent writes).

▸ 3. Storage / data model
Key entities + the store per access pattern (SQL for relations/txns, KV/NoSQL for scale, blob for media). Pick partition/shard keys.

▸ 4. High-level design
Block diagram: clients → LB/gateway → services → data stores → async (queue/stream). Show the request flow.

▸ 5. API design
Key endpoints mapping to the functional reqs (verbs, pagination, idempotency keys).

▸ 6. Detailed design (deep dives)
Zoom into the 1-2 hard components: the geo index, the dedup path, the ranking, the ledger. This is where senior signal shows.

▸ 7. Evaluation
Find the bottleneck (from your estimates), then fix it: shard, cache, replicate, queue. State trade-offs explicitly (consistency vs availability, latency vs throughput).

▸ 8. Distinctive components
The thing unique to THIS system (video → adaptive streaming + CDN; chat → millions of sockets).

▸ Latency numbers (memorize the orders of magnitude)
Memory ~100ns · SSD read ~100µs · same-DC RTT ~0.5ms · disk seek ~10ms · cross-region RTT ~50-150ms. Anything on the hot path that hits disk or crosses a region is your bottleneck.

▸ Common pitfalls
Jumping to components before requirements; skipping estimation; no trade-offs; over/under-engineering; ignoring failure & monitoring; bad time management.

🔗 Further reading
• ByteByteGo — system design 101: https://blog.bytebytego.com
• System Design Primer (donnemartin): https://github.com/donnemartin/system-design-primer
• System Design Handbook — interview guide: https://www.systemdesignhandbook.com/guides/system-design-interview/
