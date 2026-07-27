🎯 Bolt — interview prep (real reported questions)

Process: recruiter → live coding → system design (whiteboard/Excalidraw) → technical discussion → behavioral (HM). Algo-heavy, ~5 stages, ~40 days for senior. Below: actual questions recent Bolt SWE candidates reported (Prepfully/InterviewPal), each mapped to study.

▸ Coding — LeetCode medium/hard, live editor
Trees: serialize/deserialize N-ary tree · LCA (N-ary/family tree) · BST iterator · check 3 trees identical → /learn trees
Arrays/hashing: top-K frequent · majority element · longest consecutive sequence · merge two/K sorted · median of two sorted arrays → /learn arrays-hashing · /learn heap · /learn binary-search
Two-pointer/window: trapping rain water · min window substring · longest substring w/ 2 distinct → /learn two-pointers · /learn sliding-window
Strings/stack: min parentheses to make valid · valid Sudoku · decode ways → /learn stack · /learn dynamic-programming
Graphs/grid: connected components · highest edge-score node · shortest distance from all buildings (BFS) → /learn graphs · /learn union-find
DP/grid: minimum path sum · decode ways · longest increasing path → /learn dynamic-programming
Linked list: rotate by k · merge K lists → /learn linked-list
Math: palindromes in base-10 & base-k · K closest points → /learn math-geometry · /learn heap
Grind NeetCode; talk out loud; state complexity BEFORE coding.

▸ System design (whiteboard)
Distributed message queue, strict ordering per partition → /learn message-queues · /learn kafka-log-based
Proximity service (nearby) → /design proximity-service
Video transcoding pipeline → /design video-streaming
Distributed cache → /design distributed-cache
Twitter timeline / Facebook newsfeed / viral spikes → /design news-feed
Facebook Messenger / live comments → /design chat · /learn websockets-realtime
Autocomplete / typeahead → /design typeahead
Payment gateway (PCI, isolate financial data) → /design payment-system · /learn security-auth
Content personalization (real-time) → /design news-feed (ranking)
Graph partitioning & data locality → /learn partitioning · /learn advanced-graphs
Job scheduling service / design an API → /design job-scheduler · /learn api-design
OOD: Parking Lot, Minesweeper → /learn object-oriented-design · /design parking-lot
Routing / dispatch optimization (TSP/VRP) → /learn traveling-salesman
Always: /design framework (RESHADED)

▸ Technical discussion (concepts, deep-dive)
REST vs GraphQL (when to discourage GraphQL) → /learn api-design
SOLID / interface segregation smells → /learn object-oriented-design
Deadlock prevention vs avoidance vs detection · pass-by-value vs reference · concurrency control → /learn concurrency
ACID / isolation / concurrency control → /learn transactions
Forward vs reverse proxy / ingress · API gateway vs load balancer → /learn proxies-gateway · /learn load-balancing
How DNS resolution works → /learn dns
Microservices vs monolith (latency, IPC cost) → /learn microservices-patterns
SQL vs NoSQL · DB indexes (read vs write) · eventual vs strong consistency → /learn storage-retrieval · /learn consistency-consensus
REST vs RabbitMQ vs gRPC vs WebSocket → /learn communication-styles
Scaling a DB under load / high-load DB → /learn scaling-databases
Web-app security / zero-day incident response → /learn security-auth · /learn observability

▸ Behavioral (ownership, communication, growth)
Proven wrong (and learned) · took a risk · missed a deadline · difficult coworker · mentorship · tech debt vs a feature · why Bolt / tenure. They value: communicate clearly, own mistakes, ask for help.
Prep: /story · /stories · /boss behavioral

Sources: Prepfully & InterviewPal Bolt SWE candidate reports, Exponent, Glassdoor, Taro.
