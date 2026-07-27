📖 Sliding Window (algorithms)

A window over a contiguous run: expand the right edge, shrink the left to keep a constraint. Turns O(n·k) into O(n).

▸ Recognition signals
"Longest/shortest/max/min contiguous subarray or substring with condition". Contiguity is the giveaway.

▸ Fixed vs variable
- Fixed size k: slide by one, add the entering element, remove the leaving one.
- Variable: grow right until invalid, then shrink left until valid again; record the best.

▸ Pitfalls
- Track window state cheaply (a count map, a running sum) — don't recompute.
- Knowing WHEN to shrink is the whole trick; write the invariant first.

▸ Interview probes
Longest Substring Without Repeating Characters, Minimum Window Substring, Max Sum Subarray of size K, Longest Repeating Character Replacement.
