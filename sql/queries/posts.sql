-- name: InsertPost :exec
INSERT INTO posts(id, title, url, feed_id, published_at, created_at, updated_at)
VALUES (@id,
        @title,
        @url,
        @feed_id,
        @published_at,
        @created_at,
        @updated_at) ON CONFLICT (url) DO
UPDATE SET
    title = EXCLUDED.title,
    published_at = EXCLUDED.published_at,
    updated_at = EXCLUDED.updated_at;

-- name: GetPostByURL :one
SELECT *
FROM posts
WHERE url = $1
  AND feed_id = $2 LIMIT 1;

-- name: GetPostsForUser :many
SELECT posts.title, posts.url, posts.feed_id, posts.published_at
FROM feed_follows
         JOIN posts ON feed_follows.feed_id = posts.feed_id
WHERE user_id = $1
ORDER BY posts.updated_at DESC LIMIT $2;
