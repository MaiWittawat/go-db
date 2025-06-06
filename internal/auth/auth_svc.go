package auth

import (
	"context"
	"go-rebuild/internal/model"
	"go-rebuild/internal/repository"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userRepo repository.UserRepository
	authRepo AuthRepositorty
}

func NewAuthService(userRepo repository.UserRepository, authRepo AuthRepositorty) AuthService {
	return AuthService{userRepo: userRepo, authRepo: authRepo}
}

func (s *AuthService) RegisterUser(ctx context.Context, user *model.User) error {
	user.ID = primitive.NewObjectID().Hex()
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.Role = "USER"
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 4)
	user.Password = string(hash)
	if err := s.userRepo.AddUser(ctx, user); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"user_id": user.ID,
			"layer":   "auth_service",
			"step":    "RegisterUser",
		}).Error("failed to register user")
		return ErrCreateUser
	}
	return nil
}

func (s *AuthService) RegisterSeller(ctx context.Context, user *model.User) error {
	user.ID = primitive.NewObjectID().Hex()
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.Role = "SELLER"
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 4)
	user.Password = string(hash)
	if err := s.userRepo.AddUser(ctx, user); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"user_id": user.ID,
			"layer":   "auth_service",
			"step":    "RegisterSeller",
		}).Error("failed to register user")
		return ErrCreateUser
	}
	return nil
}

func (s *AuthService) GetUserByID(ctx context.Context, id string, user *model.User) error {
	if err := s.userRepo.GetUserByID(ctx, id, user); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"user_id": id,
			"layer":   "auth_service",
			"step":    "GetUserByID",
		}).Error("failed to get user by ID")
		return ErrUserNotFound
	}
	return nil
}

func (s *AuthService) Login(ctx context.Context, userLogin *model.User) (string, error) {
	if userLogin.Email == "" {
		log.Error("user email is empty, cannot create token")
		return "", ErrUserNotFound
	}

	var tempUser model.User

	if err := s.userRepo.GetUserByEmail(ctx, userLogin.Email, &tempUser); err != nil {
		log.Error("failed to create token, user not found with email: ", userLogin.Email)
		return "", err
	}

	if userLogin.Email != tempUser.Email {
		log.Error("email mismatch, cannot create token for email: ", userLogin.Email)
		return "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(tempUser.Password), []byte(userLogin.Password)); err != nil {
		log.WithError(err).Error("failed to create token, password mismatch for email: ", userLogin.Email)
		return "", ErrInvalidCredentials
	}

	token, err := s.authRepo.GenerateToken(tempUser.ID)
	if err != nil {
		log.WithError(err).Error("failed to create token for email: ", userLogin.Email)
		return "", ErrTokenCreationFailed
	}

	log.Info("token created successfully for email: ", userLogin)
	return token, nil
}

func (s *AuthService) VerifyToken(token string) (string, error) {
	if token == "" {
		log.WithError(ErrEmptyToken).WithFields(log.Fields{
			"layer": "auth_service",
			"step":  "VerifyToken",
			"token": token,
		}).Error("token is empty, cannot verify")
		return "", ErrEmptyToken
	}
	userID, err := s.authRepo.VerifyToken(token)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"layer": "auth_service",
			"step":  "VerifyToken",
			"token": token,
		}).Error("token is empty, cannot verify")
		return "", ErrInvalidToken
	}
	return userID, nil
}
