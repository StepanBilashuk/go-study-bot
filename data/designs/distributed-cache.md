🏗 Design a distributed cache (Redis/Memcached-style) [Bolt asked this]

▸ Requirements
Functional: get/set/delete by key with a TTL; scale beyond one node; survive node failures.
Non-functional: very low latency, high hit rate, horizontal scale, tolerate node add/remove.

▸ Estimation
Millions of ops/s; data spread across a cluster; each node holds a shard in memory.

▸ High-level design
Client library → hashes the key to a node → GET/SET on that node's in-memory store. A cluster of cache nodes; optional replicas per shard.

▸ Deep dives
- Sharding: consistent hashing (a ring with virtual nodes) so adding/removing a node reshuffles only ~1/N of keys, not everything.
- Replication: a replica per shard for availability; async replication (may lose recent writes).
- Eviction: LRU/LFU when memory is full; per-key TTL.
- Usage pattern: cache-aside (app reads DB on miss and populates); write-through/behind for writes.
- Hot keys: replicate the hot key / client-side cache; stampede protection (locks, jittered TTL).

▸ Trade-offs & bottlenecks
Consistency (cache vs DB) — TTL/invalidation; hot shards; cold-start after a node loss; memory vs hit-rate.

🔗 Further reading
• ByteByteGo — distributed cache (YouTube): https://www.youtube.com/@ByteByteGo
• Redis — cluster & consistent hashing: https://redis.io/docs/latest/operate/oss_and_stack/reference/cluster-spec/
• /learn caching · /learn partitioning
