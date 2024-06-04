-- +migrate Up

CREATE TABLE group_members (
    "group_id" INTEGER NOT NULL,
    "user_id" INTEGER NOT NULL,
    FOREIGN KEY ("user_id") REFERENCES users("user_id")
    FOREIGN KEY ("group_id") REFERENCES group_index("group_id")
    UNIQUE ("user_id", "group_id")
);

-- +migrate Down

DROP TABLE "group_members";