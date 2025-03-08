with updated_court as (
    update court
    set name = $1
    where id = $2
    returning id, name, creator_id, created_at
)
select
    uc.id as court_id,
    uc.name as court_name,
    account.id as creator_id,
    account.name as creator_name,
    uc.created_at as court_created_at
from updated_court uc
join account on uc.creator_id = account.id