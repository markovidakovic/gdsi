with updated_match as (
    update match 
    set score = $1, winner_id = $2
    where id = $3 and season_id = $4 and league_id = $5
    returning id, court_id, scheduled_at, player_one_id, player_two_id, winner_id, score, season_id, league_id, created_at
)
select
    um.id,
    court.id as court_id,
    court.name as court_name,
    um.scheduled_at,
    player1.id as player_one_id,
    account1.name as player_one_name,
    player2.id as player_two_id,
    account2.name as player_two_name,
    winner.id as winner_id,
    account3.name as winner_name,
    um.score,
    season.id as season_id,
    season.title as season_title,
    league.id as league_id,
    league.title as league_title,
    um.created_at
from updated_match um
join court on um.court_id = court.id
join player player1 on um.player_one_id = player1.id
join account account1 on player1.account_id = account1.id
join player player2 on um.player_two_id = player2.id
join account account2 on player2.account_id = account2.id
left join player winner on um.winner_id = winner.id
left join account account3 on winner.account_id = account3.id
join season on um.season_id = season.id
join league on um.league_id = league.id