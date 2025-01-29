package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	args := CreateAccountParams{
		OwnerID:     1,
		Currency:    "USD",
		Money:       100,
		CountryCode: 34,
	}

	account, err := testQueries.CreateAccount(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, args.OwnerID, account.OwnerID)
	require.Equal(t, args.Currency, account.Currency)
	require.Equal(t, args.Money, account.Money)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

}
