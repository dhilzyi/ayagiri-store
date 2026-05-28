-- name: CreateOrder :one
INSERT INTO
  orders (id, created_at, updated_at, table_id)
VALUES
  ($1, NOW(), NOW(), $2)
RETURNING
  *;

-- name: DeleteOrder :exec
DELETE FROM orders
WHERE
  id = $1;

-- name: GetOrderByID :one
SELECT
  *
FROM
  orders
WHERE
  id = $1;

-- name: BulkCreateOrderItem :copyfrom
INSERT INTO
  order_items (
    created_at,
    updated_at,
    order_id,
    product_id,
    quantity
  )
VALUES
  ($1, $2, $3, $4, $5);

-- name: GetOrders :many
SELECT
  *
FROM
  orders;

-- name: OrderComplete :one
UPDATE orders
SET
  order_complete = $1
WHERE
  id = $2
RETURNING
  *;

-- name: GetOrderItems :many
SELECT
  order_items.*,
  products.name AS product_name,
  orders.order_complete
FROM
  order_items
  INNER JOIN products ON order_items.product_id = products.id
  INNER JOIN orders ON order_items.order_id = orders.id;
