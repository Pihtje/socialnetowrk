-- +migrate Up
---- everything here is applied when migrate.Up is called

INSERT INTO categories (title, description) VALUES ("Public", "For posts that everyone can see");
INSERT INTO categories (title, description) VALUES ("Private", "For posts only the followers of the creator can see");
INSERT INTO categories (title, description) VALUES ("Shared", "For posts that only select users chosen by the creator can see");

-- +migrate Down
---- everything here is applied when migrate.Down is called

DELETE FROM categories WHERE (title, description) = ("Public", "For posts that everyone can see");
DELETE FROM categories WHERE (title, description) = ("Private", "For posts only the followers of the creator can see");
DELETE FROM categories WHERE (title, description) = ("Shared", "For posts that only select users chosen by the creator can see");