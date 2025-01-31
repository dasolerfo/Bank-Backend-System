// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: entries.sql

package db

import (
	"context"
	"database/sql"
)

const createEntries = `-- name: CreateEntries :one
INSERT INTO entries (
  "account_id",
  "amount"
) VALUES (
    $1, $2
) RETURNING id, account_id, amount, created_at
`

type CreateEntriesParams struct {
	AccountID sql.NullInt64 `json:"account_id"`
	Amount    float64       `json:"amount"`
}

func (q *Queries) CreateEntries(ctx context.Context, arg CreateEntriesParams) (Entry, error) {
	row := q.db.QueryRowContext(ctx, createEntries, arg.AccountID, arg.Amount)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const getEntries = `-- name: GetEntries :one
SELECT id, account_id, amount, created_at FROM entries
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetEntries(ctx context.Context, id int64) (Entry, error) {
	row := q.db.QueryRowContext(ctx, getEntries, id)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const listEntries = `-- name: ListEntries :many
SELECT id, account_id, amount, created_at FROM entries
WHERE account_id = $1
ORDER BY created_at
LIMIT $2
OFFSET $3
`

type ListEntriesParams struct {
	AccountID sql.NullInt64 `json:"account_id"`
	Limit     int32         `json:"limit"`
	Offset    int32         `json:"offset"`
}

func (q *Queries) ListEntries(ctx context.Context, arg ListEntriesParams) ([]Entry, error) {
	rows, err := q.db.QueryContext(ctx, listEntries, arg.AccountID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Entry
	for rows.Next() {
		var i Entry
		if err := rows.Scan(
			&i.ID,
			&i.AccountID,
			&i.Amount,
			&i.CreatedAt,
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
