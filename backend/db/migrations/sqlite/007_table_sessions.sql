-- +migrate Up

CREATE TABLE "sessions" (
    "session_id" TEXT NOT NULL PRIMARY KEY,
    "user_id" INTEGER,
    "email" TEXT NOT NULL,
    "datetime" DATETIME NOT NULL,
    "role_id" INTEGER NOT NULL,
    FOREIGN KEY ("user_id") REFERENCES "users"("user_id")
    FOREIGN KEY ("role_id") REFERENCES "roles"("role_id")
    FOREIGN KEY ("email") REFERENCES "users"("email")
);

-- +migrate Down
DROP TABLE sessions;