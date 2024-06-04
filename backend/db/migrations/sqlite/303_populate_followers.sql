-- +migrate Up

INSERT INTO followers (user_id, follower_id) VALUES (1, 2);
INSERT INTO followers (user_id, follower_id) VALUES (1, 3);
INSERT INTO followers (user_id, follower_id) VALUES (1, 5);
INSERT INTO followers (user_id, follower_id) VALUES (2, 4);
INSERT INTO followers (user_id, follower_id) VALUES (2, 1);

-- +migrate Down

DELETE FROM followers WHERE (user_id, follower_id) = (1, 2);
DELETE FROM followers WHERE (user_id, follower_id) = (1, 3);
DELETE FROM followers WHERE (user_id, follower_id) = (1, 5);
DELETE FROM followers WHERE (user_id, follower_id) = (2, 4);
DELETE FROM followers WHERE (user_id, follower_id) = (2, 1);