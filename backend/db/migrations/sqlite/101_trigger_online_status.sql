-- +migrate Up
-- +migrate StatementBegin
CREATE TRIGGER "online_status"
AFTER INSERT ON sessions
FOR EACH ROW
BEGIN
    UPDATE users
    SET online_status = "online"
    WHERE user_id = NEW.user_id;
END;

CREATE TRIGGER "user_offline"
BEFORE DELETE ON sessions
FOR EACH ROW
BEGIN
    UPDATE users
    SET online_status = "offline"
    WHERE user_id = OLD.user_id;
END;
-- +migrate StatementEnd

-- +migrate Down
DROP TRIGGER "online_status";
DROP TRIGGER "user_offline";
-- DROP TRIGGER "update_allowed_users";