package order

import (
	"context"
	"go-rebuild/model"
)


type OrderService interface {
	Save(ctx context.Context, o *model.Order) error
	Update(ctx context.Context, o *model.Order, id string) error
	Delete(ctx context.Context, id string) error

	GetAll(ctx context.Context) ([]model.Order, error)
	GetByID(ctx context.Context, id string, order *model.Order) error
}