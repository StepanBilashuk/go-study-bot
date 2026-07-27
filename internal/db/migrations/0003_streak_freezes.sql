-- Phase 5: streak freezes. Each user gets a small budget of missed days their
-- current streak may absorb without breaking (spec §16: "streak with freezes").
alter table users add column if not exists streak_freezes int not null default 2;
