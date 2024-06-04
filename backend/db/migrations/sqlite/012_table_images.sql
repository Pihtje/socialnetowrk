-- +migrate Up

CREATE TABLE images (
    "image_id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "user_id" INTEGER,
    "post_id" INTEGER,
    "comment_id" INTEGER,
    "image_URL" TEXT NOT NULL UNIQUE,
    FOREIGN KEY ("user_id") REFERENCES users("user_id")
    FOREIGN KEY ("post_id") REFERENCES posts("post_id")
    FOREIGN KEY ("comment_id") REFERENCES comments("comment_id")
);

-- +migrate Down

DROP TABLE images;