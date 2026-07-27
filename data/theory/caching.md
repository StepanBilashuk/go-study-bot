📖 Caching & invalidation (system design)

Keep hot/expensive data in a faster tier. The hard part isn't caching — it's invalidation.

▸ Patterns
- Cache-aside (lazy): app checks cache, on miss reads DB and populates. Most common.
- Read-through: cache library loads on miss.
- Write-through: write cache + DB synchronously (consistent, slower writes).
- Write-behind: write cache, flush to DB async (fast, can lose data on crash).

▸ Failure modes
- Stale data → TTL, explicit invalidation, versioned keys.
- Thundering herd / stampede on expiry → locks, request coalescing, jittered TTL.
- Cache penetration (misses for nonexistent keys) → cache negatives / bloom filter.

▸ Interview probes
Write-through vs write-behind vs cache-aside and their failure modes; how you invalidate; stampede protection.
