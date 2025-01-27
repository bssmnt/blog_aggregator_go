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
WHERE url = $1 AND feed_id = $2
LIMIT 1;
