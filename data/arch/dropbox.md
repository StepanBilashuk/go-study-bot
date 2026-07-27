🏛 Dropbox — file storage & sync

▸ The gist
Sync files across devices reliably and efficiently. Split metadata (who has what, versions) from block storage (the actual bytes); dedup identical blocks. Famously migrated off S3 onto their own storage ("Magic Pocket").

▸ Patterns to learn
- Separate metadata service from block storage; content-addressed blocks (hash) enable dedup + integrity.
- Chunk files into blocks; sync only changed blocks (delta sync) → bandwidth-efficient.
- Own storage at scale beats cloud rent past a threshold (Magic Pocket); erasure coding for durability.
- Notification/sync protocol to push changes to clients.

▸ Maps to
/learn object-storage · /learn scaling-databases · /design distributed-cache (content-addressing ideas)

🔗 Read the real thing
• Dropbox Tech blog: https://dropbox.tech/
• (Magic Pocket, scaling to exabytes) https://dropbox.tech/infrastructure
