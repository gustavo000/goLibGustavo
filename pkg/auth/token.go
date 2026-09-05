package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role,omitempty"`
	jwt.RegisteredClaims
}

func CreateToken(userID string, role string, secret string, expiresIn time.Duration) (string, error) {
	if userID == "" {
		return "", errors.New("user id is required")
	}
	if secret == "" {
		return "", errors.New("jwt secret is required")
	}
	if expiresIn <= 0 {
		return "", errors.New("expiration must be greater than zero")
	}

	now := time.Now()
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseToken(tokenString string, secret string) (*Claims, error) {
	if tokenString == "" {
		return nil, errors.New("token is required")
	}
	if secret == "" {
		return nil, errors.New("jwt secret is required")
	}

	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if token.Method != jwt.SigningMethodHS256 {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
