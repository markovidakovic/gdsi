select 
    account.id as account_id,
    account.name as account_name,
    account.email as account_email,
    account.dob as account_dob,
    account.gender as account_gender,
    account.phone_number as account_phone_number,
    account.role as account_role,
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
    league.id as league_id,
    league.title as league_title,
    player.created_at as player_created_at,
    account.created_at as account_created_at
from account
left join player on account.id = player.account_id
left join league on player.current_league_id = league.id
where account.id = $1