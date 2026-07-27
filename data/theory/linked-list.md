📖 Linked List (algorithms)

Pointers node-to-node. The tricks: fast/slow pointers, in-place reversal by rewiring next, and a dummy head to kill edge cases.

▸ Recognition signals
Reverse in place, cycle detection, find middle / kth-from-end, merge sorted lists, O(1) removal given a node.

▸ Core techniques
- Fast/slow (Floyd): cycle detection, middle node.
- Reversal: prev/curr/next, flip one link at a time.
- Dummy head: uniform handling of head removal/insertion.

▸ Pitfalls
- Losing the next pointer before you rewire.
- Null checks on curr and curr.next.

▸ Interview probes
Reverse Linked List, Linked List Cycle, Merge Two Sorted Lists, Remove Nth From End, LRU Cache (doubly-linked list + hash map).
