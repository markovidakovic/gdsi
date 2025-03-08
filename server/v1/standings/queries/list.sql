
select
    standing.id as standing_id,
    standing.points as standing_points,
    standing.matches_played as standing_matches_played,
    standing.matches_won as standing_matches_won,
    standing.sets_won as standing_sets_won,
    standing.sets_lost as standing_sets_lost,
    standing.games_won as standing_games_won,
    standing.games_lost as standing_games_lost,
    season.id as standing_season_id,
    season.title as standing_season_title,
    league.id as standing_league_id,
    league.title as standing_league_title,
    player.id as standing_player_id,
    account.name as standing_player_name,
    standing.created_at as standing_created_at
from standing
join season on standing.season_id = season.id
join league on standing.league_id = league.id
join player on standing.player_id = player.id
join account on player.account_id = account.id
where standing.season_id = $1 and standing.league_id = $2
order by
    standing.points desc,
    standing.matches_won desc,
    standing.sets_won desc,
    (standing.sets_won - standing.sets_lost) desc,
    standing.games_won desc,
    (standing.games_won - standing.games_lost) desc,
    standing.created_at desc