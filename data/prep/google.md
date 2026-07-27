🎯 Google — interview prep (SWE phone screen → onsite)

Phone screen: 45-60 min live coding in a shared Google Doc (NO autocomplete/compiler/run). One medium problem (or two shorter). Scored 1-4 on Communication · Problem Solving · Coding · Verification. Borderline → a 2nd screen, not a reject. Google writes its OWN questions and extends them mid-interview.

▸ How to pass (from candidate reports)
- Practice in a plain Google Doc the last 1-2 weeks — no IDE.
- State the approach + complexity BEFORE coding; narrate trade-offs continuously (that's the Communication score).
- Expect a follow-up that changes an assumption (input becomes a stream, memory-constrained, property no longer holds). Prep the follow-up, not just the first solution.
- Leave the last few minutes to verify: walk an example, edge cases, off-by-one.

▸ Reported coding questions
All pairs summing to target (→ follow-up: array too big for memory → external sort + two pointers) → /learn arrays-hashing · /learn two-pointers
Linked list cycle (→ find cycle start) · reverse a linked list → /learn linked-list
Three binary trees identical · longest equal-value path in a tree → /learn trees
Longest increasing subarray (→ if you may replace one element) · longest increasing path in a grid (memoised) → /learn dynamic-programming · /learn graphs
Sorted array → height-balanced BST · missing element via binary search (→ multiple missing) → /learn binary-search · /learn trees
Longest palindromic substring → /learn dynamic-programming

▸ Reported system design
Design search autocomplete → /design typeahead
Design user-session management for a web app → /learn security-auth · /learn caching
Payment gateway, isolate the financial data layer → /design payment-system · /learn ledger-accounting
Always: /design framework

▸ 2026 change
A pilot (junior/mid, some US teams) replaces one coding round with a Gemini-assisted CODE-COMPREHENSION round: read/debug/improve an existing codebase, use the model well, spot when it's wrong.

Sources: Prepfully Google SWE guide & question bank, Google "How we hire".
