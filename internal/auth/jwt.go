package auth

import (
	"context"
	"go-rebuild/internal/model"
)


type Jwt interface {
	Login(ctx context.Context, user *model.User) (*string, error)  
	Register(ctx context.Context, user *model.User) error
	GetRoleUserByID(ctx context.Context, userID string) (*string, error)
	GenerateToken(user *model.User) (*string, error)
	
	VerifyToken(token string) (*model.Claims, error)
	CheckAllowRoles(userID string, allowedRoles []string) bool
}
