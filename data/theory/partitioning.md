📖 Partitioning / sharding (system design)

Split data across nodes so it scales beyond one machine. Partition by key range or by hash.

▸ Range vs hash
- Range: supports range scans, but risks hot spots (e.g. by timestamp).
- Hash: spreads load evenly, but kills range queries.
Consistent hashing minimizes how much data moves when a node joins/leaves.

▸ Rebalancing
Fixed number of partitions (move whole partitions) is simplest; dynamic splitting (like HBase) adapts to size. Avoid hash-mod-N (reshuffles everything on resize).

▸ Hot keys
A single popular key overwhelms one partition → add a random suffix / split the key / cache it.

▸ Secondary indexes
Local (per-partition, scatter-gather reads) vs global (partitioned by term, cross-partition writes).

▸ Interview probes
Explain rebalancing and hot-key handling; consistent hashing; local vs global secondary indexes.
