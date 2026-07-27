📖 Message queues (system design)

Async decoupling through a broker. Producers hand off work; consumers process later, absorbing spikes.

▸ Queue vs log
- Queue (RabbitMQ, SQS): a message is delivered to one consumer and removed — great for distributing work. Dead-letter queue (DLQ) captures failures.
- Log (Kafka): messages are retained and replayable, many consumers read independently.

▸ Delivery & ordering
At-least-once is the norm → make consumers idempotent. Ordering is per-queue/partition, not global.

▸ Pitfalls
Poison messages (need a DLQ + retry cap); assuming exactly-once (it's at-least-once + dedup); unbounded queues → backpressure; ordering assumptions.

▸ Interview probes
Queue vs log; at-least-once + idempotency; DLQ and retries; ordering guarantees; when to pick a queue over Kafka.

🔗 Further reading
• ByteByteGo — message queues (YouTube): https://www.youtube.com/@ByteByteGo
• AWS SQS — how it works: https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/welcome.html
• Hussein Nasser — queues & pub/sub: https://www.youtube.com/@hnasr
