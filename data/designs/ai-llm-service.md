🏗 Design an LLM-backed service (RAG / chat assistant) [2026]

Still a system-design question — RESHADED applies. The twist: the model is your most expensive, slowest, least predictable component. Design around it.

▸ Requirements
Answer questions over the company's docs (RAG) or serve an LLM feature. Non-functional: latency (sub-second to a few seconds), cost per request (tokens/GPU), knowledge freshness, quality & safety.

▸ RAG pipeline
Offline: chunk documents → embed → store vectors in a vector DB (pgvector/Pinecone/Milvus). Online: embed the query → top-K similarity search → build a prompt with the retrieved context → call the LLM → return with citations.

▸ Deep dives
- Retrieval = the storage-schema question in disguise: chunking strategy, embedding model, and index freshness (re-embed when docs change — replication/invalidation).
- Caching, new dimensions: exact-match cache barely helps; cache embeddings, cache retrieved context, and semantic caching (serve a stored answer when a new query is close enough) — with a staleness policy.
- Cost & capacity: a model call is orders of magnitude pricier than a DB read; GPU is the bottleneck → batch requests, route easy queries to a smaller model, cap output tokens.
- Failure: models time out, degrade, hallucinate → timeouts + fallback to a smaller model, graceful degradation to a non-AI path; monitor answer quality + cost/request, not just latency.

▸ Trade-offs & bottlenecks
Latency vs retrieval depth (K); cost vs quality (model size); freshness vs re-index cost; hallucination vs guardrails (grounding + citations); privacy of prompts/embeddings.

🔗 Further reading
• pgvector — Postgres vector search: https://github.com/pgvector/pgvector
• Anthropic docs — building with Claude: https://docs.anthropic.com/
• ByteByteGo — RAG / LLM design (YouTube): https://www.youtube.com/@ByteByteGo
