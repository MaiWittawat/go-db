package user

import (
	"context"
	"go-rebuild/model"
	"time"

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
		return err
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
		return err
	}
	return nil
}

func (us *userService) Update(ctx context.Context, u *model.User, id string) error {
	var oldUser model.User
	if err := us.userRepo.GetUserByID(ctx, id, &oldUser); err != nil {
		return err
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
		return err
	}
	return nil
}

func (us *userService) Delete(ctx context.Context, id string) error {
	var user model.User
	if err := us.userRepo.GetUserByID(ctx, id, &user); err != nil {
		return err
	}

	if err := us.userRepo.DeleteUser(ctx, id, &user); err != nil {
		return err
	}
	return nil
}



func (us *userService) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	if err := us.userRepo.GetUserByID(ctx, id, &user); err != nil {
		return nil, err
	}
	return &user, nil
}



func (us *userService) GetAll(ctx context.Context) ([]model.User, error) {
	users, err := us.userRepo.GetAllUser(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}