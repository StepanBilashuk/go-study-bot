📖 Comms styles: REST vs gRPC vs RabbitMQ vs WebSocket (system design)

A very common "difference between X and Y" question. They live on two axes: sync vs async, and the interaction shape.

▸ The four
- REST: SYNC request-response over HTTP/1.1, stateless, human-readable, cacheable, ubiquitous. Chatty; over/under-fetching. Public APIs, CRUD.
- gRPC: SYNC (or streaming) RPC over HTTP/2, binary Protobuf, typed contracts, low latency, bidirectional streaming. Not browser-native. Internal service-to-service, high performance.
- RabbitMQ (message broker): ASYNC. Producer fires a message and moves on; consumers process later. Work queues + pub/sub, at-least-once, DLQ. Decoupling, spikes, background work, events. (Kafka = the log variant: retained + replayable.)
- WebSocket: persistent, bidirectional, full-duplex over one TCP connection. Real-time push both ways. Chat, live feeds, collaboration.

▸ The axes to say out loud
- Sync (caller waits: REST, gRPC) vs async (fire-and-forget: RabbitMQ/Kafka).
- Request-response (REST/gRPC) vs pub/sub (RabbitMQ/Kafka) vs streaming/bidirectional (gRPC streams, WebSocket).
- Coupling: sync couples caller to callee's availability; async decouples via the broker.
- Latency vs throughput vs decoupling.

▸ When to pick
Public/CRUD → REST. Internal high-perf/typed → gRPC. Decouple / absorb spikes / events / background jobs → RabbitMQ or Kafka. Live bidirectional → WebSocket (or SSE for one-way).

▸ Interview probes
Difference between REST/gRPC/RabbitMQ/WebSocket; sync vs async; when would you use a queue instead of a REST call; gRPC vs REST trade-offs.

🔗 Further reading
• ByteByteGo — REST vs RPC vs GraphQL vs WebSocket (YouTube): https://www.youtube.com/@ByteByteGo
• gRPC — official docs: https://grpc.io/docs/what-is-grpc/introduction/
• /learn api-design · /learn message-queues · /learn websockets-realtime
