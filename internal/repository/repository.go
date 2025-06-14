package repository

import (
	"context"
	"go-rebuild/internal/model"
)

type StockRepository interface {
	AddStock(ctx context.Context, s *model.Stock) error
	UpdateStock(ctx context.Context, s *model.Stock, id string) error
	DeleteStock(ctx context.Context, id string, s *model.Stock) error

	GetStockByProductID(ctx context.Context, productID string, stock *model.Stock) error
}

type OrderRepository interface {
	AddOrder(ctx context.Context, o *model.Order) error
	UpdateOrder(ctx context.Context, o *model.Order, id string) error
	DeleteOrder(ctx context.Context, id string, o *model.Order) error

	GetAllOrder(ctx context.Context) ([]model.Order, error)
	GetOrderByID(ctx context.Context, id string, order *model.Order) error
}

type ProductRepository interface {
	AddProduct(ctx context.Context, p *model.Product) error
	UpdateProduct(ctx context.Context, p *model.Product, id string) error
	DeleteProduct(ctx context.Context, id string, p *model.Product) error

	GetAllProduct(ctx context.Context) ([]model.Product, error)
	GetProductByID(ctx context.Context, id string, p *model.Product) error
}

type UserRepository interface {
	AddUser(ctx context.Context, u *model.User) error
	UpdateUser(ctx context.Context, u *model.User, id string) error
	DeleteUser(ctx context.Context, id string, user *model.User) error

	GetAllUser(ctx context.Context) ([]model.User, error)
	GetUserByID(ctx context.Context, id string, user *model.User) error
	GetUserByEmail(ctx context.Context, email string, user *model.User) error
}
