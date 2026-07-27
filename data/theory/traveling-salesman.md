📖 Traveling Salesman & routing (TSP/VRP) (algorithms)

TSP: visit every city exactly once and return to start, minimizing total distance. NP-hard — no known polynomial exact solution. Directly relevant to Bolt/Wolt dispatch and delivery routing.

▸ Exact solutions
- Brute force: try all permutations → O(n!). Only tiny n.
- Held-Karp (bitmask DP): dp[mask][i] = min cost to visit the set `mask` ending at city i → O(2^n · n²), feasible up to ~n≤18-20.

▸ Heuristics / approximation (what real systems use)
- Nearest-neighbor: greedily go to the closest unvisited city. Fast, ~25% above optimal.
- 2-opt / 3-opt: local search — repeatedly uncross edges to improve a tour.
- Christofides: 1.5-approximation for metric TSP (triangle inequality).
- Lin-Kernighan: strong practical heuristic.

▸ Real-world: VRP
Vehicle Routing Problem generalizes TSP to multiple vehicles + capacity + time windows — exactly the courier/driver dispatch problem. Solved with heuristics/metaheuristics (OR-Tools) at scale.

▸ Recognition & interview probes
"Shortest tour", route/delivery optimization, dispatch. Explain why it's NP-hard; sketch Held-Karp; give a practical heuristic; connect to VRP for real dispatch.

🔗 Further reading
• cp-algorithms — no direct TSP page; DP over bitmasks: https://cp-algorithms.com/
• William Fiset — TSP (Held-Karp) video: https://www.youtube.com/playlist?list=PLDV1Zeh2NRsDGO4--qE8yH72HFL1Km93P
• Google OR-Tools — routing (VRP): https://developers.google.com/optimization/routing
