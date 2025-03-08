with updated_player as (
    update player 
    set height = $1, weight = $2, handedness = $3, racket = $4
    where id = $5
    returning id, height, weight, handedness, racket, matches_expected, matches_played, matches_won, matches_scheduled, seasons_played, account_id, current_league_id, created_at
)
select 
    up.id as player_id,
    up.height as player_height,
    up.weight as player_weight,
    up.handedness as player_handedness,
    up.racket as player_racket,
    up.matches_expected as player_matches_expected,
    up.matches_played as player_matches_played,
    up.matches_won as player_matches_won,
    up.matches_scheduled as player_matches_scheduled,
    up.seasons_played as player_seasons_played,
    up.elo as player_elo,
    up.highest_elo as player_highest_elo,
    up.is_elo_provisional as player_is_elo_provisional,
    account.id as player_account_id,
    account.name as player_account_name,
    current_league.id as player_current_league_id,
    current_league.tier as player_current_league_tier,
    current_league.name as player_current_league_name,
    previous_league.id as player_previous_league_id,
    previous_league.tier as player_previous_league_tier,
    previous_league.name as player_previous_league_name,
    player.previous_league_rank as player_previous_league_rank,
    player.is_playing_next_season as player_is_playing_next_season,
    up.created_at as player_created_at
from updated_player up
join account on up.account_id = account.id
left join league current_league on up.current_league_id = current_league.id
left join league previous_league on up.previous_league_id = previous_league.id