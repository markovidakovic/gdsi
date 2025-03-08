with inserted_match as (
    insert into match (court_id, scheduled_at, player_one_id, player_two_id, winner_id, score, season_id, league_id, creator_id)
    values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
    returning id, court_id, scheduled_at, player_one_id, player_two_id, winner_id, score, season_id, league_id, created_at
)
select
    im.id,
    court.id as court_id,
    court.name as court_name,
    im.scheduled_at,
    player1.id as player_one_id,
    account1.name as player_one_name,
    player2.id as player_two_id,
    account2.name as player_two_name,
    winner.id as winner_id,
    account3.name as winner_name,
    im.score,
    season.id as season_id,
    season.title as season_title,
    league.id as league_id,
    league.title as league_title,
    im.created_at
from inserted_match im
join court on im.court_id = court.id
join player player1 on im.player_one_id = player1.id
join account account1 on player1.account_id = account1.id
join player player2 on im.player_two_id = player2.id
join account account2 on player2.account_id = account2.id
left join player winner on im.winner_id = winner.id
left join account account3 on winner.account_id = account3.id
join season on im.season_id = season.id
join league on im.league_id = league.id