-- +goose Up
CREATE TABLE posts
(
    id           UUID PRIMARY KEY,
    title        TEXT UNIQUE NOT NULL,
    url          TEXT UNIQUE NOT NULL,
    feed_id      UUID        NOT NULL,
    FOREIGN KEY (feed_id)
        REFERENCES feeds (id)
        ON DELETE CASCADE,
    published_at TIMESTAMP   NOT NULL,
    created_at   TIMESTAMP   NOT NULL,
    updated_at   TIMESTAMP   NOT NULL
);

-- +goose Down
DROP TABLE posts;
