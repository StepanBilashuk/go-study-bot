🏗 Design a URL shortener (TinyURL / bit.ly)

▸ Requirements
Functional: shorten a long URL → short code; redirect short→long; optional custom alias, expiry, analytics.
Non-functional: very read-heavy (redirects), low latency, high availability, short codes never collide.

▸ Estimation
100M new URLs/month; read:write ~100:1 → redirects dominate. 5 years ≈ 6B URLs → 7-char base62 (62^7 ≈ 3.5T) is plenty.

▸ High-level design
Client → API → Write path (generate code, store code→URL) + Read path (code→URL, redirect 301/302). Cache hot codes; CDN at edge.

▸ Deep dives
- Code generation: (a) hash(URL) + collision check, or (b) a distributed counter → base62 encode (no collisions, unpredictable via a range-per-node or Snowflake-style id). Avoid coordination per request by handing each node a block of ids.
- Storage: key-value (code → URL), 6B rows → partition by code. Read-optimized.
- Redirect caching: hot codes in Redis/CDN; 301 (permanent, cached) vs 302 (allows analytics).
- Analytics: async — emit a click event to Kafka, aggregate offline.

▸ Trade-offs & bottlenecks
Hashing (collisions) vs counter (predictability); 301 vs 302 (cacheability vs tracking); hot links (cache/CDN); custom aliases (uniqueness).

🔗 Further reading
• ByteByteGo — URL shortener (YouTube): https://www.youtube.com/@ByteByteGo
• System Design Primer — designing a URL shortener: https://github.com/donnemartin/system-design-primer
