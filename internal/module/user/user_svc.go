package user

import (
	"context"
	"encoding/json"
	"errors"
	messagebroker "go-rebuild/internal/message_broker"
	"go-rebuild/internal/model"
	"go-rebuild/internal/module"
	"go-rebuild/internal/repository"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	// error
	ErrCreateUser   = errors.New("fail to create user")
	ErrUpdateUser   = errors.New("fail to update user")
	ErrDeleteUser   = errors.New("fail to delete user")
	ErrUserNotFound = errors.New("user not found")

	ErrSendEmailMessage = errors.New("failed to send email message")
	ErrVerifyUser       = errors.New("failed to verify user")
	ErrSendWelcomeEmail = errors.New("failed to send welcome email")
	ErrMarShal          = errors.New("failed to marshal object")
)

type userService struct {
	userRepo    repository.UserRepository
	producerSvc messagebroker.ProducerService
}

// ------------------------ Constructor ------------------------
func NewUserService(userRepo repository.UserRepository, producerSvc messagebroker.ProducerService) module.UserService {
	return &userService{
		userRepo:    userRepo,
		producerSvc: producerSvc,
	}
}

// ------------------------ Method Basic CUD ------------------------
func (us *userService) Save(ctx context.Context, user *model.User) error {
	user.ID = primitive.NewObjectID().Hex()
	var baseLogFields = log.Fields{
		"user_id": user.ID,
		"layer":   "user_service",
		"method":  "user_save",
	}

	if err := user.Verify(); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("verify")
		return ErrVerifyUser
	}

	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	if err := user.SetPassword(user.Password); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("set password")
		return ErrCreateUser
	}

	if err := us.userRepo.AddUser(ctx, user); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("add user")
		return ErrCreateUser
	}

	bodyByte, err := json.Marshal(user)
	if err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("marshal")
		return ErrMarShal
	}

	mqConf := &model.MQConfig{ExchangeName: messagebroker.UserExchangeName, ExchangeType: messagebroker.UserExchangeType, QueueName: messagebroker.UserQueueName, RoutingKey: "user.create"}
	if err := us.producerSvc.Publishing(ctx, mqConf, bodyByte); err != nil {
		log.WithError(err).WithFields(baseLogFields)
		return ErrSendWelcomeEmail
	}

	return nil
}

func (us *userService) Update(ctx context.Context, req *model.User, id string) error {
	var baseLogFields = log.Fields{
		"user_id": id,
		"layer":   "user_service",
		"method":  "user_update",
	}

	if err := req.Verify(); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("verify")
		return err
	}

	var currentUser model.User
	if err := us.userRepo.GetUserByID(ctx, id, &currentUser); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("get user by id")
		return ErrUserNotFound
	}

	currentUser.SetDefaultNotNilField(req)
	if err := us.userRepo.UpdateUser(ctx, &currentUser, id); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("update user")
		return ErrUpdateUser
	}
	log.Printf("[Service]: user {%s} updated success:", currentUser.ID)

	bodyByte, err := json.Marshal(currentUser)
	if err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("marshal")
		return ErrMarShal
	}

	mqConf := &model.MQConfig{ExchangeName: messagebroker.UserExchangeName, ExchangeType: messagebroker.UserExchangeType, QueueName: messagebroker.UserQueueName, RoutingKey: "user.update"}
	if err := us.producerSvc.Publishing(ctx, mqConf, bodyByte); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("publishing")
		return ErrSendEmailMessage
	}

	return nil
}

func (us *userService) Delete(ctx context.Context, id string) error {
	var baseLogFields = log.Fields{
		"user_id": id,
		"layer":   "user_service",
		"method":  "user_delete",
	}

	var user model.User
	if err := us.userRepo.GetUserByID(ctx, id, &user); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("get user by id")
		return ErrUserNotFound
	}

	if err := us.userRepo.DeleteUser(ctx, id, &user); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("delete user")
		return ErrDeleteUser
	}
	log.Printf("[Service]: user {%s} deleted success:", user.ID)

	return nil
}

// ------------------------ Method Basic Query ------------------------
func (us *userService) GetAll(ctx context.Context) ([]model.User, error) {
	var baseLogFields = log.Fields{
		"layer":  "user_service",
		"method": "user_getAll",
	}

	users, err := us.userRepo.GetAllUser(ctx)
	if err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("get all user")
		return nil, ErrUserNotFound
	}

	log.Info("[Service]: get all user success")
	return users, nil
}

func (us *userService) GetByID(ctx context.Context, id string) (*model.User, error) {
	var baseLogFields = log.Fields{
		"user_id": id,
		"layer":   "user_service",
		"method":  "user_getByID",
	}

	var user model.User
	if err := us.userRepo.GetUserByID(ctx, id, &user); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("get user by id")
		return nil, ErrUserNotFound
	}

	log.Printf("[Service]: get user {%s} success:", user.ID)
	return &user, nil
}

func (us *userService) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var baseLogFields = log.Fields{
		"user_email": email,
		"layer":      "user_service",
		"method":     "user_getByEmail",
	}

	var user model.User
	if err := us.userRepo.GetUserByEmail(ctx, email, &user); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("get user by email")
		return nil, ErrUserNotFound
	}

	log.Printf("[Service]: get user by email {%s} success:", user.Email)
	return &user, nil
}
