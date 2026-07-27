📖 Two Pointers (algorithms)

Two indices walking a sequence together, so you avoid a nested loop. Works when the data has order you can exploit.

▸ Recognition signals
Sorted array, pair/triplet summing to a target, in-place partition/dedup, palindrome check, merging.

▸ Variants
- Opposite ends (converging): sorted pair-sum, container-with-most-water — move the pointer that can improve the answer.
- Same direction (fast/slow): remove duplicates in place, partition.

▸ Pitfalls
- Skip duplicates to avoid repeated triplets (3Sum: advance past equal values).
- Off-by-one on the loop boundary and pointer invariants.

▸ Interview probes
3Sum, Container With Most Water, Valid Palindrome, Remove Duplicates from Sorted Array, Trapping Rain Water. Say the invariant out loud ("left/right can only shrink the window").

🔗 Further reading
• NeetCode — patterns roadmap: https://neetcode.io/roadmap
• NeetCode (YouTube, per-problem): https://www.youtube.com/@NeetCode
• Tech Interview Handbook — algo cheatsheet: https://www.techinterviewhandbook.org/algorithms/study-cheatsheet/
