-- name: UserCreate :one
INSERT INTO users (
    email, 
    password, 
    created_at, 
    updated_at
    ) 
    VALUES (?, ?, ?, ?)
    RETURNING id, email, created_at, updated_at;

-- name: UserList :many
SELECT id, email, created_at, updated_at FROM users
ORDER BY id ASC;

-- name: UserRead :one
SELECT id, email, created_at, updated_at FROM users
WHERE id = ?;

-- name: UserLoginRead :one
SELECT email, password FROM users
WHERE email = ?;

-- name: UserUpdatePassword :one
UPDATE users SET password = ?, updated_at = ? WHERE id = ? 
RETURNING id, email, created_at, updated_at;

-- name: UserUpdateEmail :one
UPDATE users SET email = ?, updated_at = ? WHERE id = ? 
RETURNING id, email, created_at, updated_at;

-- name: UserUpdateIsActive :one
UPDATE users SET isactive = ?, updated_at = ? WHERE id = ? 
RETURNING id, email, created_at, updated_at;

-- name: UserUpdateIsStaff :one
UPDATE users SET isstaff = ?, updated_at = ? WHERE id = ? 
RETURNING id, email, created_at, updated_at;

-- name: UserUpdateIsAdmin :one
UPDATE users SET isadmin = ?, updated_at = ? WHERE id = ? 
RETURNING id, email, created_at, updated_at;

-- name: UserDelete :exec
DELETE FROM users WHERE id = ?;
