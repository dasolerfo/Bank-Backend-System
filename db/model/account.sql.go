// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: account.sql

package db

import (
	"context"
)

const addAccountBalance = `-- name: AddAccountBalance :one
UPDATE accounts
SET money = money + $1
WHERE id = $2
RETURNING id, owner_id, currency, created_at, money, country_code
`

type AddAccountBalanceParams struct {
	Amount int64 `json:"amount"`
	ID     int64 `json:"id"`
}

func (q *Queries) AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, addAccountBalance, arg.Amount, arg.ID)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.OwnerID,
		&i.Currency,
		&i.CreatedAt,
		&i.Money,
		&i.CountryCode,
	)
	return i, err
}

const createAccount = `-- name: CreateAccount :one
INSERT INTO accounts (
  "owner_id",
  "currency",
  "money",
  "country_code"
) VALUES (
    $1, $2, $3, $4
) RETURNING id, owner_id, currency, created_at, money, country_code
`

type CreateAccountParams struct {
	OwnerID     int64    `json:"owner_id"`
	Currency    Currency `json:"currency"`
	Money       int64    `json:"money"`
	CountryCode int32    `json:"country_code"`
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, createAccount,
		arg.OwnerID,
		arg.Currency,
		arg.Money,
		arg.CountryCode,
	)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.OwnerID,
		&i.Currency,
		&i.CreatedAt,
		&i.Money,
		&i.CountryCode,
	)
	return i, err
}

const deleteAccount = `-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1
`

func (q *Queries) DeleteAccount(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteAccount, id)
	return err
}

const getAccount = `-- name: GetAccount :one
SELECT id, owner_id, currency, created_at, money, country_code FROM accounts
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetAccount(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccount, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.OwnerID,
		&i.Currency,
		&i.CreatedAt,
		&i.Money,
		&i.CountryCode,
	)
	return i, err
}

const getAccountForUpdate = `-- name: GetAccountForUpdate :one
SELECT id, owner_id, currency, created_at, money, country_code FROM accounts
WHERE id = $1 LIMIT 1 FOR NO KEY UPDATE
`

func (q *Queries) GetAccountForUpdate(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccountForUpdate, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.OwnerID,
		&i.Currency,
		&i.CreatedAt,
		&i.Money,
		&i.CountryCode,
	)
	return i, err
}

const listAccount = `-- name: ListAccount :many
SELECT id, owner_id, currency, created_at, money, country_code FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2
`

type ListAccountParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListAccount(ctx context.Context, arg ListAccountParams) ([]Account, error) {
	rows, err := q.db.QueryContext(ctx, listAccount, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Account{}
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.OwnerID,
			&i.Currency,
			&i.CreatedAt,
			&i.Money,
			&i.CountryCode,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAccount = `-- name: UpdateAccount :one
UPDATE accounts
SET money = $2
WHERE id = $1
RETURNING id, owner_id, currency, created_at, money, country_code
`

type UpdateAccountParams struct {
	ID    int64 `json:"id"`
	Money int64 `json:"money"`
}

func (q *Queries) UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, updateAccount, arg.ID, arg.Money)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.OwnerID,
		&i.Currency,
		&i.CreatedAt,
		&i.Money,
		&i.CountryCode,
	)
	return i, err
}
