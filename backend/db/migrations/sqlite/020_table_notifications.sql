-- +migrate Up

CREATE TABLE notifications (
    notification_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    sender_id INTEGER NOT NULL,
    target_id INTEGER NOT NULL,
    description TEXT NOT NULL,
    seen_by_target INTEGER NOT NULL CHECK (seen_by_target IN (0, 1)) DEFAULT 0,
    seen_by_sender INTEGER NOT NULL CHECK (seen_by_sender IN (0, 1)) DEFAULT 0,
    notification_status TEXT NOT NULL CHECK (notification_status IN ("accepted", "rejected", "pending")) DEFAULT "pending",
    notification_value TEXT,
    FOREIGN KEY ("sender_id") REFERENCES "users"("user_id")
    FOREIGN KEY ("target_id") REFERENCES "users"("user_id")
    --UNIQUE (sender_id, target_id, (description IN ("followRequest", "groupInvite", "groupJoinRequest")))
);

CREATE UNIQUE INDEX block_duplicates
ON notifications (sender_id, target_id, description, notification_value)
WHERE description IN ('followRequest', 'groupInvite', 'groupJoinRequest');

-- +migrate Down
---- everything here is applied when migrate.Down is called

DROP TABLE notifications;

DROP INDEX block_duplicates;