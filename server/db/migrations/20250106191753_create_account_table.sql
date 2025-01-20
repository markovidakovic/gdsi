-- migrate:up
create extension if not exists "uuid-ossp";
create type gender as enum ('male', 'female');

create table account (
    id uuid primary key not null default uuid_generate_v4(),
    name varchar(250) not null,
    email varchar(250) not null unique,
    dob date not null, 
    gender gender not null,
    phone_number varchar(250) not null,
    password text not null,
    created_at timestamptz not null default current_timestamp
);

-- migrate:down
drop table if exists account;
drop type if exists gender;
