-- migrate:up
CREATE TABLE season(
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    title varchar(250) NOT NULL,
    description varchar(500),
    creator_id UUID REFERENCES account (id) NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- migrate:down
DROP TABLE IF EXISTS season;