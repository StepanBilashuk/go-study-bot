📖 Trees (algorithms)

Recursive structure. Master DFS (pre/in/post order) and BFS (level order); a BST's ordering invariant gives O(h) search.

▸ Recognition signals
Hierarchy, path/subtree problems, level-by-level processing, ordered lookups (BST).

▸ How it works
- DFS: recursion (or explicit stack); pick the order by when you process the node.
- BFS: a queue, process level by level.
- BST: left<node<right — search/insert in O(h).

▸ Pitfalls
- BST validation needs min/max bounds, not just parent comparison.
- Recursion depth on skewed trees; base cases.

▸ Interview probes
Level Order Traversal, Validate BST, Lowest Common Ancestor, Diameter of Binary Tree, Serialize/Deserialize, Balanced check.

🔗 Further reading
• NeetCode — patterns roadmap: https://neetcode.io/roadmap
• NeetCode (YouTube, per-problem): https://www.youtube.com/@NeetCode
• Tech Interview Handbook — algo cheatsheet: https://www.techinterviewhandbook.org/algorithms/study-cheatsheet/
