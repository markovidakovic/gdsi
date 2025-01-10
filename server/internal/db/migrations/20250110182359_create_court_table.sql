-- migrate:up
CREATE TABLE court(
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    name varchar(250) UNIQUE NOT NULL,
    creator_id UUID REFERENCES account (id) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- migrate:down
DROP TABLE IF EXISTS court;