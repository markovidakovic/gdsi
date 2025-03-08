update refresh_token
set is_revoked = true
where id = $1