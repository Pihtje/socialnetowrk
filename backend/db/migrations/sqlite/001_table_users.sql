-- +migrate Up

CREATE TABLE users (
    "user_id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "email" TEXT NOT NULL UNIQUE,
    "password" TEXT NOT NULL,
    "first_name" TEXT NOT NULL,
    "last_name" TEXT NOT NULL,
    "date_of_birth" TEXT NOT NULL,
    "online_status" TEXT NOT NULL,
    "role_id" INTEGER NOT NULL,
    FOREIGN KEY ("role_id") REFERENCES "roles"("role_id")
);

-- +migrate Down
DROP TABLE users;