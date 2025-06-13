package user

import (
	"context"
	"errors"
	messagebroker "go-rebuild/internal/message_broker"
	"go-rebuild/internal/model"
	"go-rebuild/internal/module"
	"go-rebuild/internal/repository"
	utils "go-rebuild/internal/utlis"

	log "github.com/sirupsen/logrus"
)

var (
	// queue
	ExchangeName = "user_exchange"
	ExchangeType = "direct"
	QueueName    = "user_queue"

	// error
	ErrCreateUser       = errors.New("fail to create user")
	ErrCreateSeller     = errors.New("fail to create seller")
	ErrUpdateUser       = errors.New("fail to update user")
	ErrDeleteUser       = errors.New("fail to delete user")
	ErrHashPassword     = errors.New("password is invalid")
	ErrUserNotFound     = errors.New("user not found")
	ErrSendEmailMessage = errors.New("failed to send email message")
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

// ------------------------ Method Basic UD ------------------------
func (us *userService) Update(ctx context.Context, req *model.User, id string) error {
	var baseLogFields = log.Fields{
		"user_id":   id,
		"layer":     "user_service",
		"operation": "user_update",
	}

	if err := req.Verify(); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to verify user")
		return err
	}

	var currentUser model.User
	if err := us.userRepo.GetUserByID(ctx, id, &currentUser); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to get user by id")
		return ErrUserNotFound
	}

	currentUser.SetDefaultNotNilField(req)
	if err := us.userRepo.UpdateUser(ctx, &currentUser, id); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to update user")
		return ErrUpdateUser
	}
	log.Printf("[Service]: user {%s} updated success:", currentUser.ID)

	packetByte, err := utils.BuildPacket("update_user", currentUser)
	if err != nil {
		return err
	}

	mqConf := messagebroker.NewMQConfig(ExchangeName, ExchangeType, QueueName, "user.update")
	if err := us.producerSvc.Publishing(ctx, mqConf, packetByte); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("fail to publish")
		return err
	}

	return nil
}

func (us *userService) Delete(ctx context.Context, id string) error {
	var baseLogFields = log.Fields{
		"user_id":   id,
		"layer":     "user_service",
		"operation": "user_delete",
	}

	var user model.User
	if err := us.userRepo.GetUserByID(ctx, id, &user); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to get user by id")
		return ErrUserNotFound
	}

	if err := us.userRepo.DeleteUser(ctx, id, &user); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to delete user")
		return ErrDeleteUser
	}
	log.Printf("[Service]: user {%s} deleted success:", user.ID)
	
	return nil
}

// ------------------------ Method Basic Query ------------------------
func (us *userService) GetAll(ctx context.Context) ([]model.User, error) {
	var baseLogFields = log.Fields{
		"layer":     "user_service",
		"operation": "user_getAll",
	}

	users, err := us.userRepo.GetAllUser(ctx)
	if err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to get all user")
		return nil, ErrUserNotFound
	}

	log.Info("[Service]: get all user success")
	return users, nil
}

func (us *userService) GetByID(ctx context.Context, id string) (*model.User, error) {
	log.Info("user id from userSvc : ", id)
	var baseLogFields = log.Fields{
		"user_id":   id,
		"layer":     "user_service",
		"operation": "user_getByID",
	}

	var user model.User
	if err := us.userRepo.GetUserByID(ctx, id, &user); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to get user by id")
		return nil, ErrUserNotFound
	}

	log.Printf("[Service]: get user {%s} success:", user.ID)
	return &user, nil
}
