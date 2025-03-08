select
    refresh_token.id,
    account.id as account_id,
    account.role as account_role,
    refresh_token.token_hash,
    refresh_token.device_id,
    refresh_token.ip_address,
    refresh_token.user_agent,
    refresh_token.issued_at,
    refresh_token.expires_at,
    refresh_token.last_used_at,
    refresh_token.is_revoked,
    player.id as player_id
from refresh_token
join account on refresh_token.account_id = account.id
join player on account.id = player.account_id
where refresh_token.token_hash = $1