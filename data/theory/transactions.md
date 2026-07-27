📖 Transactions & isolation (system design)

ACID. Isolation levels trade correctness for performance by preventing specific anomalies.

▸ Levels and what each prevents
- Read Committed: no dirty reads.
- Snapshot / Repeatable Read (MVCC): no non-repeatable reads; reads see a consistent snapshot.
- Serializable: as if transactions ran one at a time — prevents everything, incl. write skew and phantoms.

▸ Anomalies to name
Dirty read, non-repeatable read, phantom, lost update, write skew (two txns each read then write, both valid alone, together break an invariant — needs serializable).

▸ Distributed
Two-phase commit (2PC) for atomic commit across nodes; blocking on coordinator failure.

▸ Pitfalls
Picking too weak a level; assuming snapshot isolation stops write skew (it doesn't); deadlocks.

▸ Interview probes
Name each isolation level and the anomaly it prevents; explain write skew; how MVCC works.

🔗 Further reading
• Designing Data-Intensive Applications, ch.7 (Kleppmann): https://dataintensive.net
• Hermitage — isolation levels tested across DBs (Kleppmann): https://github.com/ept/hermitage
• Jepsen — consistency & isolation models: https://jepsen.io/consistency
