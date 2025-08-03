-- name: CreateMessage :one
INSERT INTO messages (chat_id, sender_id, content)
VALUES ($1, $2, $3)
RETURNING *;

-- name: ListMessagesByChat :many
SELECT * FROM messages
WHERE chat_id = $1
ORDER BY sent_at DESC
LIMIT $2 OFFSET $3;
