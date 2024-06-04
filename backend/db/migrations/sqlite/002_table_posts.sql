-- +migrate Up

CREATE TABLE posts (
    "post_id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "user_id" INTEGER NOT NULL,
    "post_datetime" DATETIME NOT NULL,
    "title" TEXT NOT NULL,
    "body" TEXT NOT NULL,
    "category_id" INTEGER NOT NULL,
    "group_id" INTEGER NOT NULL,
    FOREIGN KEY ("group_id") REFERENCES "group_index"("group_id")			
    FOREIGN KEY ("user_id") REFERENCES "users"("user_id")
    FOREIGN KEY ("category_id") REFERENCES "categories"("category_id")
);

-- +migrate Down

DROP TABLE posts;