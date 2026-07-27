📖 Distributed coordination (system design)

Getting independent nodes to agree: locks, leader election, membership, config.

▸ Distributed locks (do it safely)
A lock in Redis/etcd isn't enough — a process can acquire a lock, pause (GC), and wake up thinking it still holds it. Use fencing tokens: a monotonically increasing number the storage checks, so a stale holder's writes are rejected. Redlock is controversial for exactly this reason.

▸ Leader election & coordination
Elect one leader (single writer, one cron runner) via consensus (Raft) or a lease in ZooKeeper/etcd. Leases have a TTL and must be renewed; losing the lease means stepping down.

▸ Pitfalls
Locks without fencing (the GC-pause bug); split-brain (two leaders); relying on wall-clock time; a coordination service becoming an availability dependency.

▸ Interview probes
Why distributed locks are dangerous (fencing tokens); leader election; what ZooKeeper/etcd provide; lease TTL and renewal.

🔗 Further reading
• Kleppmann — How to do distributed locking (fencing): https://martin.kleppmann.com/2016/02/08/how-to-do-distributed-locking.html
• ByteByteGo — ZooKeeper / coordination (YouTube): https://www.youtube.com/@ByteByteGo
• etcd docs — distributed coordination: https://etcd.io/docs/
