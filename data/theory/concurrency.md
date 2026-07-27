📖 Concurrency & threading (system design)

Single-process concurrency — Bolt asks deadlocks, ACID isolation, pass-by-value. Distinct from distributed coordination (that's across machines).

▸ The core problem
Threads share memory. Unsynchronized access to shared mutable state = a race condition (result depends on timing). Fix by making critical sections mutually exclusive.

▸ Primitives
- Mutex/lock: one holder at a time. Semaphore: N permits (bounded pool). Condition variable: wait/signal.
- Atomics / CAS (compare-and-swap): lock-free updates for a single variable.
- Memory model / visibility: without a happens-before (volatile/atomic/lock), one thread may not SEE another's write. Reordering is real.

▸ Deadlock (the 4 Coffman conditions)
Mutual exclusion + hold-and-wait + no preemption + circular wait — all four must hold. Break any one:
- Prevention: lock ordering (always acquire in the same order) kills circular wait; take all locks at once kills hold-and-wait.
- Avoidance: Banker's algorithm (grant only if a safe state remains).
- Detection: build a wait-for graph, find cycles, abort a victim.
Related: livelock (threads keep retrying, no progress), starvation (unfair scheduling).

▸ Higher-level
Thread pools (bound concurrency), producer-consumer (bounded queue), async/event loop (one thread, non-blocking I/O — Node/Netty). Pass-by-value copies; pass-by-reference shares (aliasing → concurrency hazards).

▸ Interview probes
Deadlock's 4 conditions + prevention/avoidance/detection; race condition + how you'd fix it; mutex vs semaphore; what CAS/atomics buy you; why a lock is also a memory barrier.

🔗 Further reading
• Jenkov — Java concurrency & threading: https://jenkov.com/tutorials/java-concurrency/index.html
• The Little Book of Semaphores (Downey): https://greenteapress.com/wp/semaphores/
• Deadlock — Coffman conditions (Wikipedia): https://en.wikipedia.org/wiki/Deadlock
