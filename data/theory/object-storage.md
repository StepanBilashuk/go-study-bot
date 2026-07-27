📖 Object / blob storage (system design)

Store large unstructured data (images, video, backups, static assets) cheaply and durably.

▸ Object vs block vs file
- Object (S3, GCS): flat namespace of key→blob + metadata, HTTP API, massive scale, high durability (11 nines). No in-place edits.
- Block: raw disk volumes (EBS) for databases.
- File: POSIX shared filesystem (NFS/EFS).

▸ Mechanics
Multipart/chunked upload for large files; presigned URLs let clients up/download directly (offloading your servers); versioning; lifecycle rules (tiering to cold storage).

▸ Pitfalls
Small-file overhead; hot objects (add a CDN); egress cost; eventual vs strong consistency (S3 is now strong read-after-write).

▸ Interview probes
Object vs block vs file; multipart upload; presigned URLs; durability vs availability; serving media at scale (storage + CDN).

🔗 Further reading
• ByteByteGo — object storage / S3 (YouTube): https://www.youtube.com/@ByteByteGo
• AWS S3 — how it works: https://docs.aws.amazon.com/AmazonS3/latest/userguide/Welcome.html
• ByteByteGo newsletter: https://blog.bytebytego.com
