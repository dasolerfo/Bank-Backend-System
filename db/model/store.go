package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

type StoreSQL struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &StoreSQL{
		db:      db,
		Queries: New(db),
	}
}

func (store *StoreSQL) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)
	err = fn(q)

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
	}

	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

func (store *StoreSQL) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error
		//txName := ctx.Value(txKey)
		//fmt.Println(txName, " create transfer")
		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})

		if err != nil {
			return err
		}

		//fmt.Println(txName, " create entry 1")

		result.FromEntry, err = q.CreateEntries(ctx, CreateEntriesParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		//fmt.Println(txName, " create entry 2")
		result.ToEntry, err = q.CreateEntries(ctx, CreateEntriesParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		//fmt.Println(txName, " get account 1")
		//TODO: Update the accounts balance
		//account1, err := q.GetAccountForUpdate(context.Background(), arg.FromAccountID)
		//if err != nil {
		//	return err
		//}

		//fmt.Println(txName, " update account 1")

		if arg.FromAccountID < arg.ToAccountID {

			result.FromAccount, result.ToAccount, err = transferMoney(ctx, arg.FromAccountID, arg.ToAccountID, -arg.Amount, arg.Amount, q)

			if err != nil {
				return err
			}

		} else {
			result.ToAccount, result.FromAccount, err = transferMoney(ctx, arg.ToAccountID, arg.FromAccountID, arg.Amount, -arg.Amount, q)

			if err != nil {
				return err
			}
		}

		return nil
	})
	return result, err
}

func transferMoney(
	ctx context.Context,
	accountID1 int64,
	accountID2 int64,
	amount1 int64,
	amount2 int64,
	q *Queries) (account1 Account, account2 Account, err error) {

	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})
	if err != nil {
		return
	}
	return
}
