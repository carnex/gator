-- +goose Up
CREATE TABLE posts(
    id UUID Primary Key,   
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    title TEXT NOT NULL,
    url TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL,
    published_at TIMESTAMP NOT NULL,
    feed_id UUID NOT NULL,
    CONSTRAINT  posts_feeds_feed_id
    FOREIGN KEY (feed_id)
    REFERENCES feeds (id)
);

-- +goose Down
Drop Table posts;