-- Mutable state only. All definitions (topics, drills, resources, companies,
-- prompts) live in YAML under data/ and prompts/ (spec §4).
--
-- Every stateful row is scoped to a user_id (the Telegram chat id), so many
-- people can use one bot instance with independent progress.

create table if not exists users (
    user_id    bigint      primary key,          -- telegram chat id
    username   text,
    first_seen timestamptz not null default now(),
    last_seen  timestamptz not null default now(),
    push_enabled boolean   not null default true  -- receive the daily /today push
);

create table if not exists topic_progress (
    user_id      bigint not null references users(user_id) on delete cascade,
    topic_slug   text   not null,
    stage        int    not null default 0,       -- 0..4
    confidence   int    not null default 1,       -- 1..5
    attempts     int    not null default 0,
    last_touched timestamptz,
    next_due     timestamptz,
    primary key (user_id, topic_slug)
);

create table if not exists drill_log (
    id         bigserial   primary key,
    user_id    bigint      not null references users(user_id) on delete cascade,
    drill_slug text        not null,              -- estimation | contradiction | clarify | next-step
    date       timestamptz not null default now(),
    outcome    text,
    score      int                                -- 1..5
);
create index if not exists drill_log_user_idx on drill_log (user_id);

create table if not exists sessions (
    id             bigserial   primary key,
    user_id        bigint      not null references users(user_id) on delete cascade,
    date           timestamptz not null default now(),
    raw_text       text        not null,          -- exactly what was typed into /debrief
    extracted_gaps jsonb,
    topics_touched text[]
);
create index if not exists sessions_user_idx on sessions (user_id);

-- Phase 3. Created now so migrations cover every table in the spec.
create table if not exists stories (
    id           bigserial primary key,
    user_id      bigint    not null references users(user_id) on delete cascade,
    title        text,
    situation    text,
    task         text,
    action       text,
    result       text,
    metrics      text[],                          -- REQUIRED: reject a story with no number
    tech_tags    text[],
    competencies text[],                          -- ownership, mentoring, conflict, failure, tradeoff
    strength     int,                             -- 1..5, moved by behavioral mock results
    versions     jsonb,                           -- {60s, 3min, technical}
    last_used    timestamptz
);
create index if not exists stories_user_idx on stories (user_id);

create table if not exists glossary (
    id       bigserial   primary key,
    user_id  bigint      not null references users(user_id) on delete cascade,
    term     text        not null,                -- English term he could not produce under pressure
    context  text,
    added_at timestamptz not null default now()
);
create index if not exists glossary_user_idx on glossary (user_id);

-- Phase 5, behind a feature flag.
create table if not exists xp_events (
    id         bigserial   primary key,
    user_id    bigint      not null references users(user_id) on delete cascade,
    kind       text        not null,              -- topic_closed | debrief | boss | drill
    points     int         not null,
    created_at timestamptz not null default now()
);
create index if not exists xp_events_user_idx on xp_events (user_id);
