select 
    league.id,
    league.title,
    league.description,
    season.id as season_id,
    season.title as season_title,
    account.id as creator_id,
    account.name as creator_name,
    league.created_at
from league
join season on league.season_id = season.id
join account on league.creator_id = account.id
where league.season_id = $1