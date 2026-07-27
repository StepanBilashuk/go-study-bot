📖 Load balancing (system design)

Spread traffic across servers for scale and availability. It's the single entry point in front of a fleet.

▸ L4 vs L7
- L4 (transport): routes by IP/port, doesn't read the payload — fast, protocol-agnostic.
- L7 (application): understands HTTP — path/host routing, TLS termination, cookies, header rewrites. Smarter, slightly slower.

▸ Algorithms
Round-robin, weighted round-robin, least-connections, least-response-time, IP/consistent hash (stickiness).

▸ Mechanics
Health checks eject dead nodes; sticky sessions for stateful backends; the LB itself must be redundant (pairs + failover, or DNS/anycast).

▸ Pitfalls
LB as a single point of failure; stickiness vs even spread; draining connections on deploy.

▸ Interview probes
L4 vs L7; algorithms and when to use each; health checks; how you make the load balancer itself highly available.

🔗 Further reading
• ByteByteGo — load balancing (YouTube): https://www.youtube.com/@ByteByteGo
• ByteByteGo newsletter: https://blog.bytebytego.com
• Cloudflare — what is load balancing: https://www.cloudflare.com/learning/performance/what-is-load-balancing/
