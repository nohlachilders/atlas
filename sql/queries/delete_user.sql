-- name: DeleteUser :exec
delete from users
where id = $1;
