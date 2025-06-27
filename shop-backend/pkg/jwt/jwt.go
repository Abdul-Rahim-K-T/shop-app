package jwtutil

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(email, role, secret string, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
		"role":  role, // "admin" or user
		"exp":   time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
