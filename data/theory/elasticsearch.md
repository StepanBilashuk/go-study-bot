📖 Elasticsearch & ELK architecture (system design)

Distributed search + analytics engine on Lucene. Great for full-text search and log analytics — NOT a system of record (no real ACID; can lose/lag writes).

▸ Architecture
An index is split into SHARDS (each a Lucene index), spread across nodes; each shard has REPLICAS for availability + read scale. Writes go to the primary, replicate to replicas; searches scatter-gather across shards. Near-real-time (a ~1s refresh makes new docs searchable).

▸ Common architectures
- Search read-model (CQRS): your DB (Postgres) is the source of truth; you sync/denormalize documents into ES for fast full-text + faceted search. On DB change → update ES (via CDC/outbox).
- ELK / observability: Elasticsearch + Logstash/Beats (ingest) + Kibana (viz) for centralized logs and metrics.
- Analytics/aggregations over semi-structured data.

▸ Pitfalls
Mapping & analyzer choices change results; reindexing to change mappings; cluster sizing & shard count (too many shards hurts); split-brain (dedicated master nodes); treating ES as primary storage; keeping ES in sync with the source of truth.

▸ Interview probes
When to use ES vs a DB; ELK for logs; ES as a CQRS read model kept in sync via CDC; shards vs replicas; why ES isn't your system of record.

🔗 Further reading
• Elastic — architecture & inverted index: https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html
• ByteByteGo — Elasticsearch / log system (YouTube): https://www.youtube.com/@ByteByteGo
• /learn search-indexing · /learn event-driven-architecture
