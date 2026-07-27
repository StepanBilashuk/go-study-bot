📖 Event-driven architecture (system design)

Services communicate by emitting EVENTS (facts that happened) to a broker, instead of calling each other directly. Decoupled, scalable, auditable.

▸ Flavors
- Event notification: "order placed" — a thin ping; consumers fetch details.
- Event-carried state transfer: the event carries the data → consumers keep local copies, fewer callbacks.
- Event sourcing: the append-only LOG of events IS the source of truth. Current state = replay the events. Full audit + time-travel + rebuild projections.

▸ CQRS
Command Query Responsibility Segregation: separate the write model (commands mutate state / append events) from read models (projections optimized per query). Often paired with event sourcing.

▸ Coordination
Choreography (each service reacts to events, no central brain — loose, but flow is implicit) vs orchestration (a coordinator/saga drives the steps — explicit, easier to reason about).

▸ Pitfalls
Eventual consistency everywhere; event schema evolution (versioning); replay cost & snapshots; debugging a flow spread across services; the dual-write problem → outbox pattern.

▸ Interview probes
Event sourcing vs CRUD; CQRS and why; choreography vs orchestration; when EDA adds too much complexity; keeping the write DB and the event stream consistent (outbox).

🔗 Further reading
• Martin Fowler — Event Sourcing & CQRS: https://martinfowler.com/eaaDev/EventSourcing.html
• microservices.io — event-driven patterns: https://microservices.io/patterns/data/event-driven-architecture.html
• /learn kafka-log-based · /learn microservices-patterns
