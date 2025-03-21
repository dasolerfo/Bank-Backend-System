package db

import (
	"context"
	"simplebank/factory"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomOwner(t *testing.T) Owner {

	args := CreateOwnerParams{
		FirstName:     factory.RandomString(7),
		FirstSurname:  factory.RandomString(7),
		SecondSurname: factory.RandomString(8),
		BornAt:        time.Now(),
		Nationality:   int32(factory.RandomInt(1, 99)),
	}

	account, err := testQueries.CreateOwner(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, args.FirstName, account.FirstName)
	require.Equal(t, args.FirstSurname, account.FirstSurname)
	//require.WithinDuration(t, args.BornAt, account.BornAt, time.Second)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.BornAt)

	return account
}

func TestCreateOwner(t *testing.T) {
	createRandomOwner(t)
}
