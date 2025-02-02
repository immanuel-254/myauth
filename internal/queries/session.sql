-- name: SessionCreate :one
INSERT INTO sessions (
    key,
    user_id, 
    created_at
    ) 
    VALUES (?, ?, ?)
    RETURNING id, key, user_id, created_at;

-- name: SessionRead :one
SELECT id, key, user_id ,created_at FROM sessions
WHERE key = ?;

-- name: SessionList :many
SELECT id, key, user_id, created_at FROM sessions
ORDER BY id ASC;

-- name: SessionTodayList :many
SELECT id, key, user_id, created_at FROM sessions
WHERE DATE(created_at) = DATE('now');

-- name: SessionYesterdayList :many
SELECT id, key, user_id, created_at FROM sessions
WHERE DATE(created_at) = DATE('now', '-1 day');

-- name: SessionWeeklyList :many
SELECT id, key, user_id, created_at FROM sessions
WHERE DATE(created_at) >= DATE('now', 'weekday 0') AND DATE(created_at) <= DATE('now');

-- name: SessionPreviousWeeklyList :many
SELECT id, key, user_id, created_at FROM sessions
WHERE created_at >= DATE('now', 'weekday 0', '-7 days') AND created_at < DATE('now', 'weekday 0');

-- name: SessionMonthlyList :many
SELECT id, key, user_id, created_at FROM sessions
WHERE strftime('%Y-%m', created_at) = strftime('%Y-%m', 'now');

-- name: SessionPreviousMonthlyList :many
SELECT id, key, user_id, created_at FROM sessions
WHERE strftime('%Y-%m', created_at) = strftime('%Y-%m', 'now', '-1 month');

-- name: SessionDelete :exec
DELETE FROM sessions WHERE key = ?;