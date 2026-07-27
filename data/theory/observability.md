📖 Observability (system design)

Understand a running system from the outside — essential once it's distributed.

▸ Three pillars
- Metrics: cheap numeric aggregates over time (rates, latencies, gauges).
- Logs: detailed discrete events (structured > plain text).
- Traces: a request's path across services, with spans and timing.

▸ Frameworks
RED (Rate, Errors, Duration) for services; USE (Utilization, Saturation, Errors) for resources. Define SLIs → SLOs → error budgets to drive alerting.

▸ Pitfalls
Metric cardinality explosion; log volume and cost; trace sampling; alert fatigue (alert on symptoms/SLOs, not every spike).

▸ Interview probes
The three pillars and when to use each; RED vs USE; SLI/SLO/error budget; how tracing correlates a request across services (trace/span IDs).

🔗 Further reading
• ByteByteGo — observability & monitoring (YouTube): https://www.youtube.com/@ByteByteGo
• Google SRE Book — SLOs & monitoring: https://sre.google/sre-book/service-level-objectives/
• OpenTelemetry docs: https://opentelemetry.io/docs/
