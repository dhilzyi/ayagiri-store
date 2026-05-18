-- name: CreateCategory :one
INSERT INTO
  category (created_at, updated_at, name, romaji_name)
VALUES
  (NOW(), NOW(), $1, $2)
RETURNING
  *;

-- name: UpdateCategoryByID :exec
UPDATE category
SET
  update_at = NOW(),
  name = $1,
  romaji_name = $2
WHERE
  id = $3;

-- name: DeleteCategoryByID :one
DELETE FROM category
WHERE
  id = $1
RETURNING
  *;

-- name: ResetCategoryRows :exec
DELETE FROM category;

-- name: GetCategories :many
SELECT
  *
FROM
  category;
