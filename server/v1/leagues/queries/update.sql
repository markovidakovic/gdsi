
with updated_league as (
    update league
    set title = $1, description = $2
    where id = $3 and season_id = $4
    returning id, title, description, season_id, creator_id, created_at
)
select
    ul.id,
    ul.title,
    ul.description,
    season.id as season_id,
    season.title as season_title,
    account.id as creator_id,
    account.name as creator_name,
    ul.created_at
from updated_league ul
join season on ul.season_id = season.id
join account on ul.creator_id = account.id