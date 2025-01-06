-- name: SessionCreate :exec
INSERT INTO sessions (
    key,
    user_id, 
    created_at
    ) 
    VALUES (?, ?, ?);

-- name: SessionList :many
SELECT id, key, user_id, created_at FROM sessions
ORDER BY id ASC;

-- name: SessionDelete :exec
DELETE FROM sessions WHERE id = ?;