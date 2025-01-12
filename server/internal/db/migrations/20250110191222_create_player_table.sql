-- migrate:up
CREATE TYPE handedness AS ENUM ('right', 'left', 'ambidextrous');

CREATE TABLE player(
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    height numeric,
    weight numeric,
    handedness handedness,
    racket varchar(250),
    matches_expected integer NOT NULL DEFAULT 0,
    matches_played integer NOT NULL DEFAULT 0,
    matches_won integer NOT NULL DEFAULT 0,
    matches_scheduled integer NOT NULL DEFAULT 0,
    seasons_played integer NOT NULL DEFAULT 0,
    winning_ratio double precision NOT NULL DEFAULT 0.0,
    activity_ratio double precision NOT NULL DEFAULT 0.0,
    ranking integer,
    elo integer,
    account_id UUID REFERENCES account (id) ON DELETE CASCADE NOT NULL,
    current_league_id UUID REFERENCES league (id),
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- migrate:down
DROP TABLE IF EXISTS player;
DROP TYPE IF EXISTS handedness;