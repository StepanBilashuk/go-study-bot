📖 Search & indexing (system design)

Full-text search is powered by an inverted index: term → list of documents containing it. Query = intersect/union the postings lists.

▸ Pipeline
Analyze text (tokenize, lowercase, remove stop words, stem) → build postings → rank by relevance (TF-IDF or BM25). Elasticsearch/OpenSearch wrap Lucene.

▸ Scaling
Shard by document across nodes (scatter-gather queries), replicate for availability, near-real-time refresh. Separate the write (index) path from the read (search) path.

▸ Pitfalls
Index size and refresh lag vs consistency; analyzer choices change results; relevance tuning; deep pagination.

▸ Interview probes
Inverted index; tokenization/analysis; TF-IDF vs BM25; how you shard and scale search; autocomplete (trie or edge n-grams).

🔗 Further reading
• ByteByteGo — search / Elasticsearch (YouTube): https://www.youtube.com/@ByteByteGo
• Elastic — inverted index & relevance: https://www.elastic.co/guide/en/elasticsearch/reference/current/index.html
• ByteByteGo newsletter: https://blog.bytebytego.com
