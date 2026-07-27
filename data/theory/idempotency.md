📖 Idempotency & exactly-once at the API level (system design)

An idempotent operation applied twice has the same effect as once. Essential because networks retry — "exactly-once delivery" doesn't exist; you get at-least-once + dedup.

▸ How to implement
- Idempotency key: client sends a unique key; server stores key → result. On replay, return the stored result instead of re-executing.
- Natural idempotency: PUT/upsert, unique constraints (a payment row keyed by request id).

▸ Pitfalls
- Key scope and expiry (how long do you remember?).
- Concurrent requests with the same key → row lock / unique index to serialize.
- Side effects (emails, external charges) must be tied to the same dedup.

▸ Interview probes
Design an idempotent payment endpoint unprompted; where you store the key; exactly-once vs at-least-once + dedup; handling concurrent duplicates.
