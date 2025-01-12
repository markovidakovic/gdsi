-- migrate:up
CREATE TABLE standing(
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    points integer NOT NULL DEFAULT 0,
    matches_played integer NOT NULL DEFAULT 0,
    matches_won integer NOT NULL DEFAULT 0,
    sets_won integer NOT NULL DEFAULT 0,
    sets_lost integer NOT NULL DEFAULT 0,
    games_won integer NOT NULL DEFAULT 0,
    games_lost integer NOT NULL DEFAULT 0,
    season_id UUID REFERENCES season (id) ON DELETE CASCADE NOT NULL,
    league_id UUID REFERENCES league (id) ON DELETE CASCADE NOT NULL,
    player_id UUID REFERENCES player (id) ON DELETE CASCADE NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- migrate:down
DROP TABLE IF EXISTS standing;