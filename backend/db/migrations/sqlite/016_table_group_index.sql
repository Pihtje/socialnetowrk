-- +migrate Up

CREATE TABLE group_index (
    "group_id" INTEGER PRIMARY KEY AUTOINCREMENT,
    "user_id" INTEGER NOT NULL,
    "group_title" TEXT NOT NULL,
    "group_description" TEXT NOT NULL,
    FOREIGN KEY ("user_id") REFERENCES users("user_id")
);

-- +migrate Down

DROP TABLE group_index;