-- name: PersistUser :execresult
INSERT INTO integrating_tests_users (name) VALUES (?);

-- name: UpdateUserName :execresult
UPDATE integrating_tests_users
SET name=?
where id = ?;

-- name: GetUsers :many
SELECT *
from integrating_tests_users
WHERE deleted = 0;

-- name: GetUser :one
SELECT *
from integrating_tests_users
WHERE name = ?
  and deleted = 0
order by id desc;
