-- migrate:up
create table match(
    id uuid primary key not null default  uuid_generate_v4(),
    court_id uuid not null references court (id),
    scheduled_at timestamptz not null,
    player_one_id uuid not null references player (id),
    player_two_id uuid not null references player (id),
    winner_id uuid references player (id),
    score text,
    season_id uuid not null references season (id) on delete cascade,
    league_id uuid not null references league (id) on delete cascade,
    creator_id uuid not null references player (id) on delete cascade,
    created_at timestamptz not null default current_timestamp
);

-- migrate:down
drop table if exists match;