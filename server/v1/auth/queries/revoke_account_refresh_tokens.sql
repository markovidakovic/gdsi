update refresh_token
set is_revoked = true
where account_id = $1 and is_revoked = false