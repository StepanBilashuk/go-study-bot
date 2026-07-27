🏗 Design international money transfer (Wise)

▸ Requirements
Functional: send money A→B across currencies; quote an FX rate; move funds; track state; notify.
Non-functional: correctness above all (no lost/double money), auditability, regulatory compliance, hundreds of thousands of transfers/day.

▸ Estimation
300k transfers/day, peak ~10–30/s. Low QPS, extreme correctness bar.

▸ High-level design
Client → gateway → Transfer service (orchestrates a saga) → Quote/FX, Ledger (double-entry), Payin (collect from sender), Payout (bank rails/SWIFT/local), Compliance/Fraud, Notification. Kafka for events + outbox.

▸ Deep dives
- Ledger: append-only double-entry; every step posts balanced entries; balances are derived. See /learn ledger-accounting.
- Idempotency: each transfer + each external call carries an idempotency key → safe retries, no double payin/payout.
- Saga: payin → convert → payout, with compensations (refund) on failure; state machine persisted.
- Multi-currency: per-currency accounts; FX conversion is a ledger entry at the quoted rate (rate locked for a TTL).
- Exactly-once with external rails that are at-least-once → dedup on provider references.

▸ Trade-offs & bottlenecks
Correctness vs latency (async payout is fine); external rail failures (retries + reconciliation); fraud holds; compliance (KYC/AML) gating.

🔗 Further reading
• Wise — backend system design interviews: https://wise.jobs/system-design-interviews
• Modern Treasury — ledgers & payment ops: https://www.moderntreasury.com/learn/what-is-a-ledger
• Stripe — idempotency: https://stripe.com/blog/idempotency
