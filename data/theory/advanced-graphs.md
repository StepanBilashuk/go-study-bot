📖 Advanced Graphs — shortest paths & MST (algorithms)

Weighted-graph workhorses beyond BFS/DFS.

▸ Shortest path
- Dijkstra: non-negative weights, min-heap by distance, O(E log V). The default.
- Bellman-Ford: handles negative edges, detects negative cycles, O(V·E).
- Floyd-Warshall: all-pairs shortest paths, O(V³), simple DP.

▸ Minimum spanning tree
- Prim: grow the tree with a min-heap of edges.
- Kruskal: sort edges, add if they don't form a cycle (union-find).

▸ Recognition & pitfalls
Weighted shortest path, cheapest cost, min network to connect everything. Dijkstra breaks with negative edges → use Bellman-Ford. Know which tool fits.

▸ Interview probes
Network Delay Time (Dijkstra), Cheapest Flights Within K Stops (Bellman-Ford), Min Cost to Connect All Points (MST), Swim in Rising Water, Reconstruct Itinerary.

🔗 Further reading
• NeetCode — Advanced Graphs: https://neetcode.io/roadmap
• William Fiset — Graph Theory (7h): https://www.youtube.com/playlist?list=PLDV1Zeh2NRsDGO4--qE8yH72HFL1Km93P
• cp-algorithms — shortest paths & MST: https://cp-algorithms.com/
