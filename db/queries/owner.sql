-- name: CreateOwner :one
INSERT INTO owners (
  "first_name",
  "first_surname",
  "second_surname",
  "born_at",
  "nationality",
  "hashed_password",
  "email"
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: GetOwner :one
SELECT * FROM owners
WHERE id = $1 LIMIT 1;

-- name: GetOwnerByEmail :one
SELECT * FROM owners
WHERE email = $1 LIMIT 1;