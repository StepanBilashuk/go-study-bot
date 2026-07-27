📖 Arrays & Hashing (algorithms)

Trade space for time: a hash map/set gives O(1) average lookup, turning nested-loop O(n²) scans into a single O(n) pass by remembering what you've seen.

▸ Recognition signals
"Have I seen X before?", counting frequencies, dedup, finding a complement (two-sum: seen[target-x]), grouping.

▸ How it works
One pass, store each element (or its key) in a map; check membership/complement as you go. Group anagrams → key by sorted chars or char-count. Top-K → count then bucket/heap.

▸ Pitfalls
- Adversarial inputs can force hash collisions → O(n) per op.
- Mutating a key after insertion corrupts the map.
- Off-by-one when the answer is the element vs its index.

▸ Interview probes
Two Sum, Group Anagrams, Top-K Frequent, Longest Consecutive Sequence, Valid Anagram. They watch whether you jump straight to O(n) with a map instead of brute force.

🔗 Further reading
• NeetCode — patterns roadmap: https://neetcode.io/roadmap
• NeetCode (YouTube, per-problem): https://www.youtube.com/@NeetCode
• Tech Interview Handbook — algo cheatsheet: https://www.techinterviewhandbook.org/algorithms/study-cheatsheet/
