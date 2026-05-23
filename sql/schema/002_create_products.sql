-- +goose Up
CREATE TABLE products (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name TEXT NOT NULL,
  english_name TEXT NOT NULL,
  price INT NOT NULL,
  category_id INT NOT NULL REFERENCES categories (id),
  discount INT NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE products;
