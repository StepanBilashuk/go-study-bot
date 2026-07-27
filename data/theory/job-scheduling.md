📖 Job scheduling & distributed cron (system design)

Run tasks on a schedule or after a delay, reliably, across many machines. Bolt has literally asked "design a job scheduling service".

▸ Core design
A store of jobs (next-run time, payload, status) + workers that claim due jobs. Poll a time-indexed table or use a delay queue (SQS delay, Redis sorted set by timestamp). Push claimed jobs onto a queue for execution.

▸ Guarantees
- At-least-once execution → make jobs idempotent (dedup key per run).
- Exactly-one-runner for the scheduler tick → leader election (or partition the schedule) so a cron doesn't fire N times.
- Handle missed runs (downtime) and clock skew; visibility timeout so a crashed worker's job is retried.

▸ Scale
Partition jobs by id/time; separate scheduling from execution; backpressure when workers lag.

▸ Pitfalls
Duplicate fires (no leader/lease); thundering herd at round times (jitter); long jobs blocking; losing jobs on crash (persist before ack).

▸ Interview probes
At-least-once + idempotency; avoiding duplicate fires (leader election/lease); missed runs; scaling to millions of jobs; delay queue vs polling.

🔗 Further reading
• ByteByteGo — distributed job scheduler (YouTube): https://www.youtube.com/@ByteByteGo
• AWS — SQS delay queues & visibility timeout: https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-delay-queues.html
• ByteByteGo newsletter: https://blog.bytebytego.com
