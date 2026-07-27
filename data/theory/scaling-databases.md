📖 Scaling databases under load (system design)

"Design a high-load DB" / "the DB is the bottleneck — what now?" Climb this ladder in order; each step is more work than the last.

▸ The ladder
1. Optimize first: add the right indexes, fix N+1 and slow queries, and use CONNECTION POOLING (never one DB connection per request — you'll exhaust the server).
2. Cache reads: put Redis/Memcached in front (cache-aside) to offload repeated reads.
3. Read replicas: route reads to replicas for read scaling — mind replication lag (stale reads; read-your-writes → read from primary after a write).
4. Vertical scale: a bigger box. Simplest, but has a ceiling and a cost cliff.
5. Shard / partition: split data across nodes for WRITE scaling. Pick a shard key that spreads load and matches your access pattern; avoid hot partitions (a celebrity/hot key). Consistent hashing to minimize reshuffle.
6. Denormalize / CQRS: precompute read models, materialized views; separate the write model from the read model.
7. Right store per workload: OLTP (row) vs OLAP (columnar) vs KV/NoSQL for scale; blob for media.

▸ Pitfalls
Connection exhaustion (pool!); hot shard; replica lag serving stale data; cross-shard queries and transactions; unbounded growth (archive / TTL cold data); premature sharding (exhaust 1-3 first).

▸ Interview probes
How to scale a write-heavy DB; read vs write scaling; when to add replicas vs shard; connection pooling; handling hot partitions; keeping consistency across shards.

🔗 Further reading
• ByteByteGo — scaling databases / sharding (YouTube): https://www.youtube.com/@ByteByteGo
• Designing Data-Intensive Applications (Kleppmann): https://dataintensive.net
• /learn partitioning · /learn replication · /learn caching
