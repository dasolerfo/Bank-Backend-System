package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ExpiredTokenError = errors.New("token has expired")
	InvalidTokenError = errors.New("invalid token")
)

type Payload struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiredAt time.Time `json:"expired_at"`
}

// GetExpirationTime implements jwt.Claims.
func (p *Payload) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.ExpiredAt), nil
}

// GetIssuedAt implements jwt.Claims.
func (p *Payload) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(p.IssuedAt), nil
}

// GetIssuer implements jwt.Claims.
func (p *Payload) GetIssuer() (string, error) {
	if p == nil {
		return "", errors.New("payload is nil")
	}
	return p.Username, nil
}

// GetNotBefore implements jwt.Claims.
func (p *Payload) GetNotBefore() (*jwt.NumericDate, error) {
	if p == nil {
		return nil, errors.New("payload is nil")
	}
	return jwt.NewNumericDate(p.IssuedAt), nil
}

// GetSubject implements jwt.Claims.
func (p *Payload) GetSubject() (string, error) {
	if p == nil {
		return "", errors.New("payload is nil")
	}
	if p.ID == uuid.Nil {
		return "", errors.New("payload ID is nil")
	}
	return p.ID.String(), nil

}

// NewPayload creates a new Payload with a unique ID, username, issued time, and expiration time.
// It returns an error if the UUID generation fails.
func NewPayload(username string, duration time.Duration) (*Payload, error) {
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, err
	}

	payload := &Payload{
		ID:        id,
		Username:  username,
		IssuedAt:  time.Now(),
		ExpiredAt: time.Now().Add(duration),
	}

	return payload, nil
}

func (p *Payload) Valid() error {
	if time.Now().After(p.ExpiredAt) {
		return ExpiredTokenError
	}
	return nil
}

func (p *Payload) GetAudience() (jwt.ClaimStrings, error) {
	if p == nil {
		return nil, errors.New("payload is nil")
	}
	audience := jwt.ClaimStrings{}
	audience = append(audience, p.Username)
	return audience, nil
}
