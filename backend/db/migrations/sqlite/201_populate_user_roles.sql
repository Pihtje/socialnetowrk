-- +migrate Up
---- everything here is applied when migrate.Up is called

INSERT INTO roles (role_id, role_name) VALUES (0, "guest");
INSERT INTO roles (role_id, role_name) VALUES (1, "user");

-- +migrate Down
---- everything here is applied when migrate.Down is called

DELETE FROM roles WHERE (role_id, role_name) = (0, "guest");
DELETE FROM roles WHERE (role_id, role_name) = (1, "user");