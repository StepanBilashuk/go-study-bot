📖 Proxies & API gateway (system design)

Intermediaries that sit between clients and servers.

▸ Forward vs reverse proxy
- Forward proxy: sits in front of clients (egress control, anonymity, corporate filtering).
- Reverse proxy: sits in front of servers (ingress) — TLS termination, caching, compression, load balancing, request routing (nginx, Envoy, HAProxy).

▸ API gateway
A reverse proxy specialized for microservices: auth, rate limiting, routing, request aggregation/composition, protocol translation, observability — so services don't each reimplement cross-cutting concerns. A service mesh (Envoy sidecars) moves service-to-service concerns into the infra layer.

▸ Pitfalls
Gateway as SPOF/bottleneck (scale it, keep it thin); too much business logic in the gateway; added latency hop.

▸ Interview probes
Forward vs reverse proxy; what an API gateway does; gateway vs service mesh; where to put auth and rate limiting.

🔗 Further reading
• ByteByteGo — API gateway & proxies (YouTube): https://www.youtube.com/@ByteByteGo
• NGINX — reverse proxy guide: https://docs.nginx.com/nginx/admin-guide/web-server/reverse-proxy/
• ByteByteGo newsletter: https://blog.bytebytego.com
