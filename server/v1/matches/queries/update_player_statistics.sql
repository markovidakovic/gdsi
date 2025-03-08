update player
set 
    matches_played = matches_played + 1,
    matches_won = matches_won + case when id = $1 then 1 else 0 end
where id in ($2, $3)