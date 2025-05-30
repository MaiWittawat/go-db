package user

import (
	"context"
	"go-rebuild/model"
)


type UserRepository interface {
	AddUser(ctx context.Context, u model.User) error
	UpdateUser(ctx context.Context, u model.User, id string) error
    DeleteUser(ctx context.Context, id string, user *model.User) error

	
	GetAllUser(ctx context.Context) ([]model.User, error)
    GetUserByID(ctx context.Context, id string, user *model.User) error
	GetUserByEmail(ctx context.Context, email string, user *model.User) error
}       