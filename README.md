# prepbot

Personal interview-prep tracker. A Telegram bot in Go, backed by PostgreSQL,
that removes decision paralysis: `/today` never gives you more than 3 topics.

Definitions (topics, drills, resources, companies, prompts) live in YAML in this
repo. Only mutable state (progress, sessions, drill log, glossary) lives in
Postgres. See `CLAUDE.md` for the full spec.

**No API key.** The bot never calls the Anthropic API. AI commands **emit a
prompt** you run in your own Claude, then **paste the JSON reply back** here —
the bot parses and applies it (like `/newcompany` → `/importcompany`). Zero cost,
no key to manage. `/cancel` drops a pending paste-back.

**Multi-user:** one bot instance serves many people. Every stateful row is
scoped to the Telegram user id, so each user has independent progress,
calibration, drills, and glossary. Users are registered automatically on first
message.

## What it does

| Command | Behaviour |
|---------|-----------|
| `/start` | Emit a calibration prompt for all topics → paste the scores JSON back to set confidence |
| `/today` | 3 topics (algorithms · system design · review) + 1 drill + ≤2 resources each |
| `/done <slug>` | Advance a topic one stage if prerequisites are met; enters spaced repetition after stage 4 |
| `/drill` | Emit a drill of your weakest kind → paste the `{score,outcome}` JSON back |
| `/debrief <text>` | Emit an extraction prompt with your text → paste the gaps JSON back |
| `/boss [behavioral]` | Checks readiness, then emits a paste-ready mock-interview brief |
| `/quiz <slug>` | Emit a 10-item recognition quiz (run it in Claude; ≥8/10 meets the stage-2 gate) |
| `/story` | Emit a story-mining prompt for your weakest competency → paste the STAR JSON back (metric required) |
| `/stories` | Competency matrix — which competencies still have no story |
| `/ready` | Readiness score per company with its named blocker |
| `/newcompany <name>` | Emit a research prompt to paste into Claude with web search |
| `/importcompany` | Import the returned JSON → writes `data/companies/<slug>.yaml` |
| `/interview <slug> <YYYY-MM-DD>` | Set an interview date; `/today` reorders around that company (`clear` to remove) |
| `/glossary` | Review the English terms you couldn't produce (collected by `/debrief`) |
| `/progress` | Compact per-track summary (mastered · in progress · not started) |
| `/stats` | XP, level, streak (with freezes), mastery bars, book coverage (needs `GAMIFICATION=true`) |
| `/push on\|off` | Enable/disable the daily plan push for yourself |
| `/whoami` | Show your chat id |
| `/help` | Full command list |
| `/reload` | Re-read all YAML from disk (bad YAML is rejected, running set untouched) |
| `/ping` | Liveness check → `pong` |

`/boss behavioral` targets the competencies you have **no story for yet**, so
the mock probes exactly your STAR-bank gaps.

Spaced-repetition intervals are hardcoded: **1, 3, 7, 21 days**.

Phases 1–5 are all built. Gamification (`/stats`, XP) sits behind the
`GAMIFICATION` flag, off by default.

## Run with Docker Compose (DB + bot together)

The simplest way to bring up everything:

```bash
cp .env.example .env          # fill TELEGRAM_BOT_TOKEN + POSTGRES_PASSWORD (no API key needed)
docker compose up --build -d
docker compose logs -f bot
```

Compose starts Postgres (with a persistent `pgdata` volume), waits for it to be
healthy, then starts the bot — which runs its migrations automatically. `data/`
and `prompts/` are mounted read-only, so editing YAML + sending `/reload` works
without a rebuild. Stop with `docker compose down` (add `-v` to wipe the DB).

## Run locally (no Docker)

Requires Go 1.26+ and a local Postgres.

```bash
createdb prepbot
cp .env.example .env          # fill TELEGRAM_BOT_TOKEN + DATABASE_URL (no API key needed)
set -a; source .env; set +a
go run ./cmd/prepbot
```

Migrations run automatically at startup. Message the bot `/ping` → `pong`.

```bash
make test    # unit tests (scheduler + Claude JSON parsing) with -race
make vet
```

## Deploy to a VPS (≈15 minutes)

Ubuntu, 2 vCPU / 2 GB. Long polling is outbound-only — no inbound port.

```bash
# On the server (as a sudo user):
./deploy/bootstrap.sh                 # postgres, ufw, fail2ban, ssh hardening, backup cron
sudo install -m 600 /dev/null /etc/prepbot.env && sudo nano /etc/prepbot.env
sudo cp deploy/prepbot.service /etc/systemd/system/ && sudo systemctl daemon-reload
sudo systemctl enable --now prepbot

# From your laptop:
make deploy REMOTE_HOST=<ip>          # build linux/amd64, scp binary + data, restart
journalctl -u prepbot -f              # (on the server) watch it
```

Secrets live only in `/etc/prepbot.env` (mode 600), never in the repo. A daily
`pg_dump -Fc` lands in `/var/backups`; copy it off-box weekly — the accumulated
progress history is the valuable artifact, not the code.

## Layout

```
cmd/prepbot        entrypoint: config → definitions → prompts → db+migrate → bot
internal/config    all config from environment variables
internal/db        pgx pool, startup .sql migrations, all queries
internal/definitions  YAML loader + validation (topics/drills/resources/companies)
internal/prompts   Claude prompt templates, rendered with {{vars}}
internal/claude    strict JSON parsers for the pasted-back replies (tested)
internal/scheduler spaced-repetition intervals (tested)
internal/tgbot     long-polling loop, command handlers, daily push
data/ prompts/     the editable content set
deploy/            systemd unit, bootstrap script
```
