🏗 Design a chat system (WhatsApp / Messenger)

▸ Requirements
Functional: 1:1 and group messaging, delivery + read receipts, online presence, offline delivery, history.
Non-functional: low latency, ordered per-conversation, high availability, huge concurrent connections.

▸ Estimation
500M DAU, billions of messages/day; millions of persistent connections.

▸ High-level design
Clients hold WebSocket connections to a Connection/Gateway layer (sticky). Messages → Chat service → persist + route to recipient's connection (or store for offline). Presence service tracks online status. Kafka/queue for fan-out.

▸ Deep dives
- Connections at scale: a Connection service per node holds sockets; a routing layer (Redis) maps user→node so any sender reaches any recipient.
- Delivery: message id + per-conversation sequence for ordering; ack states sent→delivered→read.
- Offline: store undelivered messages; push on reconnect; mobile push notifications.
- Groups: fan-out to N members (small groups push; large groups pull).

▸ Trade-offs & bottlenecks
Millions of open sockets (memory, LB WebSocket support); ordering per conversation; presence cost at scale (heartbeat); E2E encryption; delivery guarantees.

🔗 Further reading
• ByteByteGo — WhatsApp/chat design (YouTube): https://www.youtube.com/@ByteByteGo
• Hussein Nasser — WebSockets: https://www.youtube.com/@hnasr
