📖 Reconciliation (ledger vs external) (system design)

Prove that two independently-kept records agree — your internal ledger vs the bank/PSP statement. Essential in fintech (Wise, Stripe, Enfuce) because distributed systems and external rails DRIFT: dropped messages, partial failures, timing, provider quirks. Idempotency prevents duplicates; reconciliation catches what still slips through.

▸ How it works
Pull the external record set (bank/PSP file, API) → match against your internal entries by a shared key (transaction id / provider reference / idempotency key) → classify each:
- Matched (amount + status agree) ✓
- Missing internally (external has it, you don't) — you dropped an event.
- Missing externally (you have it, provider doesn't) — it never landed.
- Amount/status mismatch — fees, FX, partial capture.
Auto-repair the safe cases; queue the rest as exceptions for a human/case-management flow.

▸ Design points
Append-only immutable ledger + idempotent postings make re-runs safe. Run scheduled (e.g. daily) plus near-real-time for critical flows. Levels: transaction, balance (sum check), position.

▸ Pitfalls
In-flight timing windows (settled after the cutoff); currency rounding / minor units; duplicates without a shared key; treating idempotency as enough (it isn't — you still reconcile).

▸ Interview probes
Why reconcile if you already have idempotency; how you match; how you classify and repair mismatches; achieving exactly-once money movement via reconciliation.

🔗 Further reading
• Modern Treasury — reconciliation & ledgers: https://www.moderntreasury.com/learn/what-is-reconciliation
• Wise — backend system design interviews: https://wise.jobs/system-design-interviews
• /learn ledger-accounting · /learn idempotency
