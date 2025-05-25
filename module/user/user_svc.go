package user

import (
	"context"
	"go-rebuild/model"
)

type UserService interface {
	Save(ctx context.Context, u *model.User) error
	GetByID(ctx context.Context, id string) (*model.User, error)
	Update(ctx context.Context, u *model.User, id string) error
	Delete(ctx context.Context, id string) error
}