-- migrate:up
create type handedness as enum ('right', 'left', 'ambidextrous');

create table player(
    id uuid primary key not null default  uuid_generate_v4(),
    height numeric,
    weight numeric,
    handedness handedness,
    racket varchar(250),
    matches_expected integer not null default 0,
    matches_played integer not null default 0,
    matches_won integer not null default 0,
    matches_scheduled integer not null default 0,
    seasons_played integer not null default 0,
    elo integer not null default 1200,
    highest_elo integer not null default 1200,
    is_elo_provisional boolean not null default true,
    account_id uuid not null references account (id) on delete cascade,
    current_league_id uuid references league (id),
    previous_league_id uuid references league (id),
    previous_league_rank integer,
    is_playing_next_season bool default true,
    created_at timestamptz not null default current_timestamp
);

-- migrate:down
drop table if exists player;
drop type if exists handedness;