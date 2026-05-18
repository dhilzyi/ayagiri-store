-- +goose Up
CREATE TABLE product (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name TEXT NOT NULL,
  romaji_name TEXT NOT NULL,
  price INT NOT NULL,
  category_id INT REFERENCES category (id),
  discount INT DEFAULT 0
);

-- +goose Down
DROP TABLE product;
