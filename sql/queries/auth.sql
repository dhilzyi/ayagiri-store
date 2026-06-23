-- name: CreateUser :one
INSERT INTO
  users (
    created_at,
    updated_at,
    username,
    password_hash,
    role
  )
VALUES
  (NOW(), NOW(), $1, $2, $3)
RETURNING
  *;

-- name: GetUserByUsername :one
SELECT
  *
FROM
  users
WHERE
  username = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE
  id = $1;

-- name: GetSession :one
SELECT
  *
FROM
  sessions
WHERE
  token = $1;

-- name: CreateSession :one
INSERT INTO
  sessions (created_at, expires_at, user_id, token)
VALUES
  (NOW(), $1, $2, $3)
RETURNING
  *;

-- name: DeleteSession :exec
DELETE FROM sessions
WHERE
  token = $1;

-- name: UpdateExpireSession :exec
UPDATE sessions
SET
  expires_at = $1
WHERE
  token = $2;

-- name: CountUsers :one
SELECT
  COUNT(id)
FROM
  users;
