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
		FirstName:      factory.RandomString(7),
		FirstSurname:   factory.RandomString(7),
		SecondSurname:  factory.RandomString(8),
		BornAt:         time.Now(),
		Nationality:    int32(factory.RandomInt(1, 99)),
		HashedPassword: "EsUnSecreto",
		Email:          factory.RandomEmail(),
	}

	owner, err := testQueries.CreateOwner(context.Background(), args)

	require.NoError(t, err)
	require.NotEmpty(t, owner)

	require.Equal(t, args.FirstName, owner.FirstName)
	require.Equal(t, args.FirstSurname, owner.FirstSurname)
	require.Equal(t, args.SecondSurname, owner.SecondSurname)
	require.Equal(t, args.Email, owner.Email)
	require.Equal(t, args.Nationality, owner.Nationality)
	require.Equal(t, args.HashedPassword, owner.HashedPassword)

	//require.WithinDuration(t, args.BornAt, account.BornAt, time.Second)

	require.NotZero(t, owner.ID)
	require.NotZero(t, owner.BornAt)

	return owner
}

func TestCreateOwner(t *testing.T) {
	createRandomOwner(t)
}

func TestGetOwner(t *testing.T) {
	owner := createRandomOwner(t)
	getOwner, err := testQueries.GetOwner(context.Background(), owner.ID)

	require.NoError(t, err)
	require.NotEmpty(t, getOwner)

	require.Equal(t, getOwner.FirstName, owner.FirstName)
	require.Equal(t, getOwner.FirstSurname, owner.FirstSurname)
	require.Equal(t, getOwner.SecondSurname, owner.SecondSurname)
	require.Equal(t, getOwner.Email, owner.Email)
	require.Equal(t, getOwner.Nationality, owner.Nationality)
	require.Equal(t, getOwner.HashedPassword, owner.HashedPassword)

	//require.WithinDuration(t, args.BornAt, account.BornAt, time.Second)

	require.NotZero(t, getOwner.ID)
	require.NotZero(t, getOwner.BornAt)

}
