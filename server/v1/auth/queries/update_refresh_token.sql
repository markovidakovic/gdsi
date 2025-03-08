update refresh_token
set last_used_at = $1
where id = $2