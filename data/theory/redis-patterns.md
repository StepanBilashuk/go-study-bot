📖 Redis beyond caching (system design)

Redis is an in-memory data-structure store — a Swiss-army tool, not just a cache.

▸ Data structures → uses
- String / hash: cache, objects, counters (INCR).
- Sorted set (ZSET): leaderboards, priority queues, rate-limit windows, "top N".
- Set: uniqueness, tags, membership.
- List: simple queues (LPUSH/BRPOP).
- Bitmap / HyperLogLog: presence flags / approximate unique counts at tiny memory.
- Stream: an append-only log with consumer groups (a lightweight Kafka).
- Geospatial: nearby lookups.

▸ Patterns
Distributed lock (SET NX PX — with fencing caveats, see /learn distributed-coordination), rate limiting (INCR+EXPIRE or token bucket in Lua), session store, pub/sub, job queues, caching (of course).

▸ Operations
Persistence: RDB snapshots + AOF log (durability vs speed trade-off). Replication + Redis Cluster (sharding via hash slots). Eviction policies (LRU/LFU) when memory is full. Single-threaded command execution → avoid O(n) commands on huge keys.

▸ Interview probes
Redis uses beyond caching; sorted set for a leaderboard; rate limiting in Redis; is Redis durable (RDB/AOF); why big O(n) commands are dangerous.

🔗 Further reading
• Redis — data types & use cases: https://redis.io/docs/latest/develop/data-types/
• ByteByteGo — Redis use cases (YouTube): https://www.youtube.com/@ByteByteGo
• /learn caching · /learn rate-limiting · /learn distributed-coordination
