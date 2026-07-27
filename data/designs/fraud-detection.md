🏗 Design a fraud-detection pipeline (Wise / fintech)

▸ Requirements
Functional: score transactions for fraud in near-real-time; block/hold risky ones; support rules + ML; feed back labels.
Non-functional: low added latency on the payment path, high recall on fraud, auditability, evolvable rules.

▸ Estimation
Same order as transfers (10s–100s/s), but each needs enrichment + scoring within ~100ms for inline decisions.

▸ High-level design
Transaction event → Kafka → Stream processor (enrich with features: velocity, device, geo, history) → Rules engine + ML model (score) → Decision (allow / review / block) → back to payment + a case-management queue for humans. Offline: feature store + model training on historical labels.

▸ Deep dives
- Inline vs async: fast rules inline (block obvious fraud) + async deep scoring for holds/review.
- Feature store: precompute aggregates (spend last 24h, count of new payees) — streaming aggregation (Flink/Kafka Streams) writing to a low-latency store.
- Idempotency & ordering: process each transaction once; per-account ordering (partition by account).
- Feedback loop: chargebacks/labels retrain the model; rules are versioned and hot-reloadable.

▸ Trade-offs & bottlenecks
Latency budget vs model depth; false positives (block good users) vs recall; feature freshness; explainability for compliance.

🔗 Further reading
• Wise — backend system design interviews: https://wise.jobs/system-design-interviews
• Confluent — streaming & Kafka Streams: https://www.confluent.io/blog/
• ByteByteGo — real-time fraud / streaming (YouTube): https://www.youtube.com/@ByteByteGo
