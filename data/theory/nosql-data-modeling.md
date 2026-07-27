📖 NoSQL data modeling (system design)

Relational modeling is entity-first then join. NoSQL is ACCESS-PATTERN-first: design the data around the exact queries you'll run, because there are no joins.

▸ Principles
- List your access patterns FIRST; model to serve them in one read.
- Denormalize and DUPLICATE data (storage is cheap; joins/scatter-gather are not). Precompute what you'll read together.
- DynamoDB single-table design: one table, partition key (PK) + sort key (SK), overload them for multiple entity types; Global Secondary Indexes (GSIs) for alternate access patterns.
- MongoDB: EMBED (nested doc) when read together and bounded; REFERENCE (separate doc) when large/unbounded or shared.
- Choose a partition key that spreads load and matches queries → avoid hot partitions.

▸ Pitfalls
Unknown future queries (NoSQL punishes them); duplicated data drifting out of sync (update all copies / use streams); huge items; unbounded arrays; hot partitions.

▸ Interview probes
How you model in DynamoDB / Mongo; single-table design; embed vs reference; why denormalize; picking a partition key.

🔗 Further reading
• AWS — DynamoDB single-table design & best practices: https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/best-practices.html
• ByteByteGo — data modeling (YouTube): https://www.youtube.com/@ByteByteGo
• /learn partitioning · /learn database-selection
