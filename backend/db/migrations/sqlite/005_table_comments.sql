-- +migrate Up

CREATE TABLE "comments" (
    "comment_id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "post_id" INTEGER NOT NULL,
    "user_id" INTEGER NOT NULL,
    "comment_datetime" DATETIME NOT NULL,
    "body" TEXT NOT NULL,				
    FOREIGN KEY ("user_id") REFERENCES "users"("user_id")
    FOREIGN KEY ("post_id") REFERENCES "posts"("post_id")
);

-- +migrate Down

DROP TABLE comments;