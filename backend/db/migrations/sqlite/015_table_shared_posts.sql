-- +migrate Up

CREATE TABLE shared_posts (
    "post_id" INTEGER NOT NULL,
    "user_id" INTEGER,
    FOREIGN KEY ("post_id") REFERENCES posts("post_id")
    FOREIGN KEY ("user_id") REFERENCES users("user_id")
    UNIQUE ("post_id", "user_id")
);

-- +migrate Down
DROP TABLE shared_posts;