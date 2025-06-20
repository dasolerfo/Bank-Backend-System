package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTMaker struct {
	secretKey string
}

const minKeySize = 32 // Minimum key size for HMAC-SHA256

// NewJWTMaker creates a new JWTMaker with the provided secret key.
func NewJWTMaker(secretKey string) (Maker, error) {
	if len(secretKey) < minKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minKeySize)
	}
	return &JWTMaker{secretKey: secretKey}, nil
}

// CreateToken creates a new JWT token for the given username and duration.
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwt.SignedString([]byte(maker.secretKey), token)
}

// VerifyToken checks if the JWT token is valid and returns the payload if it is.
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	payload, err := VerifyToken(token, maker.secretKey)
	if err != nil {
		return nil, err
	}
	if time.Now().After(payload.ExpiredAt) {
		return nil, fmt.Errorf("token has expired")
	}
	return payload, nil
}
