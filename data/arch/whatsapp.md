🏛 WhatsApp — billions of connections on Erlang

▸ The gist
Delivered messaging to hundreds of millions of users with a shockingly small team, on Erlang/BEAM — a runtime built for millions of lightweight concurrent processes and hot code swaps.

▸ Patterns to learn
- Right tool for concurrency: Erlang's per-connection lightweight processes + supervision trees → one server holds huge numbers of persistent connections.
- Keep it lean: minimal services, careful tuning of the network stack (millions of TCP connections per box).
- Store-and-forward for offline delivery; per-conversation ordering; acks (sent/delivered/read).

▸ Maps to
/learn concurrency · /learn websockets-realtime · /design chat

🔗 Read the real thing
• HighScalability — WhatsApp architecture: https://highscalability.com/
• (Erlang at scale) https://www.erlang-factory.com/
