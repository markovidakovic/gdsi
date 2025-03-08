select 
    court.id as court_id,
    court.name as court_name,
    account.id as creator_id,
    account.name as creator_name,
    court.created_at as court_created_at
from court
join account on court.creator_id = account.id
where court.id = $1