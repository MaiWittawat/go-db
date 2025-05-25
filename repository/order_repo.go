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
	if err := or.db.Create(ctx, or.collection, o); err != nil {
		return err
	}
	return nil
}

func (or *OrderRepo) GetOrderByID(ctx context.Context, id string) (*model.Order, error) {
	return nil, nil
}

func (or *OrderRepo) UpdateOrder(ctx context.Context, o *model.Order, id string) error {
	if err := or.db.Update(ctx, or.collection, o, id); err != nil {
		return err
	}
	return nil
}

func (or *OrderRepo) DeleteOrder(ctx context.Context, id string) error {
	if err := or.db.Delete(ctx, or.collection, id); err != nil {
		return err
	}
	return nil
}
