-- name: CreatePosts :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: GetPostsForUser :many
SELECT
  p.url,
  p.title,
  p.description,
  p.published_at
FROM posts p
JOIN feeds f
  ON p.feed_id = f.id
WHERE f.id IN (
  SELECT id
  FROM feeds
  WHERE feeds.user_id = $1
  UNION
  SELECT feed_id
  FROM feed_follows
  WHERE feed_follows.user_id = $1
)
ORDER BY
  f.last_fetched_at DESC NULLS LAST,
  p.published_at DESC
LIMIT $2;


