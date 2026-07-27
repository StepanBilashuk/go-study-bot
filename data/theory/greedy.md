📖 Greedy (algorithms)

Make the locally optimal choice at each step and never reconsider. Fast and simple — but only correct if you can justify it (an exchange argument).

▸ Recognition signals
"Maximum/minimum …" with a natural ordering, scheduling, intervals, "reach the end". Often: sort, then sweep making the obvious choice.

▸ How it works vs DP
Greedy commits immediately; DP explores. Greedy works when a local optimum provably leads to a global one. If a counterexample exists, you need DP.

▸ Pitfalls
- The #1 trap: assuming greedy works without proof. Always sanity-check with a small case.
- Wrong sort key.

▸ Interview probes
Jump Game I/II, Gas Station, Task Scheduler, Partition Labels, Merge Intervals, Maximum Subarray (Kadane), Hand of Straights.

🔗 Further reading
• NeetCode — Greedy: https://neetcode.io/roadmap
• NeetCode (YouTube): https://www.youtube.com/@NeetCode
• Tech Interview Handbook: https://www.techinterviewhandbook.org/algorithms/greedy/
