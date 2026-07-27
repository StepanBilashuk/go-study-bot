🏛 Discord — storing billions of messages

▸ The gist
Real-time chat + voice at massive scale. Backend services in Elixir (BEAM, like WhatsApp) for concurrency; messages moved from MongoDB → Cassandra → ScyllaDB as volume exploded into the trillions.

▸ Patterns to learn
- Pick a store for the write pattern: wide-column (Cassandra/ScyllaDB) partitioned by channel + time bucket handles relentless append + range reads.
- BEAM/Elixir for millions of concurrent stateful connections; a routing layer maps user → node.
- Careful partition key design to avoid hot channels; read repair, tombstones, compaction realities.

▸ Maps to
/learn database-selection · /learn nosql-data-modeling · /learn concurrency · /design chat

🔗 Read the real thing
• Discord — How Discord stores trillions of messages: https://discord.com/blog/how-discord-stores-trillions-of-messages
• Discord Engineering blog: https://discord.com/category/engineering
