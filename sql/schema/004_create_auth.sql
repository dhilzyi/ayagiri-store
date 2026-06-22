-- +goose Up
CREATE TABLE users (
  id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  username TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  role TEXT NOT NULL DEFAULT 'staff'
);

CREATE TABLE sessions (
  token TEXT PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  user_id INT REFERENCES users (id) ON DELETE CASCADE
);

-- +goose Down
DROP TABLE sessions;

DROP TABLE users;
