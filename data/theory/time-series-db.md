📖 Time-series databases (system design)

Purpose-built for timestamped, append-mostly data with range + aggregate queries and retention. Metrics, monitoring, IoT sensors, financial ticks.

▸ The players
- Prometheus: PULL-based metrics (scrapes targets), PromQL, built for monitoring; short-term local storage + remote write.
- InfluxDB: push-based, high write throughput.
- TimescaleDB: a PostgreSQL extension — SQL + hypertables (time-partitioned) + continuous aggregates. Nice when you already run Postgres.

▸ Why not just Postgres
TSDBs give time-partitioning, heavy compression, downsampling / continuous aggregates, and retention policies out of the box, and handle relentless append load a normal OLTP table struggles with.

▸ Pitfalls
- High CARDINALITY (too many unique label/tag combinations) blows up memory/index — the #1 TSDB failure.
- Retention & downsampling (raw for days, rollups for months).
- Push vs pull (Prometheus pulls; short-lived jobs need a pushgateway).

▸ Interview probes
Why a TSDB over a relational table; the cardinality problem; retention/downsampling; pull vs push metrics; relation to observability.

🔗 Further reading
• Prometheus — overview & data model: https://prometheus.io/docs/introduction/overview/
• TimescaleDB — time-series on Postgres: https://docs.timescale.com/
• /learn observability
