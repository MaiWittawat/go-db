package order

import (
	"context"
	"go-rebuild/model"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type orderService struct {
	orderRepo OrderRepository
}

// ------------------------ Constructor ------------------------
func NewOrderService(ordeRepo OrderRepository) OrderService {
	return &orderService{orderRepo: ordeRepo}
}

// ------------------------ Method Basic CUD ------------------------
func (os *orderService) Save(ctx context.Context, o *model.Order) error {
	o.ID = primitive.NewObjectID().Hex()
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
	var order model.Order 
	if err := os.orderRepo.GetOrderByID(ctx, id, &order); err != nil {
		return err
	}

	if err := os.orderRepo.DeleteOrder(ctx, id, &order); err != nil {
		return err
	}
	return nil
}

// ------------------------ Method Basic Query ------------------------
func (os *orderService) GetAll(ctx context.Context) ([]model.Order, error) {
	orders, err := os.orderRepo.GetAllOrder(ctx)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (os *orderService) GetByID(ctx context.Context, id string, order *model.Order) (err error) {
	err = os.orderRepo.GetOrderByID(ctx, id, order)
	if err != nil {
		return err
	}
	return nil
}