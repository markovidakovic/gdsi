select
    season.id,
    season.title,
    season.description,
    season.start_date,
    season.end_date,
    account.id as creator_id,
    account.name as creator_name,
    season.created_at
from season
join account on season.creator_id = account.id
where season.id = $1