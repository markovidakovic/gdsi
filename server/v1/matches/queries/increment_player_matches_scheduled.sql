
update player
set
    matches_scheduled = matches_scheduled + 1
where id = $1