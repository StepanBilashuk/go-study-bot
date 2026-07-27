📖 Intervals (algorithms)

Problems over ranges [start, end]. The move is almost always: sort by start (sometimes end), then sweep.

▸ Recognition signals
Overlapping ranges, merge, meeting rooms, calendars, "minimum resources to cover".

▸ How it works
- Merge: sort by start; if the next start ≤ current end, extend; else emit and move on.
- Count concurrent (meeting rooms): min-heap of end times, or a sweep line of +1/-1 events.

▸ Pitfalls
- Sort key: start vs end depends on the question.
- Boundary: does touching (end == next start) count as overlap? Clarify.

▸ Interview probes
Merge Intervals, Insert Interval, Non-overlapping Intervals, Meeting Rooms I/II (heap), Minimum Number of Arrows to Burst Balloons.

🔗 Further reading
• NeetCode — Intervals: https://neetcode.io/roadmap
• NeetCode (YouTube): https://www.youtube.com/@NeetCode
• Tech Interview Handbook: https://www.techinterviewhandbook.org/algorithms/interval/
