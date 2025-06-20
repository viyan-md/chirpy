-- +goose up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    email TEXT
);

-- +goose down
DROP TABLE IF EXISTS users;

