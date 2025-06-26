package db

import (
	"context"
	"time"

	"simplebank/factory"

	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {

	owner := createRandomOwner(t)
	account, err := testQueries.CreateAccount(context.Background(), CreateAccountParams{
		Currency:    Currency(factory.RandomCurreny()),
		OwnerID:     owner.ID,
		Money:       factory.RandomMoney(),
		CountryCode: 32,
	})

	/*maxId := slices.MaxFunc(accounts, func(i Account, c Account) int {
		return cmp.Compare(i.ID, c.ID)
	})*/

	args := CreateEntriesParams{
		//AccountID: (accounts[factory.RandomInt(0, int64(len(accounts))-1)].ID),
		AccountID: account.ID,
		Amount:    factory.RandomMoney(),
	}

	entry, err := testQueries.CreateEntries(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.AccountID, args.AccountID)
	require.Equal(t, args.Amount, entry.Amount)
	//require.WithinDuration(t, args.BornAt, account.BornAt, time.S
	return entry
}

func createRandomEntryAccount(t *testing.T, idAccount int) Entry {

	args := CreateEntriesParams{
		AccountID: int64(idAccount),
		Amount:    factory.RandomMoney(),
	}

	entry, err := testQueries.CreateEntries(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, entry.AccountID, args.AccountID)
	require.Equal(t, args.Amount, entry.Amount)
	//require.WithinDuration(t, args.BornAt, account.BornAt, time.S
	return entry
}

func TestCreateEntry(t *testing.T) {
	createRandomEntry(t)
}

func TestGetEntry(t *testing.T) {
	entry := createRandomEntry(t)

	entrys, err := testQueries.GetEntry(context.Background(), entry.ID)
	require.NotEmpty(t, entrys)
	require.NoError(t, err)

	require.Equal(t, entry.AccountID, entrys.AccountID)
	require.Equal(t, entry.Amount, entrys.Amount)
	require.Equal(t, entry.ID, entrys.ID)

	require.WithinDuration(t, entry.CreatedAt, entrys.CreatedAt, time.Second)
}

func TestListEntries(t *testing.T) {

	account := createRandomAccount(t)

	for i := 0; i < 10; i++ {
		createRandomEntryAccount(t, int(account.ID))
	}

	arg := ListEntriesParams{
		AccountID: account.ID,
		Limit:     5,
		Offset:    5,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)

	require.NoError(t, err)
	require.Len(t, entries, 5)

	for _, entry := range entries {
		require.NotEmpty(t, entry)
	}
}
