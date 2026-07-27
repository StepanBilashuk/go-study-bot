📖 AI / LLM system design (system design)

The 2026 differentiator round: "design a chat assistant / semantic search / the serving layer for an LLM feature". It's still system design — the model is just a new kind of component: the most expensive, slowest, and least predictable one you'll design around.

▸ RAG (retrieval-augmented generation)
Offline: chunk docs → embed → store vectors in a vector DB (pgvector, Pinecone, Milvus). Online: embed the query → top-K similarity search → stuff retrieved context into a prompt → call the LLM → answer with citations.

▸ What changes vs classic SD
- Caching gets a new axis: exact-match caching barely helps (everyone phrases differently) → cache embeddings, cache retrieved context, and semantic caching (serve a stored answer when a new query is "close enough"), with a staleness policy.
- Cost & capacity: a model call costs orders of magnitude more than a DB read; GPUs are the bottleneck → batch, route easy queries to a smaller model, cap output tokens.
- Freshness = replication/invalidation: re-embed when source docs change.
- Failure: models time out, degrade, hallucinate → timeouts + fallback to a smaller model, graceful degradation to a non-AI path; ground answers + cite to fight hallucination.

▸ Pitfalls
Treating it as an ML question (it's not — no attention math); forgetting distributed-systems basics once "AI" appears; ignoring cost/quality monitoring; leaking private prompts/embeddings.

▸ Interview probes
Walk the RAG pipeline; keep the vector index fresh; semantic caching and when it's unacceptable; handle model timeouts and cost; treat the model as an unreliable, expensive dependency.

🔗 Further reading
• pgvector — Postgres vector search: https://github.com/pgvector/pgvector
• Anthropic — build with Claude / retrieval: https://docs.anthropic.com/
• ByteByteGo — LLM / RAG system design (YouTube): https://www.youtube.com/@ByteByteGo
