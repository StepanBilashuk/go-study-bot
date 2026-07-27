🏗 Design a news feed (Twitter / Facebook)

▸ Requirements
Functional: post; see a feed of people you follow, newest/ranked; like/comment.
Non-functional: low feed-read latency, high availability, handle celebrities (huge fan-out).

▸ Estimation
300M DAU, avg 200 follows; reads ≫ writes. Feed reads are the hot path.

▸ High-level design
Post service → fan-out. Two models:
- Fan-out on write (push): on post, write to each follower's feed cache. Fast reads, expensive for celebrities.
- Fan-out on read (pull): build the feed at read time by merging followees' recent posts. Cheap writes, slower reads.
Hybrid: push for normal users, pull for celebrities (merge their posts in at read time).

▸ Deep dives
- Feed store: per-user timeline in Redis (list of post ids); hydrate posts from a post store.
- Ranking: chronological → ML ranking (features: recency, affinity, engagement).
- Storage: posts in a partitioned store (by post id); social graph in a graph/KV store.

▸ Trade-offs & bottlenecks
Celebrity fan-out (hybrid), write amplification (push) vs read cost (pull), feed freshness vs cost, ranking latency.

🔗 Further reading
• ByteByteGo — news feed system (YouTube): https://www.youtube.com/@ByteByteGo
• System Design Primer: https://github.com/donnemartin/system-design-primer
