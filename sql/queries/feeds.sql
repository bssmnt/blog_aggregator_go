-- name: CreateFeed :exec
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (@id,
        NOW(),
        NOW(),
        @name,
        @url,
        @user_id);