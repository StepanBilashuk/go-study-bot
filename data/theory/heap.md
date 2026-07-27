📖 Heap / Priority Queue (algorithms)

A partially-ordered tree giving O(log n) push/pop and O(1) peek of the min (or max). The go-to for "top-K", "merge-K", and streaming order.

▸ Recognition signals
Top/bottom K, k-th largest, merge k sorted streams, running median, schedule-by-priority.

▸ Patterns
- Top-K: keep a size-K heap (min-heap for K largest) → O(n log k).
- Two heaps: max-heap of the low half + min-heap of the high half → running median.
- Merge-K: heap of the current heads.

▸ Pitfalls
- Bound the heap size for top-K (don't heapify everything if K is small).
- Custom comparators / stability.

▸ Interview probes
Kth Largest Element, Top-K Frequent, Merge K Sorted Lists, Find Median from Data Stream, Task Scheduler.
