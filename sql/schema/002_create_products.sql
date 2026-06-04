-- +goose Up
CREATE TABLE products (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name TEXT NOT NULL,
  english_name TEXT,
  price INT NOT NULL,
  category_id INT NOT NULL DEFAULT 1 REFERENCES categories (id) ON DELETE SET DEFAULT,
  discount INT DEFAULT 0
);

-- +goose Down
DROP TABLE products;
