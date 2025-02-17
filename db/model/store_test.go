package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println(">> before trx: ", account1.Money, account2.Money)

	//run x concurrent transfer transictions
	x := 5
	amount := int64(10)

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < x; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})

			errs <- err
			results <- result
		}()

	}

	existed := make(map[int]bool)

	for i := 0; i < x; i++ {
		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotEmpty(t, result)

		//Check transfer
		transfer := result.Transfer
		require.NotEmpty(t, transfer)
		require.Equal(t, account1.ID, transfer.FromAccountID)
		require.Equal(t, account2.ID, transfer.ToAccountID)
		require.Equal(t, amount, transfer.Amount)
		require.NotZero(t, transfer.ID)
		require.NotZero(t, transfer.CreatedAt)

		_, err = store.GetTranfer(context.Background(), transfer.ID)
		require.NoError(t, err)

		//Check from Entry
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -amount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		//Check to Entry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, amount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		//Check accounts balance

		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		fmt.Println(">>  trx: ", fromAccount.Money, toAccount.Money)

		diff1 := account1.Money - fromAccount.Money
		diff2 := toAccount.Money - account2.Money

		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff2 > 0)
		require.True(t, int(diff1)%int(amount) == 0)

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= x)
		require.NotContains(t, existed, k)
		existed[k] = true

	}

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	fmt.Println(">> diners", account1.Money, " x:", x, " amount:", amount)
	require.Equal(t, account1.Money-(int64(x)*amount), updatedAccount1.Money)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.Equal(t, account2.Money+(int64(x)*amount), updatedAccount2.Money)

	fmt.Println(">> after trx: ", updatedAccount1.Money, updatedAccount2.Money)
}

func TestTransferTxDeadlocks(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println(">> before trx: ", account1.Money, account2.Money)

	//run x concurrent transfer transictions
	x := 10
	amount := int64(10)

	errs := make(chan error)
	//results := make(chan TransferTxResult)

	for i := 0; i < x; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			var accountID1, accountID2 int64
			if i%2 == 0 {
				accountID1 = account1.ID
				accountID2 = account2.ID
			} else {
				accountID1 = account2.ID
				accountID2 = account1.ID
			}
			_, err := store.TransferTx(ctx, TransferTxParams{
				FromAccountID: accountID1,
				ToAccountID:   accountID2,
				Amount:        amount,
			})

			errs <- err
		}()

	}

	//existed := make(map[int]bool)

	for i := 0; i < x; i++ {
		err := <-errs
		require.NoError(t, err)

	}

	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	fmt.Println(">> diners", account1.Money, " x:", x, " amount:", amount)
	require.Equal(t, account1.Money, updatedAccount1.Money)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.Equal(t, account2.Money, updatedAccount2.Money)

	fmt.Println(">> after trx: ", updatedAccount1.Money, updatedAccount2.Money)
}
