package port

import (
	"context"
	"go-rebuild/model"
)

type UserDB interface {
	Create(ctx context.Context, u *model.User) error
	Update(ctx context.Context, u *model.User, id string) error
	Delete(ctx context.Context, id string) error
	FindByID(ctx context.Context, id string) (*model.User, error)
}