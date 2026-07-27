📖 CDN & edge caching (system design)

Cache content in points-of-presence near users to cut latency and offload the origin.

▸ How it works
- Pull CDN: on a miss, fetch from origin and cache (by cache key). Push CDN: you upload ahead of time.
- Control with cache-control/TTL headers; purge/invalidate on change; design cache keys carefully (include only meaningful query params).
- Edge compute (Cloudflare Workers, Lambda@Edge) runs logic at the edge.

▸ Recognition
Static assets, media, global audience, read-heavy content, origin offload.

▸ Pitfalls
Invalidation lag; cache-key mistakes (caching per-user or per-query junk); dynamic/personalized content; cache poisoning.

▸ Interview probes
How a CDN works; cache keys and invalidation; static vs dynamic content; when edge compute helps; TTL trade-offs.

🔗 Further reading
• ByteByteGo — CDN (YouTube): https://www.youtube.com/@ByteByteGo
• Cloudflare — what is a CDN: https://www.cloudflare.com/learning/cdn/what-is-a-cdn/
• ByteByteGo newsletter: https://blog.bytebytego.com
