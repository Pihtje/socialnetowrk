-- +migrate Up

CREATE TABLE "direct_messages" (
    "message_id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "sender_id" INTEGER NOT NULL,
    "target_id" INTEGER NOT NULL,
    "message" TEXT NOT NULL,
    "message_seen" INTEGER NOT NULL,
    "message_datetime" DATETIME NOT NULL,
    "group_id" INTEGER NOT NULL,
    FOREIGN KEY ("sender_id") REFERENCES "users"("user_id")
    FOREIGN KEY ("target_id") REFERENCES "users"("user_id")
    FOREIGN KEY ("group_id") REFERENCES "group_index"("group_id")
);

-- +migrate Down

DROP TABLE direct_messages;