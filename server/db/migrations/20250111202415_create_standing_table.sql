-- migrate:up
create table standing(
    id uuid primary key not null default  uuid_generate_v4(),
    points integer not null default 0,
    matches_played integer not null default 0,
    matches_won integer not null default 0,
    sets_won integer not null default 0,
    sets_lost integer not null default 0,
    games_won integer not null default 0,
    games_lost integer not null default 0,
    season_id uuid not null references season (id) on delete cascade,
    league_id uuid not null references league (id) on delete cascade,
    player_id uuid not null references player (id) on delete cascade,
    created_at timestamptz not null default current_timestamp,
    unique (season_id, league_id, player_id)
);

-- migrate:down
drop table if exists standing;