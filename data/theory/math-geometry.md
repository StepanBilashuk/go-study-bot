📖 Math & Geometry (algorithms)

Number theory + grid/matrix manipulation. Less about a single pattern, more about careful implementation and avoiding overflow.

▸ Toolbox
- GCD/LCM via Euclid; modular arithmetic (keep numbers bounded).
- Fast exponentiation pow(x,n) in O(log n).
- Prime sieve (Eratosthenes).
- Matrix ops: rotate in place (transpose + reverse), spiral traversal, set-zeroes with markers.

▸ Recognition signals
Matrix manipulation, digit problems, combinatorics, "compute X mod 1e9+7", overflow-prone math.

▸ Pitfalls
- Integer overflow → use 64-bit / mod as you go.
- Off-by-one in grid boundaries; edge cases (0, negatives, empty).

▸ Interview probes
Rotate Image, Spiral Matrix, Set Matrix Zeroes, Happy Number, Pow(x, n), Multiply Strings, GCD of Strings.

🔗 Further reading
• NeetCode — Math & Geometry: https://neetcode.io/roadmap
• Tech Interview Handbook: https://www.techinterviewhandbook.org/algorithms/math/
• cp-algorithms — number theory: https://cp-algorithms.com/
