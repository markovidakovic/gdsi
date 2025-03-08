with inserted_season as (
    insert into season (title, description, start_date, end_date, creator_id)
    values ($1, $2, $3, $4, $5)
    returning id, title, description, start_date, end_date, creator_id, created_at
)
select s.id, s.title, s.description, s.start_date, s.end_date, account.id as creator_id, account.name as creator_name, s.created_at
from inserted_season s
join account on s.creator_id = account.id