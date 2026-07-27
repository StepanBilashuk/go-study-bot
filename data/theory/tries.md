📖 Tries (Prefix Trees) (algorithms)

A tree keyed by characters: each path from the root spells a prefix. Insert/search/prefix-check in O(L) where L is the word length, independent of how many words are stored.

▸ Recognition signals
Prefix queries, autocomplete, dictionary/word games, "words starting with…", maximum-XOR (binary trie over bits).

▸ How it works
Each node holds child pointers (array[26] or a map) and an end-of-word flag. Walk char by char, creating nodes as needed. A binary trie stores numbers bit by bit for XOR problems.

▸ Pitfalls
- Memory: array[26] is fast but heavy; a map is compact but slower.
- Forgetting the end-of-word marker → false positives.

▸ Interview probes
Implement Trie, Word Search II (trie + backtracking over the grid), Design Add/Search Words (wildcard), Replace Words, Maximum XOR of Two Numbers.

🔗 Further reading
• NeetCode — Tries: https://neetcode.io/roadmap
• NeetCode (YouTube): https://www.youtube.com/@NeetCode
• Tech Interview Handbook — trie: https://www.techinterviewhandbook.org/algorithms/trie/
