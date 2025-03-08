insert into refresh_token (account_id, token_hash, issued_at, expires_at)
values ($1, $2, $3, $4)
returning id, account_id, token_hash, device_id, ip_address, user_agent, issued_at, expires_at, last_used_at, is_revoked