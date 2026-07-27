📖 Consistency & consensus (system design)

Linearizability = the system behaves as if there's one up-to-date copy, and operations take effect in real-time order. Strong but expensive. Eventual consistency = replicas converge, but reads can be stale.

▸ CAP (state it precisely)
During a network PARTITION you must choose: stay Consistent (reject some requests) or stay Available (serve possibly-stale data). When there's no partition you can have both — it's not "pick 2 always".

▸ Consensus
Raft/Paxos let nodes agree on a value/log despite failures: leader election + replicated log + majority quorum. Used for leader election, atomic commit, config.

▸ Linearizability vs serializability
Linearizability = recency on a single object; serializability = transactions equivalent to some serial order. Different guarantees.

▸ Interview probes
State CAP precisely (not the folk version); linearizability vs serializability; how Raft elects a leader and commits.

🔗 Further reading
• Kleppmann — Please stop calling databases CP or AP: https://martin.kleppmann.com/2015/05/11/please-stop-calling-databases-cp-or-ap.html
• aphyr — Strong consistency models: https://aphyr.com/posts/313-strong-consistency-models
• Jepsen — consistency model map: https://jepsen.io/consistency
• Raft, visually explained: https://thesecretlivesofdata.com/raft/
