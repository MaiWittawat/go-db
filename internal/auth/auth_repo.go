package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists with this email")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrTokenCreationFailed = errors.New("failed to create token")
	ErrEmptyToken = errors.New("token is empty")
	ErrInvalidToken = errors.New("invalid token")
	
	ErrCreateUser = errors.New("failed to create user")
	ErrTokenVerificationFailed = errors.New("token verification failed")
	ErrUserNotFound = errors.New("user not found")
	ErrUnauthorized = errors.New("unauthorized access, please login")
)


type AuthRepositorty interface {
	GenerateToken(email string) (string, error)
	VerifyToken(token string) (string, error)
}

type AuthRepo struct {
	secretKey string
}

func NewAuthRepo(secretKey string) AuthRepositorty {
	return &AuthRepo{secretKey: secretKey,}
}

func (j *AuthRepo) GenerateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(1 * time.Hour).Unix() , // Token valid for 24 hours
	})
	tokenString, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}


func (j *AuthRepo) VerifyToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.secretKey), nil
	})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrInvalidToken
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", ErrInvalidToken
	}

	return userID, nil
}