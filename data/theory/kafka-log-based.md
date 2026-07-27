📖 Stream processing & log-based brokers (system design)

Kafka is an append-only, partitioned, replayable log. Producers append; consumers read at their own offset. This decouples producers from consumers and lets you re-read history.

▸ Why it wins
Buffering/decoupling, ordered delivery per partition, replay for reprocessing, one source feeding many consumers.

▸ Key mechanics
- Partitions give parallelism; ordering is guaranteed only WITHIN a partition.
- Consumer groups split partitions across consumers; rebalancing on membership change.
- Delivery: at-least-once by default; exactly-once via idempotent producer + transactions (dedup by producer id + sequence).
- CDC: stream a database's change log as events (Debezium) to keep systems in sync.
- Log compaction: keep the latest value per key.

▸ Pitfalls
Global ordering (only per-partition), consumer lag, rebalancing storms, exactly-once cost.

▸ Interview probes
Explain exactly-once semantics and CDC unprompted; partitioning strategy; consumer groups; when to use a log vs a queue.
