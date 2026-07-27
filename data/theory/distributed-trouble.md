📖 The trouble with distributed systems (system design)

Networks drop/delay/partition, clocks drift, and failures are partial. You can't distinguish a slow node from a dead one — only time out and guess.

▸ Clocks
Wall-clock (time-of-day) can jump backwards (NTP) and differ between nodes → unsafe for ordering events or expiring locks. Use logical clocks (Lamport, vector clocks) for causal order.

▸ Fencing
A paused node can wake up and act on a stale lock → use monotonically increasing fencing tokens so the storage rejects stale writers.

▸ Failure model
Crash-stop vs Byzantine (lying) nodes; most systems assume crash faults.

▸ Pitfalls
Assuming synchronized clocks; unbounded network delay; using timestamps for locking/ordering.

▸ Interview probes
Why wall-clock ordering is unsafe; timeouts and their trade-offs; what fencing tokens solve.

🔗 Further reading
• Designing Data-Intensive Applications, ch.8 (Kleppmann): https://dataintensive.net
• Jeff Hodges — Notes on distributed systems for young bloods: https://www.somethingsimilar.com/2013/01/14/notes-on-distributed-systems-for-young-bloods/
• aphyr — The trouble with timestamps: https://aphyr.com/posts/299-the-trouble-with-timestamps
