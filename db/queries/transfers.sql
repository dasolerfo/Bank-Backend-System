-- name: CreateTransfer :one
INSERT INTO transfers (
  "from_account_id",
  "to_account_id",
  "amount"
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: GetTranfer :one
SELECT * FROM transfers
WHERE id = $1 LIMIT 1;

-- name: ListTranfers :many
SELECT * FROM transfers
ORDER BY created_at
LIMIT $1
OFFSET $2;