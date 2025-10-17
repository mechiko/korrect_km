-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
-- Create users table
CREATE TABLE users
(
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    username   TEXT NOT NULL UNIQUE,
    email      TEXT NOT NULL UNIQUE,
    password   TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE sessions
(
    token  TEXT PRIMARY KEY,
    data   BLOB NOT NULL,
    expiry REAL NOT NULL
);
-- Index to speed up expiry-based look-ups / sweeps
CREATE INDEX sessions_expiry_idx ON sessions(expiry);
-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

-- Drop indexes
DROP INDEX IF EXISTS sessions_expiry_idx;

-- Drop tables
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS sessions;
