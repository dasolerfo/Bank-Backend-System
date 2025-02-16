-- name: CreateAccount :one
INSERT INTO accounts (
  "owner_id",
  "currency",
  "money",
  "country_code"
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1 FOR NO KEY UPDATE;


-- name: ListAccount :many
SELECT * FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts
SET money = $2
WHERE id = $1
RETURNING *;

-- name: AddAccountBalance :one
UPDATE accounts
SET money = money + sqlc.arg(amount)
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;