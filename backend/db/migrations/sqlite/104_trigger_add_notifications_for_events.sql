-- +migrate Up

-- +migrate StatementBegin
CREATE TRIGGER create_notifications_for_existing_members
AFTER INSERT ON group_events
BEGIN
    INSERT INTO notifications (sender_id, target_id, description, notification_value, seen_by_sender)
    SELECT "group", user_id, "event", NEW.event_id, "1"
    FROM group_members
    WHERE group_id = NEW.group_id;
END;
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TRIGGER create_notifications_for_new_members
AFTER INSERT ON group_members
BEGIN
    INSERT INTO notifications (sender_id, target_id, description, notification_value, seen_by_sender)
    SELECT "group", NEW.user_id, "event", event_id, "1"
    FROM group_events
    WHERE group_id = NEW.group_id;
END;
-- +migrate StatementEnd


-- +migrate Down
DROP TRIGGER create_notifications_for_existing_members;
DROP TRIGGER create_notifications_for_new_members;