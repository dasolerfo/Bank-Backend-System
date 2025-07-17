-- name: CreateSession :one
INSERT INTO sessions (
    id,
    owner_id,
    refresh_token,
    client_ip,
    user_agent,
    is_blocked,
    expires_at, 
    email
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetSession :one
SELECT * FROM sessions
WHERE id = $1 LIMIT 1;



