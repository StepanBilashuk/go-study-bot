🏗 Design a distributed job scheduler (Bolt asked this)

▸ Requirements
Functional: schedule a job at a time or after a delay; run it once; retry on failure; recurring (cron) jobs.
Non-functional: reliable (no lost jobs), no duplicate execution surprises, scale to millions of jobs, low fire latency.

▸ Estimation
10M scheduled jobs, ~10k due/s at peak.

▸ High-level design
API → Job store (jobs: id, next_run, payload, status) → Scheduler (finds due jobs) → Queue → Workers (execute, ack). 

▸ Deep dives
- Finding due jobs: time-indexed store — a delay queue (Redis sorted set scored by run-time, or SQS delay) or a partitioned "next_run <= now" scan. Push due jobs to a work queue.
- Exactly-one scheduler tick: leader election / lease (etcd/ZooKeeper) OR partition the schedule so each shard has one owner — otherwise a cron fires N times.
- At-least-once execution → jobs must be idempotent (dedup key per (job, scheduled_time)).
- Crash recovery: visibility timeout — if a worker doesn't ack in T, re-queue.
- Missed runs after downtime: on restart, scan overdue and fire (with a catch-up policy).

▸ Trade-offs & bottlenecks
Duplicate fires (leader/lease + idempotency), thundering herd at round minutes (jitter), long jobs blocking workers, clock skew, hot time buckets (partition).

🔗 Further reading
• ByteByteGo — distributed job scheduler (YouTube): https://www.youtube.com/@ByteByteGo
• AWS — SQS delay queues & visibility timeout: https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-delay-queues.html
• Kleppmann — distributed locking (fencing): https://martin.kleppmann.com/2016/02/08/how-to-do-distributed-locking.html
