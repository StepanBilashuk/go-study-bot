📖 Object-oriented / low-level design (LLD) (system design)

A different interview from distributed system design: design the CLASSES and their interactions for a bounded system — parking lot, elevator, vending machine, Minesweeper. Bolt asks these directly.

▸ The method
1. Clarify requirements + a few use cases.
2. Identify entities (nouns) → classes; give each a single responsibility (SRP).
3. Model relationships (composition/aggregation/inheritance) and interfaces.
4. Define the key operations (verbs) and their APIs.
5. Apply patterns where they fit, then show extensibility.

▸ Patterns worth knowing
Strategy (pluggable pricing/allocation), State (order/game state), Factory (create by type), Observer (events), Singleton (registry), Decorator. Program to interfaces, not implementations.

▸ Pitfalls
God classes; deep inheritance instead of composition; missing SOLID; not handling concurrency (two cars, one spot); over-patterning.

▸ Worked examples (classic prompts)
Parking Lot (see /design parking-lot), Elevator system, Vending machine, Minesweeper, Rate limiter, Deck of cards, Library.

🔗 Further reading
• Refactoring Guru — design patterns & SOLID: https://refactoring.guru/design-patterns
• Grokking the low-level design (concepts): https://github.com/tssovi/grokking-the-object-oriented-design-interview
• ByteByteGo — OOD (YouTube): https://www.youtube.com/@ByteByteGo
