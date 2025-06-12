package order

import (
	"context"

	dbRepo "go-rebuild/internal/db"
	"go-rebuild/internal/model"
	rclient "go-rebuild/internal/cache"
	"go-rebuild/internal/repository"
	"time"

	log "github.com/sirupsen/logrus"
)

type orderRepo struct {
	db         dbRepo.DB
	collection string
	cache      rclient.Cache
	keyGen     *rclient.KeyGenerator
}

// ------------------------ Constructor ------------------------
func NewOrderRepo(db dbRepo.DB, cache rclient.Cache) repository.OrderRepository {
	keyGen := rclient.NewKeyGenerator("orders")
	return &orderRepo{db: db, collection: "orders", cache: cache, keyGen: keyGen}
}

// ------------------------ Method Basic CUD ------------------------
func (or *orderRepo) AddOrder(ctx context.Context, o *model.Order) error {
	// save order to db
	if err := or.db.Create(ctx, or.collection, o); err != nil {
		return err
	}

	// set order cache
	cacheKeyID := or.keyGen.KeyID(o.ID)
	if err := or.cache.Set(ctx, cacheKeyID, o, 15*time.Minute); err != nil {
		log.Warn("failed to set order cache in AddOrder: ", err)
	}

	return nil
}

func (or *orderRepo) UpdateOrder(ctx context.Context, o *model.Order, id string) error {
	// update order in db
	if err := or.db.Update(ctx, or.collection, o, id); err != nil {
		return err
	}

	// clear order cachelist in redis
	cacheKeyList := or.keyGen.KeyList()
	if err := or.cache.Delete(ctx, cacheKeyList); err != nil {
		log.Warn("failed to clear cache orders in UpdateOrder: ", err)
	}

	//  clear order cacheKeyID in redis
	cacheKeyID := or.keyGen.KeyID(id)
	if err := or.cache.Delete(ctx, cacheKeyID); err != nil {
		log.Warn("failed to clear cache order in UpdateOrder: ", err)
	}

	return nil
}

func (or *orderRepo) DeleteOrder(ctx context.Context, id string, order *model.Order) error {
	// delete order in db
	if err := or.db.Delete(ctx, or.collection, order, id); err != nil {
		return err
	}

	// clear cache in redis
	cacheKeyID := or.keyGen.KeyID(id)
	if err := or.cache.Delete(ctx, cacheKeyID); err != nil {
		log.Warn("failed to clear cache order in DeleteOrder: ", err)
	}

	// clear cache in redis
	cacheKeyList := or.keyGen.KeyList()
	if err := or.cache.Delete(ctx, cacheKeyList); err != nil {
		log.Warn("failed to clear cache orders in DeleteOrder: ", err)
	}

	return nil
}

// ------------------------ Method Basic Query ------------------------
func (or *orderRepo) GetAllOrder(ctx context.Context) ([]model.Order, error) {
	var orders []model.Order
	cacheKeyList := or.keyGen.KeyList()

	// get orders in redis
	if err := or.cache.Get(ctx, cacheKeyList, &orders); err == nil {
		log.Info("get orders from cache: ", orders)
		return orders, nil
	}

	// get orders in db
	if err := or.db.GetAll(ctx, or.collection, &orders); err != nil {
		log.Info("get orders from db")
		return nil, err
	}

	// set orders cache
	if err := or.cache.Set(ctx, cacheKeyList, orders, 15*time.Minute); err != nil {
		log.Warn("failed to set cache orders in GetAllOrder: ", err)
	}

	return orders, nil
}

func (or *orderRepo) GetOrderByID(ctx context.Context, id string, order *model.Order) (err error) {
	// get order from redis
	cacheKeyID := or.keyGen.KeyID(id)
	if err = or.cache.Get(ctx, cacheKeyID, &order); err == nil {
		log.Info("get order from cache: ", order)
		return nil
	}

	// get order from db
	if err = or.db.GetByID(ctx, or.collection, id, order); err != nil {
		log.Info("get order from db: ", order)
		return err
	}

	// set order cache in redis
	if err := or.cache.Set(ctx, cacheKeyID, order, 15*time.Minute); err != nil {
		log.Warn("failed to set cache order in GetOrderByID: ", err)
	}

	return nil
}
