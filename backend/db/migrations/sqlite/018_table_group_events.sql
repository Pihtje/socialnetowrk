-- +migrate Up

CREATE TABLE group_events (
    event_id INTEGER PRIMARY KEY AUTOINCREMENT,
    group_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    day_time DATETIME NOT NULL,
    FOREIGN KEY (group_id) REFERENCES group_index(group_id)
);
-- +migrate Down

DROP TABLE "group_events";