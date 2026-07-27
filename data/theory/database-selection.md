📖 Choosing a database (system design)

"When Mongo vs Postgres vs ClickHouse?" Pick by DATA MODEL + ACCESS PATTERN + consistency + scale. There's no universal best — polyglot persistence is normal.

▸ The families
- Relational (PostgreSQL, MySQL): structured data, relationships, JOINs, ACID transactions, complex queries. The default for OLTP. Postgres = the feature-rich default (JSONB, extensions, GIS).
- Document (MongoDB): flexible/nested schema, denormalized aggregates read/written together, fast-evolving schema, hierarchical data. Weak multi-document transactions historically.
- Wide-column (Cassandra, ScyllaDB): massive WRITE throughput, high availability, known query patterns, no joins. Great for write-heavy / time-series-ish at scale.
- Key-value (Redis, DynamoDB): O(1) lookup by key — caching, sessions, simple high-scale lookups.
- Columnar / OLAP (ClickHouse, BigQuery, Redshift, Snowflake): scan-and-aggregate over billions of append-mostly rows. Dashboards/analytics, NOT your app's OLTP.
- Graph (Neo4j): relationship-heavy, multi-hop traversals (social, fraud rings, recommendations).
- Search (Elasticsearch): full-text + faceted search.
- Time-series (Prometheus, InfluxDB, TimescaleDB): metrics/events indexed by time.

▸ How to decide (say this)
Start relational (Postgres) unless a specific need forces otherwise. Add a specialized store ONLY for a proven access pattern (search → ES, analytics → ClickHouse, cache/leaderboard → Redis, metrics → TSDB). Keep one system of record; sync derived stores from it.

▸ Interview probes
When Mongo vs Postgres vs ClickHouse; document vs relational trade-offs; polyglot persistence; why not just one DB for everything.

🔗 Further reading
• ByteByteGo — types of databases / how to choose (YouTube): https://www.youtube.com/@ByteByteGo
• Designing Data-Intensive Applications (Kleppmann): https://dataintensive.net
• /learn storage-retrieval · /learn analytics-olap · /learn nosql-data-modeling
