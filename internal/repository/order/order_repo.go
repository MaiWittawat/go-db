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

func (or *orderRepo) UpdateOrder(ctx context.Context, o *model.Order, id string) error {
	if err := or.db.Update(ctx, or.collection, o, id); err != nil {
		return err
	}

	cacheKeyList := or.keyGen.KeyList()
	if err := or.cache.Delete(ctx, cacheKeyList); err != nil {
		log.Warn("failed to clear cache orders in UpdateOrder: ", err)
	}

	cacheKeyID := or.keyGen.KeyID(id)
	if err := or.cache.Delete(ctx, cacheKeyID); err != nil {
		log.Warn("failed to clear cache order in UpdateOrder: ", err)
	}

	log.Info("set cache in UpdateOrder success")
	return nil
}

func (or *orderRepo) DeleteOrder(ctx context.Context, id string, order *model.Order) error {
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

// ------------------------ Method Basic Query ------------------------
func (or *orderRepo) GetAllOrder(ctx context.Context) ([]model.Order, error) {
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

func (or *orderRepo) GetOrderByID(ctx context.Context, id string, order *model.Order) (err error) {
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
