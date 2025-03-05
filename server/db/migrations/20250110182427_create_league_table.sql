-- migrate:up
create table league(
    id uuid primary key not null default  uuid_generate_v4(),
    tier integer not null default 1,
    name varchar(1) not null default 'A',
    max_players integer not null default 6,
    season_id uuid not null references season (id) on delete cascade,
    creator_id uuid not null references account (id),
    created_at timestamptz not null default current_timestamp
);

-- migrate:down
drop table if exists league;