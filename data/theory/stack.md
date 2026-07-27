📖 Stack (algorithms)

LIFO. Two big uses: matching nested structure, and the monotonic stack for range "next greater/smaller" queries in O(n).

▸ Recognition signals
Parentheses/validity, nested/recursive structure, "next greater element", histogram/skyline, undo.

▸ Monotonic stack
Keep elements in increasing (or decreasing) order; when the incoming element breaks the order, pop and resolve — each element is pushed/popped once → O(n).

▸ Pitfalls
- Popping an empty stack (guard it).
- Choosing the monotonic direction (increasing vs decreasing) for the question asked.

▸ Interview probes
Valid Parentheses, Daily Temperatures, Largest Rectangle in Histogram, Next Greater Element, Min Stack (track the running min alongside values).
