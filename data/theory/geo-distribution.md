📖 Geo-distribution & multi-region (system design)

Serve users from multiple regions for low latency and regional failure tolerance. Everything is a latency-vs-consistency trade-off.

▸ Patterns
- Read replicas per region: local low-latency reads, writes go to a home region.
- Geo-partitioning: pin each user's data to their home region (also helps residency).
- Active-active (multi-leader): writes anywhere → conflict resolution; vs active-passive failover.

▸ Constraints
- Cross-region round-trips are 50-150ms → don't do synchronous cross-region writes on the hot path.
- Data residency (GDPR): some data must stay in-region.

▸ Pitfalls
Conflict resolution, split-brain on failover, replication lag across regions, chatty cross-region calls.

▸ Interview probes
Latency vs consistency trade-off per region; data residency; active-active vs active-passive; where you place the write leader.
