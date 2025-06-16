package model

import (
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Email string // Custom claim
	jwt.RegisteredClaims
}

// Default in jwt.RegisteredClaims เหมือน gorm
// Sub   string    // Subject
// Exp   time.Time // Expiration. expire in 1 hour
// Iat   time.Time // Issued at
// Iss   string    // Issuer
