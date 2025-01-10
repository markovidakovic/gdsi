-- migrate:up
CREATE TABLE refresh_token(
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    account_id UUID REFERENCES account (id) ON DELETE CASCADE NOT NULL,
    token_hash text NOT NULL,
    device_id text,
    ip_address varchar(250),
    user_agent text,
    issued_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at timestamptz NOT NULL,
    last_used_at timestamptz,
    is_revoked boolean NOT NULL DEFAULT false
    -- CONSTRAINT unique_token_per_device UNIQUE (account_id, device_id)
);

-- migrate:down
DROP TABLE IF EXISTS refresh_token;
