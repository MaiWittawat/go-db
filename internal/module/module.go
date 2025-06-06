package module

import (
	"context"
	"go-rebuild/internal/model"
)

type OrderService interface {
	Save(ctx context.Context, o *model.Order) error
	Update(ctx context.Context, o *model.Order, id string) error
	Delete(ctx context.Context, id string) error

	GetAll(ctx context.Context) ([]model.Order, error)
	GetByID(ctx context.Context, id string, order *model.Order) error
}

type ProductService interface {
	Save(ctx context.Context, p *model.Product) error
	Update(ctx context.Context, p *model.Product, id string) error
	Delete(ctx context.Context, id string) error

	GetAll(ctx context.Context) ([]model.Product, error)
	GetByID(ctx context.Context, id string, product *model.Product) error
}

type UserService interface {
	Update(ctx context.Context, u *model.User, id string) error
	Delete(ctx context.Context, id string) error

	GetAll(ctx context.Context) ([]model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
}
