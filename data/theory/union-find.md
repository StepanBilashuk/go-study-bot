📖 Union-Find / Disjoint Set (algorithms)

Track a partition of elements into groups, with near-O(1) union and find using union-by-rank/size + path compression.

▸ Recognition signals
Connectivity/grouping, "number of connected components", cycle detection in an UNDIRECTED graph, Kruskal's MST, merging accounts/emails.

▸ How it works
- find(x): follow parents to the root; compress the path on the way.
- union(a,b): attach the smaller tree under the larger (by rank/size).
Amortized ~O(α(n)) — effectively constant.

▸ Pitfalls
- Skipping path compression OR union-by-rank → degrades to O(n).
- Using it for directed-graph cycles (it's for undirected).

▸ Interview probes
Number of Connected Components, Redundant Connection (find the cycle edge), Accounts Merge, Number of Islands II, Graph Valid Tree.

🔗 Further reading
• NeetCode — Union-Find: https://neetcode.io/roadmap
• William Fiset — Union Find (YouTube): https://www.youtube.com/@WilliamFiset-videos
• Tech Interview Handbook: https://www.techinterviewhandbook.org/algorithms/graph/
