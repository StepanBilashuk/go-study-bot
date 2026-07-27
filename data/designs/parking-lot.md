🏗 Design a parking lot (low-level / OOD) [Bolt asked this]

A classic OBJECT-oriented design prompt: design the classes, not a distributed system. Show clean responsibilities + extensibility.

▸ Requirements
Multiple levels, spot types (compact/large/handicap/EV), vehicles of matching sizes, park/unpark, ticketing, pricing, find nearest free spot, track availability.

▸ Class model
- ParkingLot (aggregates Levels; entry/exit).
- Level (has ParkingSpots; tracks free counts).
- ParkingSpot (id, SpotType, occupied, fits(Vehicle)).
- Vehicle (abstract) → Car / Bike / Truck (each has a Size).
- Ticket (spot, vehicle, entryTime).
- ParkingStrategy (interface) → nearest / by-type — Strategy pattern so allocation is pluggable.
- PricingStrategy (interface) → hourly / flat — Strategy pattern.

▸ Key operations
park(vehicle) → find a fitting free spot via ParkingStrategy → mark occupied → issue Ticket.
unpark(ticket) → free the spot → PricingStrategy.price(ticket).

▸ Design signals they grade
SRP (each class one job); composition over inheritance; interfaces for the varying bits (allocation, pricing); concurrency (two cars racing for one spot → lock/atomic decrement); extensibility (new vehicle/spot type without touching core).

🔗 Further reading
• /learn object-oriented-design (SOLID + patterns)
• Refactoring Guru — patterns: https://refactoring.guru/design-patterns
• Grokking OOD interview: https://github.com/tssovi/grokking-the-object-oriented-design-interview
