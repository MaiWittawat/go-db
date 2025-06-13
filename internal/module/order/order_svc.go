package order

import (
	"context"
	"errors"
	messagebroker "go-rebuild/internal/message_broker"
	"go-rebuild/internal/model"
	"go-rebuild/internal/module"
	"go-rebuild/internal/repository"
	utils "go-rebuild/internal/utlis"
	"time"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	// stock queue
	ExchangeName = "stock_exchange"
	ExchangeType = "direct"
	QueueName    = "stock_queue"

	// error
	ErrCreateOrder = errors.New("fail to create order")
	ErrUpdateOrder = errors.New("fail to update order")
	ErrDeleteOrder = errors.New("fail to delete order")

	ErrOrderNotFound = errors.New("order not found")
	ErrChangeProduct = errors.New("can not change product")
)

type orderService struct {
	orderRepo   repository.OrderRepository
	productSvc  module.ProductService
	producerSvc messagebroker.ProducerService
}

// ------------------------ Constructor ------------------------
func NewOrderService(ordeRepo repository.OrderRepository, productSvc module.ProductService, producerSvc messagebroker.ProducerService) module.OrderService {
	return &orderService{
		orderRepo:   ordeRepo,
		productSvc:  productSvc,
		producerSvc: producerSvc,
	}
}

// ------------------------ Method Basic CUD ------------------------
func (s *orderService) Save(ctx context.Context, order *model.Order) error {
	product, err := s.productSvc.GetByID(ctx, order.ProductID)
	if err != nil {
		return err
	}

	order.ID = primitive.NewObjectID().Hex()
	order.Price = product.Price
	order.Amount = order.Quantity * order.Price
	order.Status = "PENDING"
	order.CreatedAt = time.Now()
	order.UpdatedAt = order.CreatedAt

	var baseLogFields = log.Fields{
		"order_id": order.ID,
		"layer":    "order_service",
		"step":     "order_save",
	}

	if err := s.orderRepo.AddOrder(ctx, order); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to save order")
		return ErrCreateOrder
	}
	log.Info("[Service]: Order created success:", order)

	packetByte, err := utils.BuildPacket("decrease_stock", model.Stock{ProductID: order.ProductID, Quantity: order.Quantity})
	if err != nil {
		return err
	}

	mqConf := messagebroker.NewMQConfig(ExchangeName, ExchangeType, QueueName, "stock.update")
	if err := s.producerSvc.Publishing(ctx, mqConf, packetByte); err != nil {
		return err
	}

	return nil
}

func (s *orderService) Update(ctx context.Context, orderReq *model.Order, id string) error {
	var baseLogFields = log.Fields{
		"order_id": id,
		"layer":    "order_service",
		"step":     "order_update",
	}

	product, err := s.productSvc.GetByID(ctx, orderReq.ProductID)
	if err != nil {
		return err
	}

	var currentOrder model.Order
	if err := s.orderRepo.GetOrderByID(ctx, id, &currentOrder); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to get order by id")
		return ErrOrderNotFound
	}

	if product.ID != currentOrder.ID {
		return ErrChangeProduct
	}
	
	currentOrder.UpdatedAt = time.Now()
	if err := s.orderRepo.UpdateOrder(ctx, &currentOrder, id); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to update order")
		return ErrUpdateOrder
	}

	log.Info("order updated success:", currentOrder)
	return nil
}

func (s *orderService) Delete(ctx context.Context, id string) error {
	var baseLogFields = log.Fields{
		"order_id": id,
		"layer":    "order_service",
		"step":     "order_delete",
	}

	var order model.Order
	if err := s.orderRepo.GetOrderByID(ctx, id, &order); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to get order by id")
		return ErrOrderNotFound
	}

	if err := s.orderRepo.DeleteOrder(ctx, id, &order); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to delete order")
		return ErrDeleteOrder
	}
	log.Info("order deleted success:", order)

	packetByte, err := utils.BuildPacket("increase_stock", model.Stock{ProductID: order.ProductID, Quantity: order.Quantity})
	if err != nil {
		return err
	}

	mqConf := messagebroker.NewMQConfig(ExchangeName, ExchangeType, QueueName, "stock.update")
	if err := s.producerSvc.Publishing(ctx, mqConf, packetByte); err != nil {
		return err
	}

	return nil
}

// ------------------------ Method Basic Query ------------------------
func (s *orderService) GetAll(ctx context.Context) ([]model.Order, error) {
	var baseLogFields = log.Fields{
		"layer": "order_service",
		"step":  "order_getAll",
	}

	orders, err := s.orderRepo.GetAllOrder(ctx)
	if err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to get all order")
		return nil, ErrOrderNotFound
	}

	log.Info("get all order success")
	return orders, nil
}

func (s *orderService) GetByID(ctx context.Context, id string, order *model.Order) (err error) {
	var baseLogFields = log.Fields{
		"order_id": id,
		"layer":    "order_service",
		"step":     "order_delete",
	}

	err = s.orderRepo.GetOrderByID(ctx, id, order)
	if err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("failed to get order by id")
		return ErrOrderNotFound
	}

	log.Info("get order by id success:", order)
	return nil
}
