🏛 Spotify — backend infrastructure

▸ The gist
Spotify runs hundreds of small, autonomous microservices owned by independent "squads". Music streaming is read-heavy and latency-sensitive; playlists/social are write-paths.

▸ Patterns to learn
- Microservices with strong team ownership (Conway's law on purpose): each squad owns its service end-to-end.
- Polyglot persistence: Cassandra (high-write, availability) for user data, Postgres where relations matter, plus caches.
- Event pipeline: huge event volume (plays, skips) flows through Kafka → data platform for recommendations/analytics.
- Client-side caching + CDN for audio delivery; precompute recommendations offline.

▸ Maps to
/learn microservices-patterns · /learn database-selection · /learn kafka-log-based · /learn cdn-edge · /design food-delivery (SOA shape)

🔗 Read the real thing
• Spotify Engineering — Backend infrastructure (2013): https://engineering.atspotify.com/2013/03/backend-infrastructure-at-spotify
• Spotify Engineering blog: https://engineering.atspotify.com/
