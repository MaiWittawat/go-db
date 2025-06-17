package order

import (
	"context"
	"encoding/json"
	"errors"
	messagebroker "go-rebuild/internal/message_broker"
	"go-rebuild/internal/model"
	"go-rebuild/internal/module"
	"go-rebuild/internal/repository"
	"time"

	log "github.com/sirupsen/logrus"
)

var (
	// error
	ErrCreateOrder = errors.New("fail to create order")
	ErrUpdateOrder = errors.New("fail to update order")
	ErrDeleteOrder = errors.New("fail to delete order")

	ErrOrderNotFound = errors.New("order not found")
	ErrChangeProduct = errors.New("can not change product")
	ErrPermission    = errors.New("no permission can't delete another order")
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
func (s *orderService) Save(ctx context.Context, oReq *model.OrderReq, userID string) error {
	productResp, err := s.productSvc.GetByID(ctx, oReq.ProductID)
	if err != nil {
		return ErrCreateOrder
	}

	order := oReq.ToOrder(userID, productResp)

	var baseLogFields = log.Fields{
		"order_id": order.ID,
		"layer":    "order_service",
		"method":   "order_save",
	}

	if err := s.orderRepo.AddOrder(ctx, order); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("save order")
		return ErrCreateOrder
	}
	log.Info("[Service]: Order created success:", order)

	bodyByte, err := json.Marshal(model.Stock{ProductID: order.ProductID, Quantity: order.Quantity})
	if err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("json marshal")
		return ErrCreateOrder
	}

	mqConf := &model.MQConfig{
		ExchangeName: messagebroker.StockExchangeName,
		ExchangeType: messagebroker.StockExchangeType,
		QueueName:    messagebroker.StockQueueName,
		RoutingKey:   "stock.decrease",
	}
	if err := s.producerSvc.Publishing(ctx, mqConf, bodyByte); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("publishing")
		return ErrCreateOrder
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
		return ErrUpdateOrder
	}

	var currentOrder model.Order
	if err := s.orderRepo.GetOrderByID(ctx, id, &currentOrder); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("get order by id")
		return ErrOrderNotFound
	}

	if product.ID != currentOrder.ID {
		log.WithError(ErrChangeProduct).WithFields(baseLogFields)
		return ErrUpdateOrder
	}

	currentOrder.UpdatedAt = time.Now()
	if err := s.orderRepo.UpdateOrder(ctx, &currentOrder, id); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("update order")
		return ErrUpdateOrder
	}

	log.Info("[Service]: order updated success:", currentOrder)
	return nil
}

func (s *orderService) Delete(ctx context.Context, id string, userID string) error {
	var baseLogFields = log.Fields{
		"order_id": id,
		"layer":    "order_service",
		"method":   "order_delete",
	}

	var order model.Order
	if err := s.orderRepo.GetOrderByID(ctx, id, &order); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("get order by id")
		return ErrOrderNotFound
	}

	if order.UserID != userID {
		log.WithError(ErrPermission).WithFields(baseLogFields)
		return ErrDeleteOrder
	}

	if err := s.orderRepo.DeleteOrder(ctx, id); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("delete order")
		return ErrDeleteOrder
	}
	log.Info("[Service]: order deleted success:", order)

	bodyByte, err := json.Marshal(model.Stock{ProductID: order.ProductID, Quantity: order.Quantity})
	if err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("json marshal")
		return err
	}

	mqConf := &model.MQConfig{
		ExchangeName: messagebroker.StockExchangeName,
		ExchangeType: messagebroker.StockExchangeType,
		QueueName:    messagebroker.StockQueueName,
		RoutingKey:   "stock.increase",
	}

	if err := s.producerSvc.Publishing(ctx, mqConf, bodyByte); err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("publishing")
		return err
	}

	return nil
}

// ------------------------ Method Basic Query ------------------------
func (s *orderService) GetAll(ctx context.Context) ([]model.OrderResp, error) {
	var baseLogFields = log.Fields{
		"layer":  "order_service",
		"method": "order_getAll",
	}

	orders, err := s.orderRepo.GetAllOrder(ctx)
	if err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("get all order")
		return nil, ErrOrderNotFound
	}

	var ordersResp []model.OrderResp
	for _, order := range orders {
		oResp := order.ToOrderResp()
		ordersResp = append(ordersResp, *oResp)
	}

	log.Info("[Service]: get all order success")
	return ordersResp, nil
}

func (s *orderService) GetByID(ctx context.Context, id string) (*model.OrderResp, error) {
	var order model.Order
	var baseLogFields = log.Fields{
		"order_id": id,
		"layer":    "order_service",
		"method":   "order_delete",
	}

	err := s.orderRepo.GetOrderByID(ctx, id, &order)
	if err != nil {
		log.WithError(err).WithFields(baseLogFields).Error("get order by id")
		return nil, ErrOrderNotFound
	}

	orderResp := order.ToOrderResp()

	log.Info("[Service]: get order by id success:", order)
	return orderResp, nil
}
