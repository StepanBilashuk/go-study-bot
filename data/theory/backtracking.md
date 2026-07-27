📖 Backtracking (algorithms)

DFS over a decision tree: choose → recurse → undo. Prune branches that can't lead to a solution.

▸ Recognition signals
"Generate all …": subsets, permutations, combinations; constraint satisfaction (N-Queens, Sudoku); path search in a grid.

▸ How it works
At each step make a choice, recurse, then undo it (restore state) before the next choice. Prune early when a partial solution already violates a constraint.

▸ Pitfalls
- Forgetting to undo mutated state (the #1 bug).
- Duplicates: sort inputs, skip equal siblings.
- Missing the pruning that makes it tractable.

▸ Interview probes
Subsets, Permutations, Combination Sum, N-Queens, Word Search, Palindrome Partitioning. Exponential by nature — pruning is what's judged.
