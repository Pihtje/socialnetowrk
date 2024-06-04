-- +migrate Up
-- +migrate StatementBegin
CREATE TRIGGER remove_event_responses
BEFORE DELETE ON group_events
BEGIN
    DELETE FROM group_event_responses WHERE event_id = OLD.event_id;
    DELETE FROM notifications WHERE description = "event" AND notification_value = OLD.event_id;
END;
-- +migrate StatementEnd

-- +migrate Down
---- everything here is applied when migrate.Down is called
DROP TRIGGER remove_event_responses;