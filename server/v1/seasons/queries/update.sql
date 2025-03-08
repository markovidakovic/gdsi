with updated_season as (
    update season 
    set title = $1, description = $2, start_date = $3, end_date = $4
    where id = $5
    returning id, title, description, start_date, end_date, creator_id, created_at
)
select 
    us.id as season_id,
    us.title as season_title,
    us.description as season_description,
    us.start_date as season_start_date,
    us.end_date as season_end_date,
    account.id as creator_id,
    account.name as creator_name,
    us.created_at as season_created_at
from updated_season us
join account on us.creator_id = account.id