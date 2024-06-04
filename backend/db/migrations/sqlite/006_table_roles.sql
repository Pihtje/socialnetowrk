-- +migrate Up

CREATE TABLE "roles" (
    "role_id" INTEGER NOT NULL PRIMARY KEY,
    "role_name" TEXT NOT NULL
);

-- +migrate Down

DROP TABLE roles;