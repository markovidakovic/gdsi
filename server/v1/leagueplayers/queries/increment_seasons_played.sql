with updated_player as (
    update player
    set
        seasons_played = seasons_played + 1
    where id = $1 and current_league_id = $2
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
    account.id as player_account_id,
    account.name as player_account_name,
    league.id as player_current_league_id,
    league.title as player_current_league_title,
    up.created_at as player_created_at
from updated_player up
join account on up.account_id = account.id
left join league on up.current_league_id = league.id