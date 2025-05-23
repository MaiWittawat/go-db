package service

import (
	"context"
	"go-rebuild/model"
	"go-rebuild/module/port"
)

type userService struct {
	userRepo port.UserRepository
}

func NewUserService(userRepo port.UserRepository) port.UserService {
	return &userService{userRepo: userRepo}
}

func (us *userService) Save(ctx context.Context, u *model.User) error {
	if err := us.userRepo.AddUser(ctx, *u); err != nil {
		return err
	}
	return nil
}
func (us *userService) Update(ctx context.Context, u *model.User, id string) error {
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