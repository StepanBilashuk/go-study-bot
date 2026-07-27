📖 Graph databases (system design)

Store nodes + edges as first-class citizens and traverse relationships cheaply. When your queries are about CONNECTIONS, not tables.

▸ Why they win
Index-free adjacency: each node directly references its neighbors, so a hop is ~O(1). Multi-hop traversals ("friends of friends who like X", "is this account connected to a known fraud ring within 4 hops") are cheap — the same query in SQL means exploding self-joins.

▸ Where they fit
Social graphs, fraud-ring detection, recommendations, knowledge graphs, network/dependency/permission graphs. Neo4j (Cypher query language) is the common choice; also Amazon Neptune, JanusGraph.

▸ Where they don't
Aggregate analytics over all data (use OLAP), simple CRUD (use relational), massive horizontal sharding (graphs are hard to shard — traversals cross partitions).

▸ Interview probes
When a graph DB beats relational; index-free adjacency; a multi-hop query example; why sharding a graph is hard.

🔗 Further reading
• Neo4j — graph database concepts & Cypher: https://neo4j.com/docs/getting-started/
• ByteByteGo — graph databases (YouTube): https://www.youtube.com/@ByteByteGo
• /learn database-selection
