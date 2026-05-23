-- +goose Up
CREATE TABLE orders (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  table_id INT NOT NULL,
  order_complete BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE order_items (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  order_id UUID NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
  product_id INT NOT NULL REFERENCES products (id),
  quantity INT NOT NULL DEFAULT 1
);

-- +goose Down
DROP TABLE order_items;

DROP TABLE orders;
