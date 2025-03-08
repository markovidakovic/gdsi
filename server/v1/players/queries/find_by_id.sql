select
    player.id as player_id,
    player.height as player_height,
    player.weight as player_weight,
    player.handedness as player_handedness,
    player.racket as player_racket,
    player.matches_expected as player_matches_expected,
    player.matches_played as player_matches_played,
    player.matches_won as player_matches_won,
    player.matches_scheduled as player_matches_scheduled,
    player.seasons_played as player_seasons_played,
    player.elo as player_elo,
    player.highest_elo as player_highest_elo,
    player.is_elo_provisional as player_is_elo_provisional,
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
    player.created_at as player_created_at
from player
join account on player.account_id = account.id
left join league current_league on player.current_leauge = current_league.id
left join league previous_league on player.perevious_league = previous_league.id
where player.id = $1