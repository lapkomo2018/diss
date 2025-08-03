-- name: CreateChat :one
INSERT INTO chats (name)
VALUES ($1)
RETURNING *;

-- name: GetChat :one
SELECT * FROM chats
WHERE id = $1;

-- name: ListChats :many
SELECT * FROM chats
ORDER BY id;
