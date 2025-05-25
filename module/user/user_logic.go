package user

import (
	"context"
	"go-rebuild/model"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (us *userService) Save(ctx context.Context, u *model.User) error {
	u.ID = uuid.NewString()
	u.CreatedAt = time.Now()
	u.UpdatedAt = u.CreatedAt
	hash, _ := bcrypt.GenerateFromPassword([]byte(u.Password), 4)
	u.Password = string(hash)

	if err := us.userRepo.AddUser(ctx, *u); err != nil {
		return err
	}
	return nil
}

func (us *userService) GetByID(ctx context.Context, id string) (*model.User, error) {
	user, err := us.userRepo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *userService) Update(ctx context.Context, u *model.User, id string) error {
	u.UpdatedAt = time.Now()
	if err := us.userRepo.UpdateUser(ctx, *u, id); err != nil {
		return err
	}
	return nil
}
func (us *userService) Delete(ctx context.Context, id string) error {
	if err := us.userRepo.DeleteUser(ctx, id); err != nil {
		return err
	}
	return nil
}
