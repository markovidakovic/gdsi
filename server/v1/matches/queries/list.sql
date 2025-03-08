select
    match.id,
    court.id as court_id,
    court.name as court_name,
    match.scheduled_at,
    player1.id as player_one_id,
    account1.name as player_one_name,
    player2.id as player_two_id,
    account2.name as player_two_name,
    winner.id as winner_id,
    account3.name as winner_name,
    match.score,
    season.id as season_id,
    season.title as season_title,
    league.id as league_id,
    league.title as league_title,
    match.created_at
from match
join court on match.court_id = court.id
join player player1 on match.player_one_id = player1.id
join account account1 on player1.account_id = account1.id
join player player2 on match.player_two_id = player2.id
join account account2 on player2.account_id = account2.id
left join player winner on match.winner_id = winner.id
left join account account3 on winner.account_id = account3.id
join season on match.season_id = season.id
join league on match.league_id = league.id
where match.season_id = $1 and match.league_id = $2