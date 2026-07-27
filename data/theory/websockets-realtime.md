📖 Real-time delivery (system design)

Getting data to clients as it happens. Three transports, increasing power:

▸ Transports
- Long polling: client re-requests; simple, wasteful, latency-y.
- SSE (Server-Sent Events): one-way server→client stream over HTTP; auto-reconnect; great for feeds/notifications.
- WebSocket: full-duplex, low-latency; for chat, collaboration, gaming.

▸ Scaling fan-out
Connections are stateful → sticky routing or a shared pub/sub layer (Redis pub/sub, Kafka) so any node can push to any client. Track presence; handle reconnection and backpressure.

▸ Pitfalls
Millions of open connections (memory, LB support); horizontal scaling needs shared pub/sub; missed messages on reconnect (buffer/resume); load balancers must support WebSocket upgrade.

▸ Interview probes
Polling vs SSE vs WebSocket; fan-out to millions; presence; how you scale stateful connections horizontally.

🔗 Further reading
• ByteByteGo — WebSocket / real-time (YouTube): https://www.youtube.com/@ByteByteGo
• Hussein Nasser — WebSockets deep dive: https://www.youtube.com/@hnasr
• Ably — realtime concepts: https://ably.com/topic/websockets
