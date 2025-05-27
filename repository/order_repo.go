package repository

import (
	"context"

	dbRepo "go-rebuild/db"
	"go-rebuild/model"
	module "go-rebuild/module/order"
	"go-rebuild/redis"
	"time"

	log "github.com/sirupsen/logrus"
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
		log.Warn("failed to set order cache in AddOrder: ", err)
	}

	log.Info("set cache in AddOrder success")
	return nil
}

func (or *OrderRepo) UpdateOrder(ctx context.Context, o *model.Order, id string) error {
	if err := or.db.Update(ctx, or.collection, o, id); err != nil {
		return err
	}
	
	cacheKeyID := or.keyGen.KeyID(id)
	if err := or.cache.Delete(ctx, cacheKeyID); err != nil {
		log.Warn("failed to clear cache order in UpdateOrder: ", err)
	}

	cacheKeyList := or.keyGen.KeyList()
	if err := or.cache.Delete(ctx, cacheKeyList); err != nil {
		log.Warn("failed to clear cache orders in UpdateOrder: ", err)
	}

	log.Info("set cache in UpdateOrder success")
	return nil
}

func (or *OrderRepo) DeleteOrder(ctx context.Context, id string, order *model.Order) error {
	if err := or.db.Delete(ctx, or.collection, order, id); err != nil {
		return err
	}

	cacheKeyID := or.keyGen.KeyID(id)
	if err := or.cache.Delete(ctx, cacheKeyID); err != nil {
		log.Warn("failed to clear cache order in DeleteOrder: ", err)
	}

	cacheKeyList := or.keyGen.KeyList()
	if err := or.cache.Delete(ctx, cacheKeyList); err != nil {
		log.Warn("failed to clear cache orders in DeleteOrder: ", err)
	}

	log.Info("set cache in DeleteOrder success")
	return nil
}



func (or *OrderRepo) GetAllOrder(ctx context.Context) ([]model.Order, error) {
	var orders []model.Order
	cacheKeyList := or.keyGen.KeyList()

	if err := or.cache.Get(ctx, cacheKeyList, &orders); err == nil {
		return orders, nil
	}

	if err := or.db.GetAll(ctx, or.collection, &orders); err != nil {
		return nil, err
	}

	if err := or.cache.Set(ctx, cacheKeyList, orders, 15*time.Minute); err != nil {
		log.Warn("failed to set cache orders in GetAllOrder: ", err)
	}

	log.Info("set cache in GetAllOrder success")
	return orders, nil
}

func (or *OrderRepo) GetOrderByID(ctx context.Context, id string, order *model.Order) (err error) {
	cacheKeyID := or.keyGen.KeyID(id)
	if err = or.cache.Get(ctx, cacheKeyID, &order); err == nil {
		return nil
	}

	if err = or.db.GetByID(ctx, or.collection, id, order); err != nil {
		return err
	}

	if err := or.cache.Set(ctx, cacheKeyID, order, 15*time.Minute); err != nil {
		log.Warn("failed to set cache order in GetOrderByID: ", err)
	}

	log.Info("set cache in GetOrderByID success")
	return nil
}

