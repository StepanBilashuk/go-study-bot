🏛 Uber — realtime dispatch & geo

▸ The gist
Match riders to nearby drivers in real time at city scale, ingesting constant driver-location pings. Moved from a monolith to domain-oriented microservices.

▸ Patterns to learn
- Geospatial indexing with Google S2 cells; drivers bucketed by cell for fast nearest-driver queries.
- High-write location ingestion (downsample/batch); Kafka as the event backbone.
- Matching/dispatch as a stream-processing problem; surge pricing from demand/supply per cell.
- Trip state machine + idempotent payments (double-entry ledger).

▸ Maps to
/design ride-sharing · /design proximity-service · /learn kafka-log-based · /learn ledger-accounting

🔗 Read the real thing
• Uber Engineering blog: https://www.uber.com/en-US/blog/engineering/
• (Marketplace, H3 geo, schemaless — indexed on) https://highscalability.com/
