-- +migrate Up
CREATE TABLE "categories" (
    "category_id" INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    "title" TEXT NOT NULL,
    "description" TEXT NOT NULL
);

-- +migrate Down
DROP TABLE "categories";