🏗 Design a payment system (Stripe / Enfuce)

▸ Requirements
Functional: charge a card, handle async provider callbacks, refund, ensure each charge happens once.
Non-functional: exactly-once money movement, auditability, PCI compliance, high availability.

▸ Estimation
Thousands of TPS at scale; correctness dominates.

▸ High-level design
Client → API (idempotency key required) → Payment orchestrator → PSP/card network (async), Ledger (double-entry), Webhook handler (provider callbacks), Notification. Outbox + Kafka for reliable events.

▸ Deep dives
- Idempotency: store request key → result; concurrent same-key requests serialize via a unique index; return the stored result on replay.
- Async settlement: authorize → capture; provider confirms via webhook (at-least-once → dedup on event id). State machine: pending → authorized → captured → settled / failed / refunded.
- Ledger: double-entry, immutable; a charge debits customer, credits merchant balance.
- Dual-write problem: write DB row + emit event atomically via the outbox pattern (CDC ships the event).
- Reconciliation with the PSP's reports catches drift.

▸ Trade-offs & bottlenecks
Exactly-once vs at-least-once reality (dedup); webhook ordering/retries; hot merchant accounts; secrets/PCI scope.

🔗 Further reading
• Stripe — Designing robust APIs with idempotency: https://stripe.com/blog/idempotency
• microservices.io — outbox & saga: https://microservices.io/patterns/data/transactional-outbox.html
• ByteByteGo — payment system design (YouTube): https://www.youtube.com/@ByteByteGo
