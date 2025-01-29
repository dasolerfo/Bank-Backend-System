-- name: CreateOwner :one
INSERT INTO owners (
  "first_name",
  "first_surname",
  "second_surname",
  "born_at",
  "nationality"
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;