📖 Microservices patterns (system design)

Splitting into services buys independent deploys and scaling — but you lose distributed transactions and gain failure modes. Patterns tame both.

▸ Data consistency
- Saga: a sequence of local transactions with compensating actions on failure (choreography via events, or orchestration).
- Outbox: write the state change AND an event to the same DB transaction, then ship the event via CDC — solves the dual-write problem.
- Avoid 2PC across services (blocking, brittle).

▸ Resilience
Circuit breaker (fail fast when a dependency is down), retries with backoff + jitter, bulkhead (isolate pools), timeouts, service discovery, API gateway.

▸ Pitfalls
Distributed monolith (chatty, coupled); dual-write inconsistency; saga complexity; retry storms.

▸ Interview probes
Saga vs 2PC; outbox pattern and why; circuit breaker; service discovery; how to avoid dual-write inconsistency.

🔗 Further reading
• ByteByteGo — microservices patterns (YouTube): https://www.youtube.com/@ByteByteGo
• microservices.io — pattern catalog (Chris Richardson): https://microservices.io/patterns/
• ByteByteGo newsletter: https://blog.bytebytego.com
