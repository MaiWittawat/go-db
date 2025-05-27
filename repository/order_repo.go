package repository

import (
	"context"
	"fmt"
	dbRepo "go-rebuild/db"
	"go-rebuild/model"
	module "go-rebuild/module/order"
	"go-rebuild/redis"
	"time"
)

type OrderRepo struct {
	db dbRepo.DB
	collection string
	cache redis.Cache
	keyGen *redis.KeyGenerator
}

func NewOrderRepo(db dbRepo.DB, cache redis.Cache) module.OrderRepository {
	keyGen := redis.NewKeyGenerator("orders")
	return &OrderRepo{db: db, collection: "orders", cache: cache, keyGen: keyGen}
}

func (or *OrderRepo) AddOrder(ctx context.Context, o *model.Order) error {
	if err := or.db.Create(ctx, or.collection, o); err != nil {
		return err
	}

	cacheKeyID := or.keyGen.KeyID(o.ID)
	if err := or.cache.Set(ctx, cacheKeyID, o, 15*time.Minute); err != nil {
		fmt.Println("Warning: fail to set order cache in AddOrder: ", err)
	}

	fmt.Println("set cache in AddOrder success")
	return nil
}

func (or *OrderRepo) UpdateOrder(ctx context.Context, o *model.Order, id string) error {
	if err := or.db.Update(ctx, or.collection, o, id); err != nil {
		return err
	}
	
	cacheKeyID := or.keyGen.KeyID(id)
	if err := or.cache.Delete(ctx, cacheKeyID); err != nil {
		fmt.Println("Warning: fail to clear cache order in UpdateOrder: ", err)
	}

	cacheKeyList := or.keyGen.KeyList()
	if err := or.cache.Delete(ctx, cacheKeyList); err != nil {
		fmt.Println("Warning: fail to clear cache orders in UpdateOrder: ", err)
	}

	fmt.Println("set cache in UpdateOrder success")
	return nil
}

func (or *OrderRepo) DeleteOrder(ctx context.Context, id string, order *model.Order) error {
	if err := or.db.Delete(ctx, or.collection, order, id); err != nil {
		return err
	}

	cacheKeyID := or.keyGen.KeyID(id)
	if err := or.cache.Delete(ctx, cacheKeyID); err != nil {
		fmt.Println("Warning: fail to clear cache order in DeleteOrder: ", err)
	}

	cacheKeyList := or.keyGen.KeyList()
	if err := or.cache.Delete(ctx, cacheKeyList); err != nil {
		fmt.Println("Warning: fail to clear cache orders in DeleteOrder: ", err)
	}

	fmt.Println("set cache in DeleteOrder success")
	return nil
}



func (or *OrderRepo) GetAllOrder(ctx context.Context) ([]model.Order, error) {
	var orders []model.Order
	cacheKeyList := or.keyGen.KeyList()

	if err := or.cache.Get(ctx, cacheKeyList, &orders); err == nil {
		fmt.Println("get order from cache")
		return orders, nil
	}

	if err := or.db.GetAll(ctx, or.collection, &orders); err != nil {
		return nil, err
	}

	fmt.Println("get order from db")
	if err := or.cache.Set(ctx, cacheKeyList, orders, 15*time.Minute); err != nil {
		fmt.Println("Warning: fail to set cache orders in GetAllOrder: ", err)
	}

	fmt.Println("set cache in GetAllOrder success")
	return orders, nil
}

func (or *OrderRepo) GetOrderByID(ctx context.Context, id string, order *model.Order) (err error) {
	cacheKeyID := or.keyGen.KeyID(id)
	if err = or.cache.Get(ctx, cacheKeyID, &order); err == nil {
		fmt.Println("get order from cache")
		return nil
	}

	if err = or.db.GetByID(ctx, or.collection, id, order); err != nil {
		return err
	}

	fmt.Println("get order from db") 	
	if err := or.cache.Set(ctx, cacheKeyID, order, 15*time.Minute); err != nil {
		fmt.Println("Warning: fail to set cache order in GetOrderByID: ", err)
	}

	fmt.Println("set cache in GetOrderByID success")
	return nil
}

