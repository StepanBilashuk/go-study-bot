📖 Analytics & OLAP (system design)

Analytical queries have a different shape from transactional ones, and need a different engine.

▸ OLTP vs OLAP
- OLTP: many small reads/writes on a few rows, row-oriented (Postgres, MySQL). Runs the app.
- OLAP: scan-and-aggregate over billions of rows, column-oriented (BigQuery, Redshift, Snowflake, ClickHouse). Runs the dashboards.

▸ Why columnar
Reading only the needed columns + heavy compression makes aggregations orders of magnitude faster. Model as a star schema (fact table + dimension tables).

▸ Pipelines
ETL/ELT moves data from OLTP → warehouse; batch (Spark) vs stream (Flink/Kafka Streams) for freshness vs cost.

▸ Pitfalls
Running analytics on your OLTP DB (kills it); freshness vs cost; row vs columnar mismatch.

▸ Interview probes
OLTP vs OLAP; why columnar for analytics; star schema; batch vs stream processing; how data gets from app DB to warehouse.

🔗 Further reading
• ByteByteGo — OLTP vs OLAP / data pipelines (YouTube): https://www.youtube.com/@ByteByteGo
• Designing Data-Intensive Applications, ch.3 & 10 (Kleppmann): https://dataintensive.net
• ByteByteGo newsletter: https://blog.bytebytego.com
