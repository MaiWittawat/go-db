package auth

import (
	"context"
	"errors"
	appcore_config "go-rebuild/cmd/go-rebuild/config"
	messagebroker "go-rebuild/internal/message_broker"
	"go-rebuild/internal/model"
	"go-rebuild/internal/module"
	"time"

	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var (
	// error
	ErrUserAlreadyExists  = errors.New("user already exists with this email")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrCreateToken        = errors.New("failed to create token")

	ErrSendWelcomeEmail = errors.New("failed to send welcome email")
	ErrInternalServer   = errors.New("internal server error")
	ErrVerifyToken      = errors.New("token verification failed")
)

type authService struct {
	secretKey   string
	userSvc     module.UserService
	producerSvc messagebroker.ProducerService
}

func NewAuthService(userSvc module.UserService, producerSvc messagebroker.ProducerService) Jwt {
	return &authService{
		secretKey:   appcore_config.Config.SecretKey,
		userSvc:     userSvc,
		producerSvc: producerSvc,
	}
}

func (a *authService) GenerateToken(user *model.User) (*string, error) {
	var baseLogFileds = log.Fields{
		"user_id":   user.ID,
		"layer":     "auth_service",
		"operation": "verifyToken",
	}

	claims := model.Claims{
		Email: user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "auth-service",
			Subject:   user.ID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(appcore_config.Config.SecretKey))
	if err != nil {
		log.WithError(err).WithFields(baseLogFileds)
		return nil, ErrCreateToken
	}

	return &tokenStr, nil
}

func (a *authService) VerifyToken(tokenStr string) (*model.Claims, error) {
	var baseLogFileds = log.Fields{
		"token":     tokenStr,
		"layer":     "auth_service",
		"operation": "verifyToken",
	}

	claims := &model.Claims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(appcore_config.Config.SecretKey), nil
	})

	if err != nil {
		log.WithError(err).WithFields(baseLogFileds)
		return nil, ErrVerifyToken
	}

	if !token.Valid {
		log.WithError(jwt.ErrTokenInvalidClaims).WithFields(baseLogFileds)
		return nil, ErrVerifyToken
	}

	return claims, nil
}

func (a *authService) Login(ctx context.Context, user *model.User) (*string, error) {
	log.Info("login in auth service call")
	var baseLogFileds = log.Fields{
		"user_id":   user.ID,
		"layer":     "auth_service",
		"operation": "login",
	}

	exisUser, err := a.userSvc.GetByEmail(ctx, user.Email)
	if err != nil {
		log.WithError(err).WithFields(baseLogFileds)
		return nil, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(exisUser.Password), []byte(user.Password)); err != nil {
		log.WithError(err).WithFields(baseLogFileds).Error("compare password failed")
		return nil, ErrInvalidCredentials
	}

	if exisUser.Email != user.Email {
		log.WithError(errors.New("email not match")).WithFields(baseLogFileds)
		return nil, ErrInvalidCredentials
	}

	tokenStr, err := a.GenerateToken(exisUser)
	if err != nil {
		log.WithError(err).WithFields(baseLogFileds)
		return nil, ErrInternalServer
	}

	log.Info("login in auth service call before return")
	return tokenStr, nil

}

func (a *authService) Register(ctx context.Context, user *model.User) error {
	if err := a.userSvc.Save(ctx, user); err != nil {
		return err
	}
	return nil
}

func (a *authService) GetRoleUserByID(ctx context.Context, userID string) (*string, error) {
	var baseLogFileds = log.Fields{
		"user_id":   userID,
		"layer":     "auth_service",
		"operation": "getRoleUserByID",
	}

	user, err := a.userSvc.GetByID(ctx, userID)
	if err != nil {
		log.WithError(err).WithFields(baseLogFileds)
		return nil, err
	}

	return &user.Role, nil
}

func (a *authService) CheckAllowRoles(userID string, allowedRoles []string) error {
	var baseLogFileds = log.Fields{
		"user_id":   userID,
		"layer":     "auth_service",
		"operation": "checkAllowRoles",
	}

	userRole, err := a.GetRoleUserByID(context.Background(), userID)
	if err != nil {
		log.WithError(err).WithFields(baseLogFileds)
		return errors.New("user not found")
	}
	for _, allowed := range allowedRoles {
		if *userRole == allowed {
			return nil
		}
	}

	log.WithError(errors.New("role not match")).WithFields(baseLogFileds)
	return errors.New("role not match")
}
