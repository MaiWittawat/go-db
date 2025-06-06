package auth

import (
	"context"
	"go-rebuild/internal/model"
)




type Auth interface {
	Register(ctx context.Context, user *model.User) error
	GetUserByEmail(ctx context.Context, email string, user *model.User) error
}