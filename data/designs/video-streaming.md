🏗 Design a video platform / transcoding pipeline (YouTube / Netflix) [Bolt asked this]

▸ Requirements
Functional: upload a video, transcode to multiple resolutions/formats, stream with adaptive quality, play globally.
Non-functional: huge storage + bandwidth, smooth playback (low buffering), scale reads massively.

▸ Estimation
Upload is write-heavy + CPU-heavy (transcoding); playback is extremely read-heavy → CDN dominates delivery.

▸ Upload → transcode pipeline
Upload (chunked/multipart to blob storage) → enqueue a transcode job → workers split the video into segments, transcode each segment in PARALLEL into several bitrates (240p…4K) and formats (HLS/DASH) → write segments + a manifest to blob storage. A job/DAG orchestrator tracks progress.

▸ Delivery
Segments served via CDN (edge cache near users). The player fetches the manifest, then adaptively requests the bitrate that fits current bandwidth (adaptive bitrate streaming).

▸ Deep dives & trade-offs
Parallel segment transcoding (fan-out on a queue); storage cost (tiering, keep popular resolutions hot); CDN cache-miss on viral spikes; encode-on-demand vs pre-encode; DRM; metadata + search.

🔗 Further reading
• ByteByteGo — Netflix/YouTube design (YouTube): https://www.youtube.com/@ByteByteGo
• Netflix Open Connect (CDN): https://openconnect.netflix.com/
• System Design Primer: https://github.com/donnemartin/system-design-primer
