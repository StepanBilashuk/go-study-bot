🏗 Design a ride-hailing / dispatch system (Bolt / Uber)

▸ Requirements
Functional: rider requests a ride, match a nearby driver, live-track, price, complete + pay.
Non-functional: low match latency, high availability, correct pricing/payment, city-scale geo.

▸ Estimation
1M DAU, ~2 rides/user/week; peak driver-location updates every ~4s → the dominant write load (100k active drivers × 0.25 QPS = 25k loc-updates/s).

▸ High-level design
Clients → gateway → Rider, Driver, Location (ingest + geo index), Matching/Dispatch, Trip (state machine), Pricing (surge), Payment, Notification. Kafka for events.

▸ Deep dives
- Geo index: partition the map with geohash / Google S2 cells; keep drivers in an in-memory grid (Redis) keyed by cell for O(cell) nearest-driver queries.
- Matching: on request, find candidate drivers in nearby cells, rank by ETA, offer sequentially (avoid double-assign with a short lock/lease per driver).
- Surge pricing: demand/supply per cell over a time window (stream aggregation).
- Live tracking: WebSocket/SSE per trip; downsample driver pings.

▸ Trade-offs & bottlenecks
Location write volume (batch/aggregate), hot cells (airport), match latency vs match quality, exactly-once on trip payment (idempotency + ledger), driver double-booking (lease/lock).

🔗 Further reading
• ByteByteGo — Uber system design (YouTube): https://www.youtube.com/@ByteByteGo
• Google S2 geometry: http://s2geometry.io/
• Uber Eng blog — marketplace & geo: https://www.uber.com/en-US/blog/engineering/
