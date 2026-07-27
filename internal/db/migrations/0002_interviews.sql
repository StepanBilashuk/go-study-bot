-- Phase 4: interview-date replanning. When a user has an upcoming interview,
-- /today reorders around that company's required_topics until the date passes.
create table if not exists interviews (
    user_id      bigint not null references users(user_id) on delete cascade,
    company_slug text   not null,
    date         date   not null,
    primary key (user_id, company_slug)
);
create index if not exists interviews_user_idx on interviews (user_id);
