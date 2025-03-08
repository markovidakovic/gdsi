with inserted_league as (
    insert into league (title, description, season_id, creator_id)
    values ($1, $2, $3, $4)
    returning id, title, description, season_id, creator_id, created_at
)
select
    il.id,
    il.title,
    il.description,
    season.id as season_id,
    season.title as season_title,
    account.id as creator_id,
    account.name as creator_name,
    il.created_at
from inserted_league il
join season on il.season_id = season.id
join account on il.creator_id = account.id