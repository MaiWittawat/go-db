package repository

import (
	"context"
	dbRepo "go-rebuild/db"
	"go-rebuild/model"
	module "go-rebuild/module/order"
)

type OrderRepo struct {
	db dbRepo.DB
	collection string
}

func NewOrderRepo(db dbRepo.DB) module.OrderRepository {
	return &OrderRepo{db: db, collection: "orders"}
}

func (or *OrderRepo) AddOrder(ctx context.Context, o *model.Order) error {
	return or.db.Create(ctx, or.collection, o)
}

func (or *OrderRepo) UpdateOrder(ctx context.Context, o *model.Order, id string) error {
	return or.db.Update(ctx, or.collection, o, id)
}

func (or *OrderRepo) DeleteOrder(ctx context.Context, id string) error {
	var order model.Order
	if err := or.db.GetByID(ctx, or.collection, id, &order); err != nil {
		return err
	}

	if err := or.db.Delete(ctx, or.collection, order, id); err != nil {
		return err
	}
	return nil
}



func (or *OrderRepo) GetAllOrder(ctx context.Context) ([]model.Order, error) {
	var orders []model.Order
	if err := or.db.GetAll(ctx, or.collection, orders); err != nil {
		return nil, err
	}
	return orders, nil
}

func (or *OrderRepo) GetOrderByID(ctx context.Context, id string, order *model.Order) (err error) {
	return or.db.GetByID(ctx, or.collection, id, order)
}


