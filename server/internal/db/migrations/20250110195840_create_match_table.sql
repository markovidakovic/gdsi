-- migrate:up
CREATE TABLE match(
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    court_id UUID REFERENCES court (id) NOT NULL,
    scheduled_at timestamptz NOT NULL,
    player_one_id UUID REFERENCES player (id) NOT NULL,
    player_two_id UUID REFERENCES player (id) NOT NULL,
    winner_id UUID REFERENCES player (id),
    score text,
    season_id UUID REFERENCES season (id) ON DELETE CASCADE NOT NULL,
    league_id UUID REFERENCES league (id) ON DELETE CASCADE NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- migrate:down
DROP TABLE IF EXISTS match;