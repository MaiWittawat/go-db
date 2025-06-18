package repository

import (
	"context"
	"go-rebuild/internal/model"
)

type MessageRepository interface {
	AddMesssage(ctx context.Context, msg *model.Message) error
	UpdateMessage(ctx context.Context, msg *model.Message, id string) error
	DeleteMessage(ctx context.Context, id string) error

	GetMessagesBetweenUser(ctx context.Context, senderID string, receiverID string)([]model.Message, error)
	GetMessageByID(ctx context.Context, id string, msg *model.Message) error
}

type StockRepository interface {
	AddStock(ctx context.Context, s *model.Stock) error
	UpdateStock(ctx context.Context, s *model.Stock) error
	DeleteStock(ctx context.Context, id string) error

	GetStockByProductID(ctx context.Context, productID string, stock *model.Stock) error
	GetAllStock(ctx context.Context) ([]model.Stock, error)
	GetStockByID(ctx context.Context, productID string, s *model.Stock) error
}

type OrderRepository interface {
	AddOrder(ctx context.Context, o *model.Order) error
	UpdateOrder(ctx context.Context, o *model.Order, id string) error
	DeleteOrder(ctx context.Context, id string) error

	GetAllOrder(ctx context.Context) ([]model.Order, error)
	GetOrderByID(ctx context.Context, id string, order *model.Order) error
}

type ProductRepository interface {
	AddProduct(ctx context.Context, p *model.Product) error
	UpdateProduct(ctx context.Context, p *model.Product, id string) error
	DeleteProduct(ctx context.Context, id string) error

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
