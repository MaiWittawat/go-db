package order

import (
	"context"
	"errors"
	"go-rebuild/internal/model"
	"go-rebuild/internal/module"
	"go-rebuild/internal/repository"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrCreateOrder = errors.New("fail to create order")
	ErrUpdateOrder = errors.New("fail to update order")
	ErrDeleteOrder = errors.New("fail to delete order")

	ErrOrderNotFound = errors.New("order not found")
)


type orderService struct {
	orderRepo repository.OrderRepository
}

// ------------------------ Constructor ------------------------
func NewOrderService(ordeRepo repository.OrderRepository) module.OrderService {
	return &orderService{orderRepo: ordeRepo}
}

// ------------------------ Method Basic CUD ------------------------
func (os *orderService) Save(ctx context.Context, order *model.Order) error {
	order.ID = primitive.NewObjectID().Hex()
	order.CreatedAt = time.Now()
	order.UpdatedAt = order.CreatedAt

	var baseLogFields = log.Fields{
		"order_id": order.ID,
		"layer":    "order_service",
		"step":     "order_save",
	}

	if err := os.orderRepo.AddOrder(ctx, order); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to save order")
		return ErrCreateOrder
	}

	log.Info("order created success:", order)
	return nil
}

func (os *orderService) Update(ctx context.Context, order *model.Order, id string) error {
	var baseLogFields = log.Fields{
		"order_id": id,
		"layer":    "order_service",
		"step":     "order_update",
	}

	order.UpdatedAt = time.Now()
	if err := os.orderRepo.GetOrderByID(ctx, id, order); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to get order by id")
		return ErrOrderNotFound
	}

	if err := os.orderRepo.UpdateOrder(ctx, order, id); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to update order")
		return ErrUpdateOrder
	}

	log.Info("order updated success:", order)
	return nil
}

func (os *orderService) Delete(ctx context.Context, id string) error {
	var baseLogFields = log.Fields{
		"order_id": id,
		"layer":    "order_service",
		"step":     "order_delete",
	}

	var order model.Order
	if err := os.orderRepo.GetOrderByID(ctx, id, &order); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to get order by id")
		return ErrOrderNotFound
	}

	if err := os.orderRepo.DeleteOrder(ctx, id, &order); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to delete order")
		return ErrDeleteOrder
	}

	log.Info("order deleted success:", order)
	return nil
}

// ------------------------ Method Basic Query ------------------------
func (os *orderService) GetAll(ctx context.Context) ([]model.Order, error) {
	var baseLogFields = log.Fields{
		"layer": "order_service",
		"step":  "order_getAll",
	}

	orders, err := os.orderRepo.GetAllOrder(ctx)
	if err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to get all order")
		return nil, ErrOrderNotFound
	}

	log.Info("get all order success")
	return orders, nil
}

func (os *orderService) GetByID(ctx context.Context, id string, order *model.Order) (err error) {
	var baseLogFields = log.Fields{
		"order_id": id,
		"layer":    "order_service",
		"step":     "order_delete",
	}

	err = os.orderRepo.GetOrderByID(ctx, id, order)
	if err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to get order by id")
		return ErrOrderNotFound
	}

	log.Info("get order by id success:", order)
	return nil
}
