select 
    account.id as account_id, 
    account.name as account_name, 
    account.email as account_email, 
    account.dob as account_dob, 
    account.gender as account_gender, 
    account.phone_number as account_phone_number, 
    account.password as account_password, 
    account.role as account_role, 
    player.id as player_id, 
    account.created_at as account_created_at
from account
left join player on player.account_id = account.id
where account.email = $1