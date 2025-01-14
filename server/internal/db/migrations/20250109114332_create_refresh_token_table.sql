-- migrate:up
create table refresh_token(
    id uuid primary key not null default  uuid_generate_v4(),
    account_id uuid not null references account (id) on delete cascade,
    token_hash text not null,
    device_id text,
    ip_address varchar(250),
    user_agent text,
    issued_at timestamptz not null default current_timestamp,
    expires_at timestamptz not null,
    last_used_at timestamptz,
    is_revoked boolean not null default false
    -- constraint unique_token_per_device unique (account_id, device_id)
);

-- migrate:down
drop table if exists refresh_token;
