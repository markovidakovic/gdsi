
insert into standing (points, matches_played, matches_won, sets_won, sets_lost, games_won, games_lost, season_id, league_id, player_id)
values ($1, 1, $2, $3, $4, $5, $6, $7, $8, $9)
on conflict (season_id, league_id, player_id) do update
set
    points = standing.points + $1,
    matches_played = standing.matches_played + 1,
    matches_won = standing.matches_won + $2,
    sets_won = standing.sets_won + $3,
    sets_lost = standing.sets_lost + $4,
    games_won = standing.games_won + $5,
    games_lost = standing.games_lost + $6