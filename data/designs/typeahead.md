🏗 Design search autocomplete / typeahead (Google suggest)

▸ Requirements
Functional: as the user types a prefix, return top-K suggestions ranked by popularity.
Non-functional: very low latency (<100ms), read-heavy, fresh-ish popularity.

▸ Estimation
Billions of queries/day; every keystroke is a request → extreme read QPS. Must be cached/edge-served.

▸ High-level design
Query logs → offline aggregation (counts per term) → build a suggestion index (trie with top-K per node, or precomputed prefix→suggestions map) → serve from an in-memory/edge cache.

▸ Deep dives
- Index: a trie where each node stores the top-K completions of its prefix (precomputed) → O(prefix) lookup. Or a prefix→[suggestions] KV, sharded by prefix.
- Ranking: popularity (query counts), recency, personalization; rebuild periodically (hourly/daily) from logs.
- Latency: cache aggressively; serve from memory/CDN; debounce on the client.
- Freshness: trending terms need faster update (a fast path merging recent counts).

▸ Trade-offs & bottlenecks
Precompute (fast reads, stale) vs on-the-fly (fresh, slow); index size vs K; personalization cost; typo tolerance (edit distance / fuzzy).

🔗 Further reading
• ByteByteGo — autocomplete/typeahead (YouTube): https://www.youtube.com/@ByteByteGo
• System Design Primer: https://github.com/donnemartin/system-design-primer
