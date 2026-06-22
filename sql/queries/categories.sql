-- name: CreateCategory :one
INSERT INTO
  categories (created_at, updated_at, name, english_name)
VALUES
  (NOW(), NOW(), $1, $2)
RETURNING
  *;

-- name: BulkCreateCategories :copyfrom
INSERT INTO
  categories (created_at, updated_at, name, english_name)
VALUES
  ($1, $2, $3, $4);

-- name: UpdateCategoryByID :exec
UPDATE categories
SET
  updated_at = NOW(),
  name = $1,
  english_name = $2
WHERE
  id = $3;

-- name: DeleteCategoryByID :one
DELETE FROM categories
WHERE
  id = $1
RETURNING
  *;

-- name: ResetCategoryRows :exec
DELETE FROM categories;

-- name: GetCategories :many
SELECT
  *
FROM
  categories
ORDER BY
  id ASC;

-- name: DeleteCategoriesByID :exec
DELETE FROM categories
WHERE
  id = ANY ($1::int[]);
