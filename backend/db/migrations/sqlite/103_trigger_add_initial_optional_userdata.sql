-- +migrate Up

-- +migrate StatementBegin
CREATE TRIGGER add_initial_optional_userdata
AFTER INSERT ON users
BEGIN
    INSERT INTO "optional_userdata" ("user_id", "image_URL", "nickname", "about_me") VALUES (NEW.user_id, NULL, NULL, NULL);
END;
-- +migrate StatementEnd


-- +migrate Down
DROP TRIGGER add_initial_optional_userdata;