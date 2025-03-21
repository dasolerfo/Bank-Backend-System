package db

import (
	"context"
	"database/sql"
	"simplebank/factory"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {

	args := CreateAccountParams{
		OwnerID:     factory.RandomOwner(),
		Currency:    Currency(factory.RandomCurreny()),
		Money:       factory.RandomMoney(),
		CountryCode: int32(factory.RandomInt(1, 99)),
	}

	account, err := testQueries.CreateAccount(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, args.OwnerID, account.OwnerID)
	require.Equal(t, args.Currency, account.Currency)
	require.Equal(t, args.Money, account.Money)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	accountS := createRandomAccount(t)
	accountF, err := testQueries.GetAccount(context.Background(), accountS.ID)

	require.NoError(t, err)
	require.NotEmpty(t, accountF)

	require.Equal(t, accountS.ID, accountF.ID)
	require.Equal(t, accountS.OwnerID, accountF.OwnerID)
	require.Equal(t, accountS.Money, accountF.Money)
	require.Equal(t, accountS.Currency, accountF.Currency)
	require.WithinDuration(t, accountS.CreatedAt, accountF.CreatedAt, time.Second)

}

func TestUpdateAccount(t *testing.T) {
	accountS := createRandomAccount(t)

	args := UpdateAccountParams{
		ID:    accountS.ID,
		Money: factory.RandomMoney(),
	}

	accountF, err := testQueries.UpdateAccount(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, accountF)

	require.Equal(t, accountS.ID, accountF.ID)
	require.Equal(t, accountS.OwnerID, accountF.OwnerID)
	require.Equal(t, args.Money, accountF.Money)
	require.Equal(t, accountS.Currency, accountF.Currency)
	require.WithinDuration(t, accountS.CreatedAt, accountF.CreatedAt, time.Second)

}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	account, err := testQueries.GetAccount(context.Background(), account1.ID)

	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account)

}

func TestListAccount(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccount(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
