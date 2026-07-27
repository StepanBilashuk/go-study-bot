📖 Replication (system design)

Keep copies of data on multiple nodes for availability and read scaling. Three shapes: single-leader, multi-leader, leaderless.

▸ Shapes
- Single-leader: writes go to the leader, replicated to followers (sync or async). Simple; failover risk.
- Multi-leader: writes accepted in multiple DCs → conflicts need resolution (LWW, CRDTs).
- Leaderless (Dynamo/Cassandra): quorum reads/writes, R+W>N for overlap.

▸ Replication lag anomalies (know 3 + fixes)
- Read-your-own-writes → read from leader / sticky routing after a write.
- Monotonic reads (going back in time) → pin a user to one replica.
- Consistent prefix (causal order broken) → causal tracking / same partition.

▸ Pitfalls
Sync (safe, slow) vs async (fast, can lose writes on failover); split-brain.

▸ Interview probes
Name 3 lag anomalies and their fixes; sync vs async trade-off; quorum math R+W>N.
