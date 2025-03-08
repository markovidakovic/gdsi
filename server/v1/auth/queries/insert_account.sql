insert into account (name, email, dob, gender, phone_number, password)
values ($1, $2, $3, $4, $5, $6)
returning id, name, email, dob, gender, phone_number, password, role, NULL as player_id, created_at