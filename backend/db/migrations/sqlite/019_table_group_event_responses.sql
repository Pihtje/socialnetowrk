-- +migrate Up
---- everything here is applied when migrate.Up is called

CREATE TABLE group_event_responses (
    response_id INTEGER PRIMARY KEY AUTOINCREMENT,
    event_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    response TEXT NOT NULL CHECK( response IN ('Attending', 'Not attending', 'Pending') ),
    FOREIGN KEY (event_id) REFERENCES group_events(event_id),
    FOREIGN KEY (user_id) REFERENCES users(user_id),
    UNIQUE (event_id, user_id)
);

-- +migrate Down
---- everything here is applied when migrate.Down is called

DROP TABLE group_event_responses;