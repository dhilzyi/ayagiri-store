-- name: CreateProduct :one
INSERT INTO
  product (created_at, update_at, name, price, discount)
VALUEs
  (NOW(), NOW(), $1, $2, $3)
RETURNING
  *;

-- name: UpdateDiscount :exec
UPDATE product
SET
  updated_at = NOW(),
  discount = $1
WHERE
  id = $2;

-- name: DeleteProductByID :exec
DELETE FROM product
WHERE
  id = $1;

-- name: GetProducts :many
SElECT
  *
FROM
  product;

-- name: GetProductByID :one
SELECT
  *
FROM
  product
WHERE
  id = $1;

-- name: ResetProductRows :exec
DELETE FROM product;
