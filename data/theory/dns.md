📖 DNS & name resolution (system design)

Turns a name into an IP — the internet's directory, and a quiet but powerful traffic-steering layer.

▸ Resolution path
Client → recursive resolver → root → TLD (.com) → authoritative server. Results are cached along the way per TTL.

▸ Records & steering
A/AAAA (IP), CNAME (alias), MX (mail), TXT, NS. TTL controls cache freshness. GeoDNS / latency-based / weighted routing steer users to the nearest or healthiest region (also used for failover and blue-green).

▸ Pitfalls
TTL vs failover speed (low TTL = faster failover but more queries); propagation delay; clients/OS caching stale records; DNS as an availability dependency.

▸ Interview probes
The resolution path; record types; TTL trade-offs; how GeoDNS / weighted routing enables multi-region and failover.

🔗 Further reading
• ByteByteGo — how DNS works (YouTube): https://www.youtube.com/@ByteByteGo
• Cloudflare — what is DNS: https://www.cloudflare.com/learning/dns/what-is-dns/
• ByteByteGo newsletter: https://blog.bytebytego.com
