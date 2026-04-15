package middleware

import (
	"aIBuildService/aPI/models"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type UserClaims struct {
	UUID  string `json:"uuid"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func NewUserTokenClaims(user *models.User, duration time.Duration) (*UserClaims, error) {
	strID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("error generating token ID: %w", err)
	}
	return &UserClaims{
		UUID:  user.UUID,
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        strID.String(),
			Issuer:    "https://www.lts.co.uk",
			Subject:   user.Email,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}, nil
}

func RenewAccessTokenClaims(UUID string, email string, duration time.Duration) (*UserClaims, error) {
	tokenID, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("error generating token ID: %w", err)
	}
	return &UserClaims{
		UUID:  UUID,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        tokenID.String(),
			Issuer:    "https://www.lts.co.uk",
			Subject:   email,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}, nil
}
