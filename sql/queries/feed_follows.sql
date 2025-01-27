-- name: CreateFeedFollow :one
WITH new_follow AS (
INSERT
INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES ($1, $2, $3, $4, $5)
    RETURNING *
    )
SELECT new_follow.*,
       feeds.name as feed_name,
       users.name as user_name
FROM new_follow
         JOIN feeds ON feeds.id = new_follow.feed_id
         JOIN users ON users.id = new_follow.user_id;

-- name: GetFeedFollowsForUser :many
SELECT feed_follows.*,
       feeds.name AS feed_name
FROM feed_follows
         JOIN feeds ON feeds.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1;


-- name: UnfollowFeed :exec
DELETE FROM feed_follows
WHERE feed_id = $1 AND user_id = $2;