📖 Bit Manipulation (algorithms)

Operate on the binary representation directly — often O(1) space and very fast.

▸ Core tricks
- XOR: a^a=0, a^0=a → find the single unpaired number.
- n & (n-1): clears the lowest set bit (count bits, power-of-two check).
- Masks: represent a subset of ≤ ~20 items as an integer (bitmask DP).
- Shifts: multiply/divide by powers of two.

▸ Recognition signals
"Without extra memory", powers of two, subsets, parity, adding without +.

▸ Pitfalls
- Signed vs unsigned shifts and overflow.
- Operator precedence (& below ==) — parenthesize.

▸ Interview probes
Single Number, Number of 1 Bits, Counting Bits, Reverse Bits, Sum of Two Integers (no +), Missing Number, Subsets via bitmask.

🔗 Further reading
• NeetCode — Bit Manipulation: https://neetcode.io/roadmap
• Tech Interview Handbook — bits: https://www.techinterviewhandbook.org/algorithms/binary/
• cp-algorithms — bit tricks: https://cp-algorithms.com/algebra/bit-manipulation.html
