📖 Graphs (algorithms)

Nodes + edges. Four workhorses: BFS (shortest unweighted), DFS (connectivity/cycles), topological sort (order a DAG), Dijkstra (weighted shortest path).

▸ Recognition signals
A grid is a graph; dependencies/prerequisites (topo sort); connectivity/components; shortest path.

▸ How it works
- BFS: queue, layer by layer → shortest hops.
- DFS: recursion/stack → components, cycle detection (colors).
- Topo sort: Kahn's (in-degrees) or DFS post-order; detects cycles.
- Dijkstra: min-heap by distance (non-negative weights).

▸ Pitfalls
- Track visited (or you loop forever); directed vs undirected; disconnected components.
- Dijkstra fails with negative edges (use Bellman-Ford).

▸ Interview probes
Number of Islands, Course Schedule (topo), Clone Graph, Network Delay Time (Dijkstra), Word Ladder.

🔗 Further reading
• NeetCode — Graphs: https://neetcode.io/roadmap
• William Fiset — Graph Theory (7h course): https://www.youtube.com/playlist?list=PLDV1Zeh2NRsDGO4--qE8yH72HFL1Km93P
• Tech Interview Handbook — graphs: https://www.techinterviewhandbook.org/algorithms/graph/

📚 Practice banks
• Top-100 coding questions — graph section (BFS/DFS, islands, bipartite, Dijkstra, Tarjan): https://shirsh94.medium.com/top-100-interview-programming-questions-that-asks-many-times-5c5bf36449ab
• DopplerHQ / awesome-interview-questions: https://github.com/DopplerHQ/awesome-interview-questions
