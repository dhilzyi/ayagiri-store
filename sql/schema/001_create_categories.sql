-- +goose Up
CREATE TABLE categories (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name TEXT NOT NULL,
  english_name TEXT NOT NULL
);

-- +goose Down
DROP TABLE categories;
