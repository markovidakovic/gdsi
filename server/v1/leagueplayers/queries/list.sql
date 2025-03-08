select
    player.id,
    player.height,
    player.weight,
    player.handedness,
    player.racket,
    player.matches_expected,
    player.matches_played,
    player.matches_won,
    player.matches_scheduled,
    player.seasons_played,
    account.id as account_id,
    account.name as account_name,
    league.id as current_league_id,
    league.title as current_league_title,
    player.created_at
from player
join account on player.account_id = account.id
left join league on player.current_league_id = league.id
where player.current_league_id = $1