-- +migrate Up

-- +migrate StatementBegin
CREATE TRIGGER create_response_for_event_on_new_event
AFTER INSERT ON group_events
BEGIN
    INSERT INTO group_event_responses (event_id, user_id, response)
    SELECT NEW.event_id, user_id, "Pending"
    FROM group_members
    WHERE group_id = NEW.group_id;
END;
-- +migrate StatementEnd

-- +migrate StatementBegin
CREATE TRIGGER create_response_for_event_on_new_member
AFTER INSERT ON group_members
BEGIN
    INSERT INTO group_event_responses (event_id, user_id, response)
    SELECT event_id, NEW.user_id, "Pending"
    FROM group_events
    WHERE group_id = NEW.group_id;
END;
-- +migrate StatementEnd


-- +migrate Down
DROP TRIGGER create_response_for_event_on_new_event;
DROP TRIGGER create_response_for_event_on_new_member;