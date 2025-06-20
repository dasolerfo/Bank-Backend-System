package token

import (
	"simplebank/factory"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestJWTMaker(t *testing.T) {
	// Create a new JWT maker
	maker, err := NewJWTMaker(factory.RandomString(32)) // Ensure the key is at least 32 characters
	require.NoError(t, err)

	// Test token creation and verification
	username := factory.RandomString(32)
	duration := time.Minute

	issuedAt := time.Now()
	expiredAt := issuedAt.Add(duration)

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	payload, err := maker.VerifyToken(token)
	require.NoError(t, err)
	require.NotEmpty(t, payload)
	require.NotZero(t, payload.ID)
	require.Equal(t, username, payload.Username)
	require.WithinDuration(t, issuedAt, payload.IssuedAt, time.Second)
	require.WithinDuration(t, expiredAt, payload.ExpiredAt, time.Second)

	if payload.Username != username {
		t.Errorf("Expected username %s, got %s", username, payload.Username)
	}
}

func TestExpiredJWTToken(t *testing.T) {
	// Create a new JWT maker
	maker, err := NewJWTMaker(factory.RandomString(32)) // Ensure the key is at least 32 characters
	require.NoError(t, err)

	// Test token creation with a short duration
	username := factory.RandomString(32)
	duration := -time.Minute // Negative duration to simulate expiration

	token, err := maker.CreateToken(username, duration)
	require.NoError(t, err)
	require.NotEmpty(t, token)

	// Verify the expired token
	payload, err := maker.VerifyToken(token)
	require.Error(t, err)
	require.Nil(t, payload)
	require.EqualError(t, err, ExpiredTokenError.Error())
}

func TestInvalidJWTToken(t *testing.T) {
	// Create a new JWT maker
	maker, err := NewJWTMaker(factory.RandomString(32)) // Ensure the key is at least 32 characters
	require.NoError(t, err)

	// Test with an invalid token
	invalidToken := "this.is.an.invalid.token"

	payload, err := maker.VerifyToken(invalidToken)
	require.Error(t, err)
	require.Nil(t, payload)
	require.EqualError(t, err, InvalidTokenError.Error())
}
