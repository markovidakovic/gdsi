-- migrate:up
create table season(
    id uuid primary key not null default  uuid_generate_v4(),
    title varchar(250) not null,
    description varchar(500),
    creator_id uuid not null references account (id),
    created_at timestamptz not null default current_timestamp
);

-- migrate:down
drop table if exists season;