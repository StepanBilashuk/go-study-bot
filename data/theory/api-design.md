📖 API design: REST, gRPC, GraphQL (system design)

The contract between client and service. Pick the protocol for the use case.

▸ Styles
- REST: resource + HTTP verbs, cacheable, simple, ubiquitous. Chatty; over/under-fetching.
- gRPC: binary Protobuf over HTTP/2, typed, streaming, low latency — great service-to-service. Not browser-native.
- GraphQL: client selects exactly the fields it needs, one endpoint. Flexible; N+1 and caching are harder.

▸ Cross-cutting
Pagination (cursor beats offset at scale), versioning (URL vs header), idempotency keys, rate limits, auth, consistent error shapes.

▸ Pitfalls
Offset pagination on huge tables; breaking changes without versioning; GraphQL N+1; ignoring idempotency on writes.

▸ Interview probes
REST vs gRPC vs GraphQL trade-offs; cursor vs offset pagination; versioning strategy; designing a clean, evolvable API.

🔗 Further reading
• ByteByteGo — API design & REST vs gRPC vs GraphQL: https://www.youtube.com/@ByteByteGo
• ByteByteGo newsletter: https://blog.bytebytego.com
• Google API Design Guide: https://cloud.google.com/apis/design
