package order

import (
	"context"

	"go-rebuild/internal/cache"
	dbRepo "go-rebuild/internal/db"
	"go-rebuild/internal/model"
	"go-rebuild/internal/repository"
	"time"

	log "github.com/sirupsen/logrus"
)

type orderRepo struct {
	db         dbRepo.DB
	collection string
	cacheSvc   cache.Cache
	keyGen     *cache.KeyGenerator
}

// ------------------------ Constructor ------------------------
func NewOrderRepo(db dbRepo.DB, cacheSvc cache.Cache) repository.OrderRepository {
	keyGen := cache.NewKeyGenerator("orders")
	return &orderRepo{
		db:         db,
		collection: "orders",
		cacheSvc:   cacheSvc,
		keyGen:     keyGen,
	}
}

// ------------------------ Method Basic CUD ------------------------
func (r *orderRepo) AddOrder(ctx context.Context, o *model.Order) error {
	// save order to db
	if err := r.db.Create(ctx, r.collection, o); err != nil {
		return err
	}

	// set order cache
	cacheKeyID := r.keyGen.KeyID(o.ID)
	if err := r.cacheSvc.Set(ctx, cacheKeyID, o, 15*time.Minute); err != nil {
		log.Warn("[Repo]: failed to set order cache in AddOrder: ", err)
	}

	return nil
}

func (r *orderRepo) UpdateOrder(ctx context.Context, o *model.Order, id string) error {
	// update order in db
	if err := r.db.Update(ctx, r.collection, o, id); err != nil {
		return err
	}

	// clear order cachelist in redis
	cacheKeyList := r.keyGen.KeyList()
	if err := r.cacheSvc.Delete(ctx, cacheKeyList); err != nil {
		log.Warn("[Repo]: failed to clear cache orders in UpdateOrder: ", err)
	}

	//  clear order cacheKeyID in redis
	cacheKeyID := r.keyGen.KeyID(id)
	if err := r.cacheSvc.Delete(ctx, cacheKeyID); err != nil {
		log.Warn("[Repo]: failed to clear cache order in UpdateOrder: ", err)
	}

	return nil
}

func (r *orderRepo) DeleteOrder(ctx context.Context, id string) error {
	// delete order in db
	if err := r.db.Delete(ctx, r.collection, &model.Order{}, id); err != nil {
		return err
	}

	// clear cache in redis
	cacheKeyID := r.keyGen.KeyID(id)
	if err := r.cacheSvc.Delete(ctx, cacheKeyID); err != nil {
		log.Warn("[Repo]: failed to clear cache order in DeleteOrder: ", err)
	}

	// clear cache in redis
	cacheKeyList := r.keyGen.KeyList()
	if err := r.cacheSvc.Delete(ctx, cacheKeyList); err != nil {
		log.Warn("[Repo]: failed to clear cache orders in DeleteOrder: ", err)
	}

	return nil
}

// ------------------------ Method Basic Query ------------------------
func (r *orderRepo) GetAllOrder(ctx context.Context) ([]model.Order, error) {
	var orders []model.Order
	cacheKeyList := r.keyGen.KeyList()

	// get orders in redis
	if err := r.cacheSvc.Get(ctx, cacheKeyList, &orders); err == nil {
		log.Info("[Repo]: get orders from cache: ", orders)
		return orders, nil
	}

	// get orders in db
	if err := r.db.GetAll(ctx, r.collection, &orders); err != nil {
		log.Info("[Repo]: get orders from db")
		return nil, err
	}

	// set orders cache
	if err := r.cacheSvc.Set(ctx, cacheKeyList, orders, 15*time.Minute); err != nil {
		log.Warn("[Repo]: failed to set cache orders in GetAllOrder: ", err)
	}

	return orders, nil
}

func (r *orderRepo) GetOrderByID(ctx context.Context, id string, order *model.Order) (err error) {
	// get order from redis
	cacheKeyID := r.keyGen.KeyID(id)
	if err = r.cacheSvc.Get(ctx, cacheKeyID, &order); err == nil {
		log.Info("[Repo]: get order from cache: ", order)
		return nil
	}

	// get order from db
	if err = r.db.GetByID(ctx, r.collection, id, order); err != nil {
		log.Info("[Repo]: get order from db: ", order)
		return err
	}

	// set order cache in redis
	if err := r.cacheSvc.Set(ctx, cacheKeyID, order, 15*time.Minute); err != nil {
		log.Warn("[Repo]: failed to set cache order in GetOrderByID: ", err)
	}

	return nil
}
