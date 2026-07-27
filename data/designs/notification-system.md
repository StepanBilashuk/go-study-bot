🏗 Design a notification system (push / email / SMS)

▸ Requirements
Functional: send notifications across channels (push, email, SMS, in-app); templates; user preferences; retries.
Non-functional: high throughput, at-least-once delivery, respect opt-outs, observable delivery rates.

▸ Estimation
Billions of notifications/day at scale; bursty (marketing sends).

▸ High-level design
Producers → Notification API → queue (Kafka) → per-channel workers (push/APNs-FCM, email/SES, SMS/Twilio) → provider. Preference service filters; template service renders; a tracking service records sent/delivered/failed.

▸ Deep dives
- Decoupling & buffering: a queue absorbs spikes; per-channel workers scale independently; DLQ for poison messages.
- Delivery guarantees: at-least-once + idempotency (dedup key per (user, event)) so retries don't double-send.
- Rate limits & backpressure: respect provider limits; prioritize transactional over marketing.
- Fan-out & preferences: check opt-out/quiet-hours; batch digests.
- Tracking: providers confirm via webhooks (at-least-once → dedup); expose delivery metrics.

▸ Trade-offs & bottlenecks
Exactly-once myth (dedup); provider outages (retry + failover); ordering; prioritization; unsubscribe compliance. (This maps to your notification-delivery 85%→99% story.)

🔗 Further reading
• ByteByteGo — notification system design (YouTube): https://www.youtube.com/@ByteByteGo
• System Design Primer: https://github.com/donnemartin/system-design-primer
