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
