-- name: CreateFeed :exec
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (@id,
        NOW(),
        NOW(),
        @name,
        @url,
        @user_id);

-- name: GetFeeds :many
SELECT feeds.name AS feed_name, feeds.url AS feed_url, users.name AS users_name
FROM feeds
         INNER JOIN users
                    ON feeds.user_id = users.id;

-- name: GetFeedByURL :one
SELECT *
FROM feeds
WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = NOW(),
    updated_at      = NOW()
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
WHERE last_fetched_at IS NULL
   OR last_fetched_at < NOW() - INTERVAL '15 minutes'
ORDER BY last_fetched_at NULLS FIRST LIMIT 1;


