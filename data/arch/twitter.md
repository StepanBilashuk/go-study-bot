🏛 Twitter/X — the timeline (fan-out)

▸ The gist
The home timeline is the hard part: hundreds of millions of users, each following many accounts, expecting a fast, fresh feed — while some accounts have tens of millions of followers.

▸ Patterns to learn
- Fan-out on write (push): precompute each follower's timeline in a cache on tweet — fast reads, expensive for celebrities.
- Fan-out on read (pull): merge followees' recent tweets at read time — cheap writes, slower reads.
- Hybrid: push for normal users, pull-and-merge for celebrity accounts. Timelines in Redis; tweets in a partitioned store.
- Ranking layer (ML) on top of chronological.

▸ Maps to
/design news-feed · /learn caching · /learn scaling-databases · /learn redis-patterns

🔗 Read the real thing
• ByteByteGo — design a news feed / Twitter: https://blog.bytebytego.com/
• HighScalability — Twitter timelines: https://highscalability.com/
