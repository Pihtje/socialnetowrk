-- +migrate Up
---- everything here is applied when migrate.Up is called

-- +migrate StatementBegin
CREATE TRIGGER remove_deleted_posts_from_shared_posts_table
BEFORE DELETE ON posts
BEGIN
    DELETE FROM shared_posts WHERE post_id = OLD.post_id;
END;

-- +migrate StatementEnd

-- +migrate Down
---- everything here is applied when migrate.Down is called

DROP TRIGGER remove_deleted_posts_from_shared_posts_table;