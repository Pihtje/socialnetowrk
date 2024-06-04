-- +migrate Up
---- everything here is applied when migrate.Up is called

INSERT INTO comments (post_id, user_id, comment_datetime, body) VALUES (1, 1, CURRENT_TIMESTAMP, "Likeeee fr yknoww xdd");
INSERT INTO comments (post_id, user_id, comment_datetime, body) VALUES (4, 1, CURRENT_TIMESTAMP, "Bad post");
INSERT INTO comments (post_id, user_id, comment_datetime, body) VALUES (4, 2, CURRENT_TIMESTAMP, "remove this");
INSERT INTO comments (post_id, user_id, comment_datetime, body) VALUES (5, 3, CURRENT_TIMESTAMP, "Yessss please remove these posts!");
INSERT INTO comments (post_id, user_id, comment_datetime, body) VALUES (1, 3, CURRENT_TIMESTAMP, "Why are you so meannnn T.T");
INSERT INTO comments (post_id, user_id, comment_datetime, body) VALUES (3, 3, CURRENT_TIMESTAMP, "Chill people");

-- +migrate Down
---- everything here is applied when migrate.Down is called

DELETE FROM comments WHERE (post_id, user_id, comment_datetime, body) = (1, 1, CURRENT_TIMESTAMP, "Likeeee fr yknoww xdd");
DELETE FROM comments WHERE (post_id, user_id, comment_datetime, body) = (4, 1, CURRENT_TIMESTAMP, "Bad post");
DELETE FROM comments WHERE (post_id, user_id, comment_datetime, body) = (4, 2, CURRENT_TIMESTAMP, "remove this");
DELETE FROM comments WHERE (post_id, user_id, comment_datetime, body) = (5, 3, CURRENT_TIMESTAMP, "Yessss please remove these posts!");
DELETE FROM comments WHERE (post_id, user_id, comment_datetime, body) = (1, 3, CURRENT_TIMESTAMP, "Why are you so meannnn T.T");
DELETE FROM comments WHERE (post_id, user_id, comment_datetime, body) = (3, 3, CURRENT_TIMESTAMP, "Chill people");