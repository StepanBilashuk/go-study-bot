🏗 Design a proximity service (Yelp / nearby drivers) [Bolt asked this]

▸ Requirements
Functional: given a location + radius, return nearby entities (restaurants, drivers, friends), ranked by distance; entities update location.
Non-functional: low-latency reads over millions of coordinates, high availability, frequent location updates.

▸ Estimation
100M places or 1M moving drivers; read-heavy for search, write-heavy if tracking moving objects (pings every few seconds).

▸ Core problem: fast geo search
Naive "distance to every point" is O(n) — too slow. Index space:
- Geohash: encode lat/lng into a string prefix; nearby points share prefixes → range-scan a prefix + neighbors.
- Quadtree: recursively split the map into 4 quadrants until each cell has ≤ K points; search the cell + neighbors.
- Google S2 / H3 hex cells: production-grade cell ids.
Redis GEO (geohash-backed) or PostGIS handle this out of the box.

▸ High-level design
Client → LocationSearch service (geo index in Redis/PostGIS) + a store of entity details. For moving objects, a Location-ingest service updates the index; downsample pings.

▸ Deep dives & trade-offs
Cell size vs result count (dense cities vs rural); border cases (search neighbor cells); static places (precompute) vs moving drivers (hot writes); ranking by ETA not just distance.

🔗 Further reading
• Redis — geospatial commands: https://redis.io/docs/latest/develop/data-types/geospatial/
• Google S2 geometry: http://s2geometry.io/
• ByteByteGo — proximity/nearby service (YouTube): https://www.youtube.com/@ByteByteGo
