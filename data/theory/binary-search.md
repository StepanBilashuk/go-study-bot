📖 Binary Search (algorithms)

Halve the search space each step against a monotonic predicate. Not just for sorted arrays — also "binary search on the answer".

▸ Recognition signals
Sorted input, "find the boundary/first/last", or minimize/maximize a value where feasibility is monotonic (if X works, X+1 works).

▸ Templates
- Classic: lo/hi, mid=lo+(hi-lo)/2, shrink toward the target.
- On answer: define feasible(x); binary-search the smallest/largest feasible x.

▸ Pitfalls
- Boundary handling and infinite loops (mid rounding, lo=mid vs lo=mid+1).
- Overflow with (lo+hi)/2 — use lo+(hi-lo)/2.
- Getting the predicate direction wrong.

▸ Interview probes
Search in Rotated Sorted Array, Find First/Last Position, Koko Eating Bananas, Median of Two Sorted Arrays.
