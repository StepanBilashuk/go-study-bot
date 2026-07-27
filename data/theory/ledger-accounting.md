📖 Ledgers & double-entry (money movement) (system design)

Moving money correctly is the core fintech problem (Wise, Stripe, Enfuce). The answer is an append-only, double-entry ledger — not mutable balance columns.

▸ Double-entry
Every movement is at least two entries that sum to zero: debit one account, credit another. Balances are derived by summing entries, never updated in place. This makes the system auditable and self-checking.

▸ Correctness
- Immutability: entries are append-only; corrections are new reversing entries, never edits.
- Idempotent postings: a transfer carries an idempotency key so retries don't double-post.
- Atomicity: both legs commit in one transaction (or a saga with compensations across services).
- Reconciliation: periodically prove internal ledger == external reality (bank statements).

▸ Pitfalls
Storing a single mutable balance (lost updates, no audit); floating-point money (use integer minor units); missing idempotency → double charges; multi-currency needs per-currency accounts + FX entries.

▸ Interview probes
Why double-entry and immutability; idempotent money transfer; how you handle multi-currency; reconciliation; ensuring exactly-once money movement.

🔗 Further reading
• Wise — backend system design interviews: https://wise.jobs/system-design-interviews
• Modern Treasury — the ledgers guide: https://www.moderntreasury.com/learn/what-is-a-ledger
• Martin Kleppmann — DDIA (transactions, ch.7): https://dataintensive.net
