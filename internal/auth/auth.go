package auth

import (
	"context"
	"encoding/json"
	"errors"
	appcore_config "go-rebuild/cmd/go-rebuild/config"
	messagebroker "go-rebuild/internal/message_broker"
	"go-rebuild/internal/model"
	"go-rebuild/internal/repository"
	"time"

	"github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

var (
	// user queue
	ExchangeName = "user_exchange"
	ExchangeType = "direct"
	QueueName    = "user_queue"

	// error
	ErrUserAlreadyExists  = errors.New("user already exists with this email")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrCreateToken        = errors.New("failed to create token")
	ErrEmptyToken         = errors.New("token is empty")
	ErrInvalidToken       = errors.New("invalid token")

	ErrSendWelcomeEmail = errors.New("failed to send welcome email")
	ErrInternalServer   = errors.New("internal server error")
	ErrCreateUser       = errors.New("failed to create user")
	ErrVerifyUser       = errors.New("failed to verify user")
	ErrVerifyToken      = errors.New("token verification failed")
	ErrUserNotFound     = errors.New("user not found")
	ErrUnauthorized     = errors.New("unauthorized access, please login")
)

type authService struct {
	secretKey   string
	userRepo    repository.UserRepository
	producerSvc messagebroker.ProducerService
}

func NewAuthService(userRepo repository.UserRepository, producerSvc messagebroker.ProducerService) Jwt {
	return &authService{
		secretKey:   appcore_config.Config.SecretKey,
		userRepo:    userRepo,
		producerSvc: producerSvc,
	}
}

func (a *authService) GenerateToken(user *model.User) (*string, error) {
	var baseLogFileds = log.Fields{
		"user_id":   user.ID,
		"layer":     "auth_service",
		"operation": "verifyToken",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,                              // Subject
		"email": user.Email,                           // Custom claim
		"exp":   time.Now().Add(1 * time.Hour).Unix(), // Expiration. expire in 1 hour
		"iat":   time.Now().Unix(),                    // Issued at
		"iss":   "auth-service",                       // Issuer
	})

	tokenStr, err := token.SignedString([]byte(appcore_config.Config.SecretKey))
	if err != nil {
		log.WithError(err).WithFields(baseLogFileds)
		return nil, ErrCreateToken
	}

	return &tokenStr, nil
}

func (a *authService) VerifyToken(tokenStr string) error {
	var baseLogFileds = log.Fields{
		"token":     tokenStr,
		"layer":     "auth_service",
		"operation": "verifyToken",
	}

	// Parse token และ validate
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// ตรวจสอบว่าใช้ signing method ที่ถูกต้อง
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.WithError(jwt.ErrSignatureInvalid).WithFields(baseLogFileds)
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(appcore_config.Config.SecretKey), nil
	})

	if err != nil {
		log.WithError(err).WithFields(baseLogFileds)
		return ErrVerifyToken
	}

	if !token.Valid {
		log.WithError(jwt.ErrTokenInvalidClaims).WithFields(baseLogFileds)
		return ErrVerifyToken
	}

	return nil
}

func (a *authService) Login(ctx context.Context, user *model.User) (*string, error) {
	log.Info("login in auth service call")
	var baseLogFileds = log.Fields{
		"user_id":   user.ID,
		"layer":     "auth_service",
		"operation": "login",
	}

	var exisUser model.User
	if err := a.userRepo.GetUserByEmail(ctx, user.Email, &exisUser); err != nil {
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

	tokenStr, err := a.GenerateToken(&exisUser)
	if err != nil {
		log.WithError(err).WithFields(baseLogFileds)
		return nil, ErrInternalServer
	}

	log.Info("login in auth service call before return")
	return tokenStr, nil

}

func (a *authService) Register(ctx context.Context, user *model.User) error {
	user.ID = primitive.NewObjectID().Hex()
	var baseLogFileds = log.Fields{
		"user_id":   user.ID,
		"layer":     "auth_service",
		"operation": "register",
	}

	if err := user.Verify(); err != nil {
		log.WithError(err).WithFields(baseLogFileds)
		return ErrVerifyUser
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	if err := user.SetPassword(user.Password); err != nil {
		log.WithError(err).WithFields(baseLogFileds)
		return ErrCreateUser
	}

	if err := a.userRepo.AddUser(ctx, user); err != nil {
		log.WithError(err).WithFields(baseLogFileds)
		return ErrCreateUser
	}

	body, err := json.Marshal(user)
	if err != nil {
		return err
	}

	packet := model.Envelope{
		Type:    "create_user",
		Payload: body,
	}

	packetByte, err := json.Marshal(packet)
	if err != nil {
		return err
	}

	mqConf := messagebroker.NewMQConfig(ExchangeName, ExchangeType, QueueName, "user.create")
	if err := a.producerSvc.Publishing(ctx, mqConf, packetByte); err != nil {
		log.WithError(err).WithFields(baseLogFileds)
		return ErrSendWelcomeEmail
	}

	return nil
}

func (a *authService) GetRoleUserByID(ctx context.Context, userID string) (*string, error) {
	var baseLogFileds = log.Fields{
		"user_id":   userID,
		"layer":     "auth_service",
		"operation": "getRoleUserByID",
	}
	var user model.User
	if err := a.userRepo.GetUserByID(ctx, userID, &user); err != nil {
		log.WithError(err).WithFields(baseLogFileds)
		return nil, err
	}
	return &user.Role, nil
}

func (a *authService) CheckAllowRoles(userID string, allowedRoles []string) bool {
	userRole, err := a.GetRoleUserByID(context.Background(), userID)
	if err != nil {
		return false
	}
	for _, allowed := range allowedRoles {
		if *userRole == allowed {
			return true
		}
	}
	return false
}


func (a *authService) GetUserIDFromToken(tokenStr string) (string, error) {
	secretKey := appcore_config.Config.SecretKey
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil || !token.Valid {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok {
		return "", errors.New("user_id not found in token")
	}

	return userID, nil
}