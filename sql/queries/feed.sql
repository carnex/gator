-- name: CreateFeed :one
INSERT INTO feeds(id, created_at, updated_at, name, url, user_id)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT feeds.name, feeds.url, users.name FROM feeds
INNER JOIN users
ON users.id = feeds.user_id;

-- name: CreateFeedFollow :one
WITH feed_follows AS (
    INSERT INTO feed_follows(id, created_at, updated_at, user_id, feed_id)
    VALUES(
    $1,
    $2,
    $3,
    $4,
    $5 
)
RETURNING *
)
SELECT
    feed_follows.*,
    feeds.name AS feed_name,
    users.name AS user_name
FROM feed_follows
INNER JOIN users ON users.id = feed_follows.user_id
INNER JOIN feeds ON feeds.id = feed_follows.feed_id;


-- name: GetFeed :one
SELECT * from feeds 
WHERE url = $1;

-- name: GetFeedsFollowsForUser :many
SELECT users.name, feeds.name
FROM feed_follows
INNER JOIN users ON users.id = feed_follows.user_id
INNER JOIN feeds ON feeds.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1;

-- name: Unfollow :exec
DELETE FROM feed_follows
WHERE user_id = $1
  AND feed_id = $2;