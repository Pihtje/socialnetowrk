-- +migrate Up

-- +migrate StatementBegin
CREATE TRIGGER popoulate_user_visibility_table
AFTER INSERT ON users
BEGIN
    INSERT INTO "user_visibility" ("user_id", "visibility") VALUES (NEW.user_id, "private");
END;


-- +migrate StatementEnd


-- +migrate Down
DROP TRIGGER popoulate_user_visibility_table;