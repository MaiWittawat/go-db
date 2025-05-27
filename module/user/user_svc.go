package user

import (
	"context"
	"go-rebuild/model"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (us *userService) SaveUser(ctx context.Context, u *model.User) error {
	u.ID = primitive.NewObjectID().Hex()
	u.Role = "USER"
	u.CreatedAt = time.Now()
	u.UpdatedAt = u.CreatedAt
	hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 4)
	u.Password = string(hash)

	if err := us.userRepo.AddUser(ctx, *u); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"user_id": u.ID,
			"layer":   "service",
			"step":    "SaveUser",
		}).Error("failed to save user")
		return ErrCreateUser
	}
	return nil
}

func (us *userService) SaveSeller(ctx context.Context, u *model.User) error {
	u.ID = primitive.NewObjectID().Hex()
	u.Role = "SELLER"
	u.CreatedAt = time.Now()
	u.UpdatedAt = u.CreatedAt
	hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 4)
	u.Password = string(hash)

	if err := us.userRepo.AddUser(ctx, *u); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"user_id": u.ID,
			"layer":   "service",
			"step":    "SaveSeller",
		}).Error("failed to save seller")
		return ErrCreateSeller
	}
	return nil
}

func (us *userService) Update(ctx context.Context, u *model.User, id string) error {
	var oldUser model.User
	if err := us.userRepo.GetUserByID(ctx, id, &oldUser); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"user_id": id,
			"layer":   "service",
			"step":    "Update",
		}).Error("failed to get user by id")
		return ErrUserNotFound
	}

	if u.Username == "" {
		u.Username = oldUser.Username
	}

	if u.Password == "" {
		u.Password = oldUser.Password
	}

	if u.Email == "" {
		u.Email = oldUser.Email
	}

	u.Role = oldUser.Role
	u.UpdatedAt = time.Now()
	if err := us.userRepo.UpdateUser(ctx, *u, id); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"user_id": id,
			"layer":   "service",
			"step":    "Update",
		}).Error("failed to update user")
		return ErrUpdateUser
	}
	return nil
}

func (us *userService) Delete(ctx context.Context, id string) error {
	var user model.User
	if err := us.userRepo.GetUserByID(ctx, id, &user); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"user_id": id,
			"layer":   "service",
			"step":    "Delete",
		}).Error("failed to get user by id")
		return ErrUserNotFound
	}

	if err := us.userRepo.DeleteUser(ctx, id, &user); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"user_id": id,
			"layer":   "service",
			"step":    "Delete",
		}).Error("failed to delete user")
		return ErrDeleteUser
	}
	return nil
}

func (us *userService) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	if err := us.userRepo.GetUserByID(ctx, id, &user); err != nil {
		log.WithError(err).WithFields(log.Fields{
			"user_id": id,
			"layer":   "service",
			"step":    "GetByID",
		}).Error("failed to get user by id")
		return nil, ErrUserNotFound
	}
	return &user, nil
}

func (us *userService) GetAll(ctx context.Context) ([]model.User, error) {
	users, err := us.userRepo.GetAllUser(ctx)
	if err != nil {
		log.WithError(err).WithFields(log.Fields{
			"layer":   "service",
			"step":    "Update",
		}).Error("failed to get all user")
		return nil, ErrUserNotFound
	}
	return users, nil
}
