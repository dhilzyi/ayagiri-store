-- +goose Up
CREATE TABLE category (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name TEXT NOT NULL,
  romaji_name TEXT NOT NULL
);

-- +goose Down
DROP TABLE category;
