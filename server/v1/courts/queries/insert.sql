with inserted_court as (
    insert into court (name, creator_id)
    values ($1, $2)
    returning id, name, creator_id, created_at			
)
select 
    ic.id as court_id, 
    ic.name as court_name, 
    account.id as creator_id, 
    account.name as creator_name, 
    ic.created_at as court_created_at
from inserted_court ic
join account on ic.creator_id = account.id