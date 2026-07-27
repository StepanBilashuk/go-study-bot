🏛 Stripe — payments API & reliability

▸ The gist
A developer-first payments API where correctness and reliability are the product. Idempotency, immutable ledgers, and rigorous API versioning are first-class.

▸ Patterns to learn
- Idempotency keys on every mutating request → safe retries, no double charges.
- Immutable, double-entry ledger as the source of truth for money movement + reconciliation with providers.
- API versioning that never breaks integrations (versioned request/response transforms).
- Rate limiting to protect the API; async webhooks (at-least-once → dedup).

▸ Maps to
/design payment-system · /learn idempotency · /learn ledger-accounting · /learn reconciliation · /learn api-design · /learn rate-limiting

🔗 Read the real thing
• Stripe — Designing robust APIs with idempotency: https://stripe.com/blog/idempotency
• Stripe — Scaling your API with rate limiters: https://stripe.com/blog/rate-limiters
