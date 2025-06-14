package module

import (
	"context"
	"go-rebuild/internal/model"
)

type StockService interface {
	Save(ctx context.Context, productID string, quantity int) error
	Update(ctx context.Context, productID string, quantity int) error
	IncreaseQuantity(ctx context.Context, q int, id string) error
	DecreaseQuantity(ctx context.Context, q int, id string) error
	Delete(ctx context.Context, id string) error
}

type OrderService interface {
	Save(ctx context.Context, o *model.Order) error
	Update(ctx context.Context, o *model.Order, id string) error
	Delete(ctx context.Context, id string, userID string) error

	GetAll(ctx context.Context) ([]model.Order, error)
	GetByID(ctx context.Context, id string, order *model.Order) error
}

type ProductService interface {
	Save(ctx context.Context, p *model.ProductReq, userID string) error
	Update(ctx context.Context, p *model.ProductReq, id string, userID string) error
	Delete(ctx context.Context, id string) error

	GetAll(ctx context.Context) ([]model.ProductRes, error)
	GetByID(ctx context.Context, id string) (*model.ProductRes, error)
}

type UserService interface {
	Update(ctx context.Context, u *model.User, id string) error
	Delete(ctx context.Context, id string) error

	GetAll(ctx context.Context) ([]model.User, error)
	GetByID(ctx context.Context, id string) (*model.User, error)
}
