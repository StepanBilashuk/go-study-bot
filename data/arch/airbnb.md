🏛 Airbnb — from monolith to SOA + search

▸ The gist
Migrated a Rails monolith to a service-oriented architecture; search/availability and pricing are core, latency-sensitive problems over a geospatial inventory.

▸ Patterns to learn
- Strangler migration: peel services off the monolith incrementally behind a gateway.
- Search over geo + availability + ranking (Elasticsearch + a serving layer); precompute and cache.
- Data pipelines (their open-source Airflow) for offline computation feeding online services.
- Standardized service infra so many teams move independently.

▸ Maps to
/learn microservices-patterns · /learn elasticsearch · /learn search-indexing · /learn scaling-databases

🔗 Read the real thing
• Airbnb Engineering: https://medium.com/airbnb-engineering
• (SOA migration, search infra) https://highscalability.com/
