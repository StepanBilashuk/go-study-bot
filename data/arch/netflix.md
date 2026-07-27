🏛 Netflix — microservices + Open Connect CDN

▸ The gist
Hundreds of microservices on AWS for control-plane (browse, recommendations, playback auth), and Netflix's OWN CDN (Open Connect) — appliances inside ISPs — for the actual video bytes.

▸ Patterns to learn
- Resilience engineering: circuit breakers (Hystrix), bulkheads, timeouts, and Chaos Monkey (deliberately kill instances to prove fault tolerance).
- Open Connect: push popular content to ISP-local caches ahead of demand → serves the vast majority of traffic off-origin.
- Async pipelines + big-data for recommendations; per-service data stores.

▸ Maps to
/learn microservices-patterns · /learn cdn-edge · /learn observability · /design video-streaming

🔗 Read the real thing
• Netflix Tech Blog: https://netflixtechblog.com/
• Open Connect (Netflix CDN): https://openconnect.netflix.com/
