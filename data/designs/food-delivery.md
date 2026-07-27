🏗 Design a food-delivery platform (Wolt / Bolt Food)

▸ Requirements
Functional: browse restaurants near me, place an order, match a courier, track the courier live, pay.
Non-functional: low latency for browse & tracking, high availability, correctness for orders/payments, one-city→many-cities scale.

▸ Estimation
1M DAU · ~0.3 orders/user/day → ~300k orders/day, peak ~20–50 orders/s (lunch/dinner spikes 3–5×). Live-location updates from couriers every ~4s dominate write QPS.

▸ High-level design
Clients → API gateway → services: Restaurant/Catalog (read-heavy, cached + CDN for images), Order (state machine), Dispatch/Matching (assign courier), Location (ingest courier pings), Payment (idempotent), Notification. Kafka between them for events.

▸ Order lifecycle
State machine: created → paid → accepted → cooking → picked_up → delivered. Persist transitions; emit events. Use idempotency keys on create/pay.

▸ Deep dives
- Geo search for nearby restaurants/couriers: geohash / S2 cells or a geo index (PostGIS, Redis GEO).
- Dispatch: match order↔courier by ETA/distance; a matching service consuming location + order streams.
- Live tracking: couriers push location to a Location service → fan out to the customer via WebSocket/SSE (per-order channel).
- Payments: see /design payment-system.

▸ Trade-offs & bottlenecks
Location write volume (batch, downsample), hot restaurants (cache), matching latency vs quality, exactly-once on money (idempotency + ledger).

🔗 Further reading
• ByteByteGo — Uber/DoorDash-style design (YouTube): https://www.youtube.com/@ByteByteGo
• Redis — geospatial (GEO) commands: https://redis.io/docs/latest/develop/data-types/geospatial/
• Uber Eng — dispatch & geo: https://www.uber.com/en-US/blog/engineering/
