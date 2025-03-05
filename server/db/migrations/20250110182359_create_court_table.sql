-- migrate:up
create table court(
    id uuid primary key not null default uuid_generate_v4(),
    name varchar(250) not null unique,
    created_at timestamptz not null default current_timestamp
);

-- migrate:down
drop table if exists court;