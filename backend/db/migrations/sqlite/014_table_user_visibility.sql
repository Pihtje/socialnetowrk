-- +migrate Up

CREATE TABLE "user_visibility" (
    "user_id" INTEGER NOT NULL,
    "visibility" TEXT NOT NULL,
    FOREIGN KEY ("user_id") REFERENCES users("user_id")
);

-- +migrate Down

DROP TABLE user_visibility;