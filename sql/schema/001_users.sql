-- +goose Up
CREATE TABLE users(
    id UUID Primary Key,   
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    name TEXT UNIQUE
);

-- +goose Down
Drop Table users;