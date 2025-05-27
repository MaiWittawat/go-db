package order

import (
	"context"
	"go-rebuild/model"
)

type OrderRepository interface {
	AddOrder(ctx context.Context, o *model.Order) error
	UpdateOrder(ctx context.Context, o *model.Order, id string) error
	DeleteOrder(ctx context.Context, id string, o *model.Order) error
	
	GetAllOrder(ctx context.Context) ([]model.Order, error)
	GetOrderByID(ctx context.Context, id string, order *model.Order) error
}