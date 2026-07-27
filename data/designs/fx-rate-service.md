🏗 Design a real-time FX rate service (Wise)

▸ Requirements
Functional: serve current exchange rates for currency pairs; lock a rate into a quote for a TTL; high read fan-out.
Non-functional: low latency, freshness, availability, and consistency (a locked quote must be honoured).

▸ Estimation
Hundreds of pairs, huge read QPS (every transfer/quote reads rates), moderate write rate (rate updates per second from providers).

▸ High-level design
Rate providers → Ingest/normalize → Rate store (in-memory + cache) → Rate API (read) + Quote service (locks a rate). Publish rate changes on Kafka for subscribers.

▸ Deep dives
- Read path: rates are tiny and hot → keep in memory / Redis, refreshed on provider push; CDN/edge for public rates.
- Quotes: when a user starts a transfer, lock rate R with an expiry (e.g., 30 min); store the locked rate with the transfer so settlement uses R, not the live rate.
- Staleness vs freshness: serve last-known on provider outage; mark stale; circuit-break bad providers.
- Consistency: the locked quote is the source of truth for that transfer (immutability).

▸ Trade-offs & bottlenecks
Freshness vs cache TTL; provider outages (fallback + stale flag); spread/markup logic; honouring expired quotes (reject vs re-quote).

🔗 Further reading
• Wise — backend system design interviews: https://wise.jobs/system-design-interviews
• ByteByteGo — caching & read-heavy systems (YouTube): https://www.youtube.com/@ByteByteGo
• Redis — caching patterns: https://redis.io/docs/latest/develop/
