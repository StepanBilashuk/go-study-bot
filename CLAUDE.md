Prep Bot — Full Project Spec
Personal interview-prep tracker. Telegram bot, Go, PostgreSQL, deployed on a 2 vCPU / 2 GB Ubuntu VPS.

This file is the complete brief. Drop it in the repo root and give it to Claude Code.


Table of contents
Hard constraints
Who this is for
Architecture
Database schema
YAML shapes
The five-stage ladder
Commands
Topics — algorithms
Topics — system design
Process drills
Resources map
Companies
Prompts — bot internal
Prompts — Claude Code build steps
Deployment
Roadmap and success criterion


1. Hard constraints
Violating any of these defeats the purpose of the project.

Phase 1 must be buildable in one weekend. No FSRS, no gamification, no admin UI, no web frontend, no Docker. If a feature is not in Phase 1, do not build it — not even groundwork.
Definitions live in YAML in the repo. State lives in Postgres. Never store topics, companies, prompts or resources in the database.
/today never returns more than 3 topics. This is the single most important product rule. The entire project exists to remove decision paralysis; a longer list recreates it.
No job scraping, no auto-applying, no LinkedIn integration.
Single binary, systemd service, Telegram long polling. No inbound ports.
Minimal dependencies: pgx, yaml.v3, a Telegram Bot API library, an HTTP client for the Anthropic API. Anything else requires explicit approval.
No frameworks, no DI containers, no repository layers, no clean architecture. Flat packages, direct SQL. This is a ~1500-line personal project.
Tests only for the spaced-repetition scheduler and for parsing Claude's JSON responses.


2. Who this is for
Backend engineer, 4 years commercial experience. NestJS / TypeScript / Node, Go in production, PostgreSQL, Redis, RabbitMQ, Kubernetes, GCP and AWS. Based in Tallinn. Preparing for senior backend roles at European product companies, main application wave March–April 2027, calibration interviews September–October 2026.

Diagnosed weaknesses this bot exists to attack:

Weakness
Evidence
How the bot attacks it
Does not self-initiate load estimation
System design mock scored 4.5/10
Daily estimation drills; mock brief forbids the interviewer from prompting
Does not self-catch contradictions in own design
Same mock
Daily contradiction drills
English ~B1.4, target B2
Preply assessment
All content, debriefs and mocks in English; glossary of terms he could not produce
Kafka / log-based systems
Highest-ROI gap across every target company
DDIA ch.11 is priority 1 in the SD track
STAR bank has no conflict / failure / trade-off story
Story audit
Weekly story-mining question (Phase 3)

Existing strengths to protect and surface, not re-learn:

Solo end-to-end ownership of a logistics platform backend and DevOps
HIPAA / FHIR healthcare platform on GCP
Card issuance in crypto-fintech
Three anchor STAR stories with numbers: notification delivery 85% → 99%+, a silent 10 000-record cap that was truncating clinical search, a saga+outbox library built after a provider outage lost ~100 card creations

Most candidates cannot produce measurable results. He can. The bot must require a metric on every new story.


3. Architecture
Layer
Component
Tech
Responsibility
Phase
runtime
bot binary
Go, single static binary
All logic: long-polling loop + scheduler goroutine
1
runtime
supervisor
systemd
Restart=always, EnvironmentFile for secrets
1
data
definitions
YAML in git
Topics, drills, resources, companies, prompts
1
data
state
PostgreSQL (apt, local)
Progress, sessions, drill log, stories, glossary
1
data
loader
yaml.v3
Reads data/** on start and /reload; invalid YAML prevents startup
1
integration
Telegram
Bot API, long polling
All user interaction
1
integration
Claude API
HTTP client
Calibration, debrief extraction, drills, mock briefs
1
logic
scheduler
goroutine
Daily push, next_due recalculation
1
logic
stage machine
Go
Enforces the 5-stage ladder and gates
1
logic
drill engine
Go + Claude API
Serves the weakest of 4 drill kinds
1
logic
readiness scorer
Go
Weighted coverage of company required_topics
4
ops
backup
cron + pg_dump -Fc
Daily dump, weekly copy off-box
1
ops
firewall
ufw + fail2ban
OpenSSH only, password auth off
1
ops
deploy
Makefile
build → scp → systemctl restart
1

Deliberately not built: web frontend (Telegram is the whole UI), Docker (pointless layer on one service), FSRS (hardcoded intervals suffice), admin CRUD UI (YAML + Claude Code replaces it).


4. Database schema
Mutable state only. Everything else is YAML.

topic_progress (

  topic_slug   text primary key,

  stage        int,          -- 0..4

  confidence   int,          -- 1..5

  attempts     int,

  last_touched timestamptz,

  next_due     timestamptz

);

drill_log (

  id         bigserial primary key,

  drill_slug text,           -- estimation | contradiction | clarify | next-step

  date       timestamptz,

  outcome    text,

  score      int             -- 1..5

);

sessions (

  id             bigserial primary key,

  date           timestamptz,

  raw_text       text,       -- exactly what was typed into /debrief

  extracted_gaps jsonb,

  topics_touched text[]

);

stories (                    -- phase 3

  id           bigserial primary key,

  title        text,

  situation    text,

  task         text,

  action       text,

  result       text,

  metrics      text[],       -- REQUIRED: reject a story with no number

  tech_tags    text[],

  competencies text[],       -- ownership, mentoring, conflict, failure, tradeoff

  strength     int,          -- 1..5, moved by behavioral mock results

  versions     jsonb,        -- {60s, 3min, technical}

  last_used    timestamptz

);

glossary (

  id       bigserial primary key,

  term     text,             -- English term he could not produce under pressure

  context  text,

  added_at timestamptz

);

xp_events (                  -- phase 5, behind a feature flag

  id         bigserial primary key,

  kind       text,           -- topic_closed | debrief | boss | drill

  points     int,

  created_at timestamptz

);


5. YAML shapes
data/

  topics/algorithms.yaml

  topics/system-design.yaml

  drills/process-drills.yaml

  companies/*.yaml

  resources.yaml

  books.yaml

prompts/

  *.yaml

topics/*.yaml

- slug: sliding-window

  track: algorithms          # algorithms | system-design

  name: Sliding Window

  priority: 3                # lower = earlier

  depends_on: [two-pointers]

  gate: "Fixed and variable window from memory + 3 unseen mediums"

  est_hours: 5

drills/process-drills.yaml

- slug: estimation

  kind: estimation           # estimation | contradiction | clarify | next-step

  name: Back-of-envelope estimation

  duration_min: 5

  prompt: prompts/drill-estimation.yaml

books.yaml — edition is not decorative: Russian translations paginate differently from the originals.

ddia:

  title: "Высоконагруженные приложения"

  author: Kleppmann

  edition: "ru-2018"

grokking:

  title: "Грокаем алгоритмы"

  author: Bhargava

  edition: "ru-2017"

resources.yaml

- topic: partitioning

  stage: 0                   # which ladder stage this serves

  type: book                 # book | video | article

  source: ddia

  chapter: 6

  section: "Partitioning"

  pages: null                # fill by hand from your edition

  est_min: 300

companies/*.yaml

slug: aiven

name: Aiven

locations: [Helsinki]

stack: [Go, Python, Kafka, PostgreSQL, Redis]

interview_process:

  - {stage: 1, name: Recruiter screen, format: call, duration_min: 30}

required_topics: [kafka-log-based, replication, partitioning, storage-retrieval]

values: [ownership, autonomy]

referral: false

researched_at: 2026-07-27

confidence: medium


6. The five-stage ladder
Every content topic climbs the same ladder. Gates are enforced — the bot refuses to advance a topic whose gate is unmet.

Stage
Algorithms
System design
Gate
0 Learn
study the pattern, write the template by hand
read the DDIA chapter + one video
reproduce the skeleton from memory
1 Guided
3–5 NeetCode problems, hints allowed
apply the block in a small design
solved with ≤1 hint
2 Quiz
pattern recognition, 60 s per item, no solving
trade-offs, numbers, failure modes
8/10 correct
3 Solo
3 unseen LeetCode mediums, 25 min timer, no hints
one full problem using the framework
2/3 in time
4 Review
paste own code to Claude for review
full mock with rubric
debrief logged

After stage 4 the topic enters spaced repetition and returns at stage 3, never at stage 0.

Phase 1 intervals are hardcoded: 1, 3, 7, 21 days.

Why stage 2 matters most. Interviews are failed on recognition, not implementation — the candidate can write a sliding window but does not see that the problem calls for one. Nothing off-the-shelf trains this in isolation, which is why the bot must generate it.

Claude review must not become "give me the answer." Rule: only after an honest attempt, and the first message pastes his own code. Ask for four things — missed edge cases, complexity check, Go idiomaticity, and what an interviewer would have probed here. The fourth is the only one an editorial cannot provide.


7. Commands
Command
Phase
Behaviour
/start
1
One-time calibration: ~40 topics, "explain in two sentences", Claude scores each and sets initial confidence
/today
1
3 topics (algorithms, system design, review) + 1 drill + max 2 resources per topic
/done <slug>
1
Advance a stage if the gate is met, otherwise name what is missing
/drill
1
One process drill of the kind with the weakest recent scores
/debrief
1
Free text in; Claude extracts gaps, updates confidence, reschedules
/boss
1
Check readiness, then emit a paste-ready mock brief
/reload
1
Re-read all YAML from disk
/quiz <topic>
2
Pattern-recognition quiz, 10 items, 60 s each
/story
3
Weekly story-mining question
/stories
3
Competency matrix — which competencies still have no story
/newcompany <name>
4
Emit a research prompt to paste into Claude with web search
/importcompany
4
Import the returned JSON, validate, write companies/<slug>.yaml
/ready
4
Readiness score per company with named blocker
/stats
5
XP, per-track levels, mastery bars, book coverage, streak

Boss gate: at least 5 topics in the block at confidence ≥ 4, none older than 10 days. Until then /boss shows what is missing — e.g. "to unlock the SD boss: kafka-log-based 2/5, partitioning 3/5".


8. Topics — algorithms
NeetCode order. Do not reshuffle. Realistic throughput at this depth is one topic per week.

#
Slug
Name
Depends on
Gate
Hrs
1
arrays-hashing
Arrays & Hashing
—
3 unseen mediums in 25 min each
5
2
two-pointers
Two Pointers
arrays-hashing
template from memory + 3 mediums
4
3
sliding-window
Sliding Window
two-pointers
fixed + variable window from memory
5
4
stack
Stack
arrays-hashing
monotonic stack from memory
4
5
binary-search
Binary Search
arrays-hashing
boundary-safe template from memory
5
6
linked-list
Linked List
two-pointers
reverse, cycle detect, merge, no hints
4
7
trees
Trees
linked-list
DFS/BFS traversals from memory
6
8
heap
Heap / Priority Queue
trees
top-K and merge-K without hints
4
9
backtracking
Backtracking
trees
subsets, permutations, combination sum
5
10
graphs
Graphs
trees, heap
BFS/DFS/topological sort + Dijkstra
8
11
dynamic-programming
Dynamic Programming
graphs
1D and 2D DP, 3 unseen mediums
10

Marginal return collapses after category 9. Only Bolt's algorithmic gauntlet justifies deep graphs and DP. Do not grind volume for the feeling of progress — 15 topics at stage 4 beats 60 at stage 1.


9. Topics — system design
#
Slug
Name
Book
Gate
Hrs
1
kafka-log-based
Stream processing, log-based brokers, Kafka, CDC
DDIA ch.11
explain exactly-once semantics and CDC unprompted
8
2
storage-retrieval
LSM trees, B-trees, indexes
DDIA ch.3
explain when LSM beats B-tree and why
6
3
replication
Replication, lag, multi-leader, leaderless
DDIA ch.5
name 3 replication-lag anomalies and their fixes
6
4
partitioning
Partitioning, consistent hashing, rebalancing
DDIA ch.6
explain rebalancing and hot-key handling
6
5
transactions
Transactions, isolation levels, race conditions
DDIA ch.7
name each isolation level and the anomaly it prevents
6
6
encoding-evolution
Encoding and schema evolution
DDIA ch.4
forward vs backward compatibility with an example
4
7
distributed-trouble
Unreliable clocks, networks, partial failure
DDIA ch.8
explain why wall-clock ordering is unsafe
5
8
consistency-consensus
Linearizability, CAP, consensus
DDIA ch.9
state CAP precisely, not the folk version
6
9
caching
Caching strategies and invalidation
—
write-through/behind/aside with failure modes
4
10
rate-limiting
Rate limiting and backpressure
—
token bucket vs leaky bucket, distributed counters
3
11
idempotency
Idempotency and exactly-once at API level
—
design an idempotent payment endpoint unprompted
4
12
geo-distribution
Geo-distribution and multi-region
—
latency/consistency trade-off per region
4

Chapter 11 goes first. Kafka is the top gap and Aiven is a target employer.


10. Process drills
The system design failure was process, not knowledge. Process is a habit, and habits are trained by isolated repetition of the sub-skill — not by whole problems once a week.

One drill every day, 3–5 minutes, regardless of which topic is scheduled.

Slug
Trains
Example
estimation
Self-initiating load estimation — the #1 weakness
"A food delivery app in one city. Do NOT design it. Produce only DAU, peak QPS, storage per year, instance count, and your assumptions. 5 minutes."
contradiction
Self-catching inconsistencies — the #2 weakness
A 6-line design with 1–2 planted contradictions (synchronous writes to one Postgres primary + "handles 50k RPS"; sharded by user_id + the main query is a time-range scan). Find them in 90 seconds.
clarify
The first 3 minutes of an interview, where half the impression forms
"Design a notification system." Produce 5 questions that narrow the problem. 2 minutes.
next-step
Making the framework sequence automatic under pressure
A truncated interview transcript. What is the correct next step, and why?

Drills produce reflexes. The weekly boss remains mandatory — 45 uninterrupted minutes in English under pushback is a separate skill that drills do not cover.

Rhythm: drill daily (5 min) · content topic on schedule · one full problem weekly (45 min + 30 min debrief).


11. Resources map
Attach resources to topic + stage, not just topic. Stage 0 wants a chapter and a video; stage 4 wants an article on pitfalls.

Never show more than 2 resources. A list of ten links recreates the paralysis the bot exists to remove.
DDIA → topics
Chapter
Topic
Priority
11
kafka-log-based
critical
3
storage-retrieval
high
5
replication
high
6
partitioning
high
7
transactions
high
4
encoding-evolution
medium
8
distributed-trouble
medium
9
consistency-consensus
medium
Grokking Algorithms → topics
Only four chapters are worth his time; chapters 1–5 are below his level.

Chapter
Topic
6 (BFS)
graphs
7 (Dijkstra)
graphs
8 (greedy)
dynamic-programming
9 (DP)
dynamic-programming
To fill by hand
Page numbers — pagination differs by edition; record edition in books.yaml, keep chapter + section as the stable reference
Video URLs — add as you actually watch, not in bulk, or the bank becomes a landfill. Suggested channels: NeetCode for algorithms; ByteByteGo, Hussein Nasser, Arpit Bhayani for system design
Rate each resource after use; quality scores let the useful ones surface over time


12. Companies
Slug
Company
Location
Stack
Process
Required topics
Referral
Status
wolt
Wolt
Helsinki / Stockholm
Kotlin, Python, Go
take-home
arrays-hashing, trees, partitioning, caching
yes
Wave 1, highest probability
bolt
Bolt
Tallinn
Go, Java
5 stages, algorithmic gauntlet
graphs, dynamic-programming, partitioning, geo-distribution
yes
Wave 1, 50/50 — algo tail is the blocker
aiven
Aiven
Helsinki
Go, Python, Kafka, PostgreSQL, Redis
standard loop
kafka-log-based, replication, partitioning, storage-retrieval
no
Best domain fit; Kafka gap is the blocker
verda
Verda
Helsinki
TypeScript, NestJS
standard
idempotency, caching, transactions
no
Strongest stack match, ready soonest
enfuce
Enfuce
Espoo
card issuing, fintech
standard
idempotency, transactions, kafka-log-based
no
Direct match to his card-issuance work
attio
Attio
London / EU remote
TypeScript, Node, GCP
EM 30 min → CTO 45 min
idempotency, caching, storage-retrieval
no
EU-remote track, €95–125K
pipedrive
Pipedrive
Tallinn
—
cognitive test + home challenge + presentation
arrays-hashing, trees, caching
no
Blocker is English, not algorithms
stripe
Stripe
Dublin
Ruby, Java, Go
high bar
idempotency, transactions, consistency-consensus, kafka-log-based
no
Stretch, 2027+

Readiness score = weighted coverage of required_topics by current confidence.

Verda    85%  → ready

Aiven    70%  → blocker: kafka-log-based

Bolt     55%  → blocker: graphs, dynamic-programming

Stripe   40%  → blocker: consistency-consensus

When an interview date is set, /today reorders around that company's profile for N days.

Before investing effort in Attio-type EU-remote roles, verify Estonia is covered by their EOR. Remote-first startups cover a limited set of jurisdictions.


13. Prompts — bot internal
Store each as a YAML file under prompts/. Read from disk at runtime so they can be tuned without redeploying.
calibration.yaml — /start, once
You are assessing a backend engineer's knowledge to calibrate a study plan.

For each topic below, ask the candidate to explain it in two sentences, or

to say when they would apply it. Ask ONE topic at a time.

Score each answer 1-5:

1 = no working knowledge

2 = recognises the term only

3 = can define it, cannot apply it

4 = can apply it, has used it

5 = can teach it and knows the trade-offs

Do not teach. Do not correct. Only assess.

Return strict JSON, no prose:

{"scores":[{"topic":"<slug>","confidence":1-5,"note":"<max 12 words>"}]}
debrief.yaml — the core loop
The single most important prompt in the project. Without it the bot is a static checklist.

You process a study debrief from a backend engineer preparing for interviews.

Input is free text, e.g. "solved sliding window, two went fine, struggled

with variable-size window".

Extract what actually happened. Do not encourage, do not advise.

Known topic slugs: {{topic_slugs}}

Return strict JSON only:

{

  "topics_touched": ["<slug>"],

  "gaps": [{"topic":"<slug>","detail":"<specific>","severity":"low|medium|high"}],

  "confidence_updates": [{"topic":"<slug>","new_confidence":1-5}],

  "new_glossary_terms": [{"term":"<english term>","context":"<where>"}]

}

Rules:

- Only use slugs from the provided list

- Lower confidence on any topic where a gap was reported

- A topic mentioned without difficulty gets +0, not +1

- new_glossary_terms only for English words the candidate could not produce
mock-sd.yaml — /boss, system design
You are a senior engineer at {{company}} running a 45-minute system design

interview. Conduct it entirely in English.

Candidate: backend engineer, 4 years, NestJS/TypeScript/Go, fintech and

healthcare background, has run production systems solo.

Task: {{task}}

STRICT RULES:

- Do NOT prompt the candidate to do load estimation. If they do not

  initiate it themselves, note it silently and move on.

- Do NOT point out contradictions in their design. Note them silently.

- Push back on vague answers. Ask "why" at least three times.

- Deep-dive on: {{weak_topics}}

- Stay in role. Do not teach during the interview.

At the end score 1-10 on each: self-initiated estimation, self-caught

contradictions, depth on one component, failure-mode awareness, English

clarity under pressure.

Then quote the three weakest moments verbatim.

The first rule is load-bearing. His failure mode is not initiating estimation. An interviewer who prompts produces a falsely good score and the bot learns the wrong thing.
mock-behavioral.yaml — /boss, behavioral
You are a hiring manager at {{company}} running a 30-minute behavioral

interview in English. Company values: {{values}}.

Ask 4 questions, one at a time, probing these competencies: {{competencies}}.

After each answer ask at least two follow-ups that dig into the candidate's

specific contribution.

Watch for and flag:

- "we" instead of "I"

- missing measurable result

- STAR structure collapsing under follow-up

Score 1-10 on: structure, specificity, measurable outcome, ownership signal,

English clarity. List every English term the candidate visibly reached for

and could not find.
drill-estimation.yaml
Give the candidate ONE system and ask only for back-of-envelope numbers.

Explicitly forbid designing anything.

Required outputs: DAU, peak QPS, storage per year, rough instance count,

and the assumptions behind each.

Do not provide any numbers yourself. Do not hint.

After the answer, score 1-5 on: did they state assumptions, did they

sanity-check magnitudes, did they cover all four outputs.
drill-contradiction.yaml
Write a 6-line system design description containing exactly 1-2 deliberate

internal contradictions. Make them realistic, not absurd.

Examples of the kind of contradiction to plant:

- synchronous writes to a single Postgres primary + "handles 50k RPS"

- sharded by user_id + primary access pattern is a time-range scan

- "eventually consistent" + a read-your-own-writes requirement

Present it. Give the candidate 90 seconds. Then reveal and score 1-5.
company-research.yaml — /newcompany, pasted into Claude with web search
Research {{company}} as a backend engineering employer. Use web search.

Return ONLY valid JSON, no prose.

{

  "name":"", "locations":[], "remote_policy":"",

  "visa_sponsorship":"yes|no|unknown", "stack":[],

  "interview_process":[{"stage":1,"name":"","format":"","duration_min":0}],

  "known_question_themes":[], "values":[],

  "salary_band_eur":{"min":0,"max":0,"source":""},

  "recent_news":[], "sentiment_summary":"",

  "required_topics":[], "new_topics":[],

  "confidence":"high|medium|low", "sources":[]

}

Available topic slugs: {{topic_slugs}}

Rules:

- required_topics MUST use only the provided slugs

- anything not covered goes into new_topics as free text

- salary_band must cite a source; if none found use nulls

- set confidence to low if interview-process data is older than 12 months

Without the slug constraint you get orphan strings like "distributed systems knowledge" that map to nothing, and the readiness score silently breaks.

Store researched_at and confidence per company; flag records older than 3 months as stale. Web data on interview processes and salaries is a hint, not a fact — one conversation through a referral beats an hour of agent research.


14. Prompts — Claude Code build steps
Paste one per step. Commit after each. Run the bot after each.
Step 0 — setup
You are building me a personal Telegram bot for technical interview prep.

The full spec is in CLAUDE.md at the repo root - read it fully before starting.

Language: Go. Database: PostgreSQL. Deploy: single binary + systemd on an

Ubuntu VPS with 2 vCPU / 2 GB RAM.

WORKING RULES - follow strictly:

1. We build ONLY Phase 1 from the spec. Anything marked Phase 2-5 is out of

   scope, including "groundwork" abstractions and interfaces for later.

2. Minimal dependencies. Allowed: pgx, yaml.v3, a Telegram Bot API library,

   an HTTP client for the Anthropic API. Ask me before adding anything else.

3. No frameworks, no DI containers, no repository layers, no clean

   architecture. Flat packages, direct SQL. This is a 1500-line personal

   project, not an enterprise system.

4. Tests only for the spaced-repetition scheduler and for parsing Claude's

   JSON responses. Nothing else.

5. Work step by step. After each step, stop, show what you did, and wait for

   my "ok". Never chain multiple steps.

6. If anything in the spec is ambiguous, ASK. Do not decide for me.

Confirm you have read the spec and list the Phase 1 steps you propose.

Do not write code yet.
Step 1 — skeleton
Step 1. Project skeleton:

- directory structure per the spec (cmd/, internal/, data/, prompts/)

- go.mod, config from environment variables

- migrations for every table in the Postgres section (plain .sql files run at

  startup - no goose, no migrate)

- YAML loader: reads data/**, validates, returns in-memory structs; a

  validation error must prevent startup with a clear message

- Telegram stub: bot starts and replies "pong" to /ping

Nothing else. Show me the file tree and the key file contents.
Step 2 — seed data
Step 2. Populate data/ from sections 8-12 of the spec.

Produce:

- data/topics/algorithms.yaml and data/topics/system-design.yaml

- data/drills/process-drills.yaml

- data/books.yaml

- data/resources.yaml

- data/companies/<slug>.yaml for each company

Do not invent page numbers - write null, I will fill them from my own

edition. Same for video URLs.
Step 3 — commands
Step 3. Implement the Phase 1 commands: /start, /today, /done, /drill,

/debrief, /boss, /reload. Behaviour is defined in section 7.

Claude prompts live in prompts/*.yaml and are read from disk - do not

hardcode them in Go.

Parse Claude responses strictly as JSON against the schemas in section 13.

On invalid JSON: one retry with a clarifying instruction, then a clear error

to the user.

Spaced repetition intervals are hardcoded: 1, 3, 7, 21 days.

Enforce the rule that /today never returns more than 3 topics.
Step 4 — deployment
Step 4. Deployment artifacts:

- Makefile: build for linux/amd64, deploy via scp + systemctl restart

- systemd unit with Restart=always and EnvironmentFile

- .env.example listing every variable

- server bootstrap script: postgres from apt, ufw (OpenSSH only), fail2ban,

  disable SSH password auth

- cron entry for daily pg_dump -Fc to /var/backups

- README: bring it up from scratch in 15 minutes

Secrets only via EnvironmentFile with mode 600. Nothing secret in the repo.


15. Deployment
Ubuntu VPS, 2 vCPU / 2 GB. A Go binary plus Postgres with this workload uses well under half the RAM.
Verify first — before writing any code
curl -sI https://api.telegram.org  | head -1

curl -sI https://api.anthropic.com | head -1

If either hangs or fails, the region blocks them and the instance must be recreated elsewhere. Find this out now, not after building.
Hardening
sudo sed -i 's/^#*PasswordAuthentication.*/PasswordAuthentication no/' /etc/ssh/sshd_config

sudo systemctl restart ssh

sudo ufw default deny incoming

sudo ufw allow OpenSSH

sudo ufw enable

sudo apt install -y fail2ban unattended-upgrades

A public IP is scanned by bots within hours. The bot itself needs no inbound port — long polling is outbound only.
Layout
Postgres from apt, bound to 127.0.0.1
Go binary as a systemd service, Restart=always
Secrets in /etc/prepbot.env, mode 600, via EnvironmentFile — never in the repo, never in code
No nginx, no exposed ports beyond SSH
Backups
pg_dump -Fc prepbot > /var/backups/prepbot-$(date +\%F).dump

Daily via cron, weekly copy off-box. The accumulated progress history is the valuable artifact, not the code — six months of "what you closed, where you failed, how confidence moved" cannot be regenerated.

The instance is paid through June 2027. With pg_dump and one binary, migrating to any other VPS takes half an hour, so avoid anything Tencent-specific.


16. Roadmap and success criterion
Phase
Scope
Build time
Gate to next
1
Core loop: schema, YAML loader, all Phase 1 commands, hardcoded intervals, deploy
1 weekend
use it for 3 weeks
2
Quiz generator, contradiction-drill generator
1 evening
2 weeks of real use
3
STAR bank, competency matrix, weekly story-mining, behavioral boss
1 evening
2 weeks of real use
4
Company research import, readiness score, interview-date replanning
1 evening
2 weeks of real use
5
Gamification: XP, per-track levels, mastery bars, streak with freezes, book coverage — behind a feature flag
1 evening
—

Wait two weeks between phases and actually use what exists.
The real risk
Building this is more enjoyable than using it. FSRS, inline keyboards, deployment — engaging. NeetCode in Go at 22:00 after a full workday — not. The realistic failure mode is six weekends spent on the tool, a beautiful system by October, and zero problems solved — while it subjectively feels like preparation.

Phase 1 is capped at one weekend for exactly this reason.
Success criterion — check at week 3
Two numbers: study sessions logged and mocks completed.

Fewer than 15 sessions and fewer than 2 mocks → the bot did not work. Fall back to a markdown file and a calendar reminder, without regret.
More than that → build the next phase.

This criterion tests whether the system works on this particular person, which no amount of design reasoning can establish in advance.
What the bot does and does not do
Moves the needle: decision paralysis (8/10), turning mock failures into a concrete next-day plan (9/10), forgetting (7/10), knowing which company you are ready for today (8/10), closing STAR gaps by surfacing what is already there (8/10).

Does not: create hours, raise English (that is speaking volume), or fix the system design process by itself — the daily drills do that, and the weekly boss verifies it.

Honest estimate: ~40% improvement on the system design track, 20–30% overall. The other 70% is hours that still have to be sat through. The bot's real contribution is raising the probability that they are — roughly from 40% to 70% over a 35-week run.
