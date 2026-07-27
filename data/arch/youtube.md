🏛 YouTube — video platform at scale

▸ The gist
Upload → transcode into many resolutions/formats → store in blob → serve globally via CDN with adaptive bitrate. Massive read fan-out; metadata + recommendations on the side.

▸ Patterns to learn
- Transcoding pipeline: split a video into segments, transcode in PARALLEL to multiple bitrates (HLS/DASH), write to blob storage.
- CDN-first delivery (Google's edge / Open Connect-style); adaptive bitrate streaming from a manifest.
- Metadata in a scalable store; view counts via async event aggregation (eventually consistent).
- Thumbnails, search, recommendations as separate services.

▸ Maps to
/design video-streaming · /learn cdn-edge · /learn object-storage · /learn message-queues

🔗 Read the real thing
• ByteByteGo — Design a system like YouTube: https://blog.bytebytego.com/p/ep130-design-a-system-like-youtube
• ByteByteGo newsletter: https://blog.bytebytego.com/
