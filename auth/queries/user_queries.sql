-- name: GetUserCount :one
    SELECT count(*) FROM user;

-- name: GetUser :one
SELECT * FROM  user
WHERE email = ? LIMIT 1;

-- name: ListUsers :many
SELECT * FROM user
ORDER BY name;

-- name: CreateUsers :one
INSERT INTO user (
    name, email, password, uuid
) VALUES (
  ?,?,?,?
         )
RETURNING *;

-- name: UpdateUser :exec
UPDATE user
set name = ?,
    email = ?,
    password = ?
WHERE pk = ?;

-- name: DeleteUser :exec
DELETE FROM user
WHERE pk = ?;