-- name: CreateProduct :one
INSERT INTO
  products (
    created_at,
    updated_at,
    name,
    english_name,
    price,
    category_id,
    discount
  )
VALUES
  (NOW(), NOW(), $1, $2, $3, $4, $5)
RETURNING
  *;

-- name: BulkCreateProducts :copyfrom
INSERT INTO
  products (
    created_at,
    updated_at,
    name,
    english_name,
    price,
    category_id,
    discount
  )
VALUES
  ($1, $2, $3, $4, $5, $6, $7);

-- name: UpdateDiscount :exec
UPDATE products
SET
  updated_at = NOW(),
  discount = $1
WHERE
  id = $2;

-- name: DeleteProductByID :exec
DELETE FROM products
WHERE
  id = $1;

-- name: GetProducts :many
SElECT
  *
FROM
  products;

-- name: GetProductByID :one
SELECT
  *
FROM
  products
WHERE
  id = $1;

-- name: ResetProductsRows :exec
DELETE FROM products;

-- name: UpdateProductByID :exec
UPDATE products
SET
  updated_at = NOW(),
  name = $1,
  price = $2,
  category_id = $3,
  discount = $4
WHERE
  id = $5;

-- name: GetProductByCategoryID :many
SELECT
  *
FROM
  products
WHERE
  category_id = $1;
