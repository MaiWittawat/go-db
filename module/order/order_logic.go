package order

import (
	"context"
	"go-rebuild/model"
	"time"
)

type orderService struct {
	orderRepo OrderRepository
}

func NewOrderService(ordeRepo OrderRepository) OrderService {
	return &orderService{orderRepo: ordeRepo}
}

func (os *orderService) Save(ctx context.Context, o *model.Order) error {
	o.CreatedAt = time.Now()
	o.UpdatedAt = o.CreatedAt
	if err := os.orderRepo.AddOrder(ctx, o); err != nil {
		return err
	}
	return nil
}

func (os *orderService) Update(ctx context.Context, o *model.Order, id string) error {
	o.UpdatedAt = time.Now()
	if err := os.orderRepo.UpdateOrder(ctx, o, id); err != nil {
		return err
	}
	return nil
}

func (os *orderService) Delete(ctx context.Context, id string) error {
	if err := os.orderRepo.DeleteOrder(ctx, id); err != nil {
		return err
	}
	return nil
}