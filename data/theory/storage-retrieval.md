📖 Storage & retrieval: LSM vs B-tree (system design)

Two engine families. B-trees (Postgres, MySQL/InnoDB) update pages in place → read-optimized, predictable. LSM-trees (Cassandra, RocksDB, LevelDB) buffer writes in a memtable, flush sorted SSTables, compact in the background → write-optimized.

▸ When LSM beats B-tree
High write throughput: writes are sequential (fast on disk/SSD), no in-place random writes. Better write and space efficiency; downside is read and compaction cost.

▸ Amplification (the vocabulary they want)
- Write amp: one logical write causes many physical writes (B-tree: WAL + page; LSM: compaction).
- Read amp: a read may check several SSTables (mitigated by bloom filters).
- Space amp: stale copies before compaction.

▸ Indexes
Secondary indexes speed reads at the cost of writes. Clustered vs non-clustered.

▸ Interview probes
Explain when LSM beats a B-tree and why; the three amplifications; how bloom filters help LSM reads.

🔗 Further reading
• Designing Data-Intensive Applications, ch.3 (Kleppmann): https://dataintensive.net
• Arpit Bhayani — storage engines / LSM & B-trees: https://arpitbhayani.me
• ByteByteGo — SSTable & LSM Tree (YouTube): https://www.youtube.com/@ByteByteGo
