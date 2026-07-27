📖 Dynamic Programming (algorithms)

Overlapping subproblems + optimal substructure. Solve each subproblem once (memoize top-down, or tabulate bottom-up). The hard part is defining the state.

▸ Recognition signals
"Count the ways", "min/max cost", "can you reach", with a choice at each step and reuse of subresults.

▸ Method
1) Define state (what the subproblem answer depends on). 2) Write the transition (how states combine). 3) Base cases. 4) Optionally reduce space (keep only the last row/two).

▸ Pitfalls
- Wrong/insufficient state is the usual failure.
- Off-by-one in indices and base cases.

▸ Interview probes
Climbing Stairs, Coin Change, House Robber, Longest Common Subsequence, Edit Distance, 0/1 Knapsack, Word Break. Say your state definition before coding.

🔗 Further reading
• NeetCode — patterns roadmap: https://neetcode.io/roadmap
• NeetCode (YouTube, per-problem): https://www.youtube.com/@NeetCode
• Tech Interview Handbook — algo cheatsheet: https://www.techinterviewhandbook.org/algorithms/study-cheatsheet/
• Errichto — Dynamic Programming (YouTube): https://www.youtube.com/@Errichto
