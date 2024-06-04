-- +migrate Up

CREATE TABLE followers (
    "user_id" INTEGER NOT NULL,
    "follower_id" INTEGER NOT NULL,
    FOREIGN KEY ("user_id") REFERENCES "users"("user_id")
    UNIQUE ("user_id", "follower_id")
);

-- +migrate Down
DROP TABLE followers;
